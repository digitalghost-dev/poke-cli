import re
import subprocess  # nosec
import time
from pathlib import Path

import dagster as dg
import polars as pl
import requests
from dagster import Backoff, RetryPolicy
from psycopg2.extras import Json, execute_values
from sqlalchemy import create_engine, text
from sqlalchemy.exc import OperationalError

from ..utils.secret_retriever import fetch_secret

WORLDS_TCG_ID = "0000129"
WORLDS_VG_ID = "0000115"
VALID_EVENT_TYPES = ["Regional", "International", "Worlds"]
EVENT_SOURCES = (
    ("tcg", "TCG", WORLDS_TCG_ID),
    ("vg", "VGC", WORLDS_VG_ID),
)

REQUEST_DELAY_SECONDS = 0.2
TOP_N_PLACEMENTS = 256
COUNTRY_PATTERN = re.compile(r"^(.*?)\s*\[([A-Za-z]{2,3})\]\s*$")


def call_api(url: str, session: requests.Session | None = None) -> dict:
    client = session or requests
    r = client.get(url, timeout=60)
    r.raise_for_status()
    return r.json()


def infer_season(start_date: str) -> int:
    year, month = int(start_date[:4]), int(start_date[5:7])
    if month >= 9:
        return year + 1
    return year


def parse_player_name(raw: str) -> tuple[str, str | None]:
    m = COUNTRY_PATTERN.match(raw)
    if m:
        return m.group(1).strip(), m.group(2).upper()
    return raw.strip(), None


def _run_soda_scan(checks_filename: str) -> None:
    current_file_dir = Path(__file__).parent
    result = subprocess.run(  # nosec
        [
            "soda", "scan", "-d", "supabase",
            "-c", "../soda/configuration.yml",
            f"../soda/{checks_filename}",
        ],
        capture_output=True,
        text=True,
        cwd=current_file_dir,
    )
    if result.stdout:
        print(result.stdout)
    if result.stderr:
        print(result.stderr)
    if result.returncode != 0:
        raise Exception(f"Soda check {checks_filename} failed with return code {result.returncode}")


_CREATE_EVENTS = """
    CREATE TABLE IF NOT EXISTS staging.comp_events (
        id          BIGSERIAL PRIMARY KEY,
        pokedata_id TEXT,
        game_type   TEXT,
        name        TEXT,
        start_date  TEXT,
        end_date    TEXT,
        season      BIGINT,
        count       BIGINT,
        rounds      BIGINT,
        last_updated TEXT,
        UNIQUE (pokedata_id, game_type)
    )
"""

_UPSERT_EVENTS = """
    INSERT INTO staging.comp_events (pokedata_id, game_type, name, start_date, end_date, season, count, rounds, last_updated)
    VALUES (:pokedata_id, :game_type, :name, :start_date, :end_date, :season, :count, :rounds, :last_updated)
    ON CONFLICT (pokedata_id, game_type) DO UPDATE SET
        name         = EXCLUDED.name,
        start_date   = EXCLUDED.start_date,
        end_date     = EXCLUDED.end_date,
        season       = EXCLUDED.season,
        count        = EXCLUDED.count,
        rounds       = EXCLUDED.rounds,
        last_updated = EXCLUDED.last_updated
"""


def build_events_dataframe(data: dict) -> pl.DataFrame:
    rows = []
    for source_key, game_type, min_event_id in EVENT_SOURCES:
        for event in data[source_key]["data"]:
            if event["id"] <= min_event_id or not any(
                event_type in event["name"] for event_type in VALID_EVENT_TYPES
            ):
                continue

            start_date = event["date"]["start"]
            rows.append({
                "pokedata_id": event["id"],
                "game_type": game_type,
                "name": event["name"],
                "start_date": start_date,
                "end_date": event["date"]["end"],
                "season": infer_season(start_date),
                "count": int(event["players"]["masters"]),
                "rounds": event["roundNumbers"]["masters"],
                "last_updated": event["lastUpdated"],
            })

    return pl.DataFrame(rows)


@dg.asset(kinds={"API", "Polars"}, name="create_events_dataframe")
def create_events_dataframe() -> pl.DataFrame:
    data = call_api("https://www.pokedata.ovh/apiv2/tournaments")
    return build_events_dataframe(data)


@dg.asset(
    kinds={"Supabase", "Postgres"},
    retry_policy=RetryPolicy(max_retries=3, delay=2, backoff=Backoff.EXPONENTIAL),
)
def load_events_data(create_events_dataframe: pl.DataFrame) -> None:
    database_url: str = fetch_secret()
    try:
        engine = create_engine(database_url)
        with engine.begin() as conn:
            conn.execute(text(_CREATE_EVENTS))
            conn.execute(text(_UPSERT_EVENTS), create_events_dataframe.to_dicts())
        print(" ✓ Data upserted into staging.comp_events")
    except OperationalError as e:
        print(f" ✖ Connection error in load_events_data(): {e}")
        raise


@dg.asset(deps=[load_events_data], kinds={"Soda"}, name="data_quality_checks_on_comp_events")
def data_quality_checks_on_comp_events() -> None:
    _run_soda_scan("checks_comp_events.yml")


_CREATE_PLAYERS = """
    CREATE TABLE IF NOT EXISTS staging.comp_players (
        id            BIGSERIAL PRIMARY KEY,
        pokedata_id   TEXT NOT NULL,
        game_type     TEXT NOT NULL,
        player_name   TEXT NOT NULL,
        country       TEXT,
        placement     BIGINT,
        wins          BIGINT,
        losses        BIGINT,
        ties          BIGINT,
        resistance_self    NUMERIC(4, 2),
        resistance_opp     NUMERIC(4, 2),
        resistance_oppopp  NUMERIC(4, 2),
        dropped_round BIGINT,
        trainer_name  TEXT,
        UNIQUE (pokedata_id, game_type, player_name)
    )
"""

_CREATE_ROUNDS = """
    CREATE TABLE IF NOT EXISTS staging.comp_rounds (
        id            BIGSERIAL PRIMARY KEY,
        pokedata_id   TEXT NOT NULL,
        game_type     TEXT NOT NULL,
        player_name   TEXT NOT NULL,
        round_number  BIGINT NOT NULL,
        opponent_name TEXT,
        result        TEXT,
        table_number  TEXT,
        UNIQUE (pokedata_id, game_type, player_name, round_number)
    )
"""

_CREATE_DECKLISTS_TEMPLATE = """
    CREATE TABLE IF NOT EXISTS staging.{table_name} (
        id           BIGSERIAL PRIMARY KEY,
        pokedata_id  TEXT NOT NULL,
        game_type    TEXT NOT NULL,
        player_name  TEXT NOT NULL,
        decklist     JSONB,
        UNIQUE (pokedata_id, game_type, player_name)
    )
"""

_CREATE_VG_DECKLISTS = _CREATE_DECKLISTS_TEMPLATE.format(table_name="comp_vg_decklists")
_CREATE_TCG_DECKLISTS = _CREATE_DECKLISTS_TEMPLATE.format(table_name="comp_tcg_decklists")

_INSERT_PLAYERS_SQL = """
    INSERT INTO staging.comp_players (
        pokedata_id, game_type, player_name, country, placement,
        wins, losses, ties,
        resistance_self, resistance_opp, resistance_oppopp,
        dropped_round, trainer_name
    ) VALUES %s
    ON CONFLICT (pokedata_id, game_type, player_name) DO NOTHING
"""

_INSERT_ROUNDS_SQL = """
    INSERT INTO staging.comp_rounds (
        pokedata_id, game_type, player_name, round_number,
        opponent_name, result, table_number
    ) VALUES %s
    ON CONFLICT (pokedata_id, game_type, player_name, round_number) DO NOTHING
"""

_INSERT_DECKLISTS_SQL_TEMPLATE = """
    INSERT INTO staging.{table_name} (
        pokedata_id, game_type, player_name, decklist
    ) VALUES %s
    ON CONFLICT (pokedata_id, game_type, player_name) DO NOTHING
"""

_INSERT_VG_DECKLISTS_SQL = _INSERT_DECKLISTS_SQL_TEMPLATE.format(
    table_name="comp_vg_decklists"
)
_INSERT_TCG_DECKLISTS_SQL = _INSERT_DECKLISTS_SQL_TEMPLATE.format(
    table_name="comp_tcg_decklists"
)

_PLAYER_TABLE_DDLS = (
    _CREATE_PLAYERS,
    _CREATE_ROUNDS,
    _CREATE_VG_DECKLISTS,
    _CREATE_TCG_DECKLISTS,
)


def fetch_events_to_process(conn) -> list[dict]:
    """Finished events not yet loaded into comp_players. 2-day buffer past end_date."""
    query = """
        SELECT e.pokedata_id, e.game_type
        FROM staging.comp_events e
        WHERE e.end_date::date < CURRENT_DATE - INTERVAL '1 day'
          AND NOT EXISTS (
              SELECT 1 FROM staging.comp_players p
              WHERE p.pokedata_id = e.pokedata_id
                AND p.game_type   = e.game_type
          )
        ORDER BY e.end_date
    """

    return [dict(row._mapping) for row in conn.execute(text(query)).all()]


def build_player_rows(
    data: dict,
    pid: str,
    gt: str,
) -> tuple[list[tuple], list[tuple], list[tuple], list[tuple]]:
    players: list[tuple] = []
    rounds: list[tuple] = []
    decklists_vg: list[tuple] = []
    decklists_tcg: list[tuple] = []

    for div in data.get("tournament_data", []):
        for p in div.get("data", []):
            placement = p.get("placing")
            if placement is None or placement > TOP_N_PLACEMENTS:
                continue

            name, country = parse_player_name(p["name"])

            players.append((
                pid, gt, name, country, placement,
                p["record"]["wins"], p["record"]["losses"], p["record"]["ties"],
                p["resistances"]["self"], p["resistances"]["opp"], p["resistances"]["oppopp"],
                p.get("drop", -1), p.get("Trainer name"),
            ))

            for round_num, r in (p.get("rounds") or {}).items():
                rounds.append((
                    pid, gt, name, int(round_num),
                    r.get("name"), r.get("result"), str(r.get("table", "")),
                ))

            deck = p.get("decklist")
            if deck:
                deck_row = (pid, gt, name, Json(deck))
                if gt == "VGC":
                    decklists_vg.append(deck_row)
                else:
                    decklists_tcg.append(deck_row)

    return players, rounds, decklists_vg, decklists_tcg


def process_event(
    cur,
    pid: str,
    gt: str,
    session: requests.Session | None = None,
) -> tuple[int, int, int]:
    """Fetch one event from the API and bulk-insert its data via execute_values."""
    url_segment = "tcg" if gt == "TCG" else "vg"
    url = f"https://www.pokedata.ovh/apiv2/{url_segment}/id/{pid}/division/masters"
    data = call_api(url, session=session)
    players, rounds, decklists_vg, decklists_tcg = build_player_rows(data, pid, gt)

    insert_batches = (
        (_INSERT_PLAYERS_SQL, players, 500),
        (_INSERT_ROUNDS_SQL, rounds, 1000),
        (_INSERT_VG_DECKLISTS_SQL, decklists_vg, 500),
        (_INSERT_TCG_DECKLISTS_SQL, decklists_tcg, 500),
    )
    for insert_sql, rows, page_size in insert_batches:
        if rows:
            execute_values(cur, insert_sql, rows, page_size=page_size)

    return len(players), len(rounds), len(decklists_vg) + len(decklists_tcg)


@dg.asset(
    deps=[dg.AssetKey("data_quality_checks_on_comp_events")],
    kinds={"API", "Supabase", "Postgres"},
    retry_policy=RetryPolicy(max_retries=3, delay=2, backoff=Backoff.EXPONENTIAL),
)
def load_players_data() -> None:
    database_url: str = fetch_secret()

    try:
        engine = create_engine(database_url)

        with engine.begin() as conn:
            for ddl in _PLAYER_TABLE_DDLS:
                conn.execute(text(ddl))

        with engine.connect() as conn:
            events = fetch_events_to_process(conn)

        total = len(events)
        print(f" → {total} events to process")

        if total == 0:
            print(" ✓ Nothing to load")
            return

        overall_start = time.time()
        successes = 0
        failures = 0
        totals = {"players": 0, "rounds": 0, "decklists": 0}

        with engine.connect() as conn, requests.Session() as session:
            for i, ev in enumerate(events, 1):
                pid = ev["pokedata_id"]
                gt = ev["game_type"]
                start = time.time()

                try:
                    with conn.begin():
                        cur = conn.connection.cursor()
                        n_players, n_rounds, n_decklists = process_event(
                            cur, pid, gt, session
                        )
                    successes += 1
                    totals["players"] += n_players
                    totals["rounds"] += n_rounds
                    totals["decklists"] += n_decklists
                except requests.HTTPError as e:
                    print(f" ✖ [{i:3d}/{total}] {gt} {pid}: API error {e}; skipping")
                    failures += 1
                    time.sleep(REQUEST_DELAY_SECONDS)
                    continue
                except Exception as e:  # noqa: BLE001
                    print(f" ✖ [{i:3d}/{total}] {gt} {pid}: {type(e).__name__}: {e}; skipping")
                    failures += 1
                    time.sleep(REQUEST_DELAY_SECONDS)
                    continue

                elapsed = time.time() - start
                total_elapsed = time.time() - overall_start
                eta_min = (total_elapsed / i) * (total - i) / 60
                print(
                    f" ✓ [{i:3d}/{total}] {gt} {pid}: "
                    f"{n_players} players, {n_rounds} rounds, {n_decklists} decklists "
                    f"in {elapsed:.1f}s (ETA {eta_min:.1f}min)"
                )

                time.sleep(REQUEST_DELAY_SECONDS)

        total_min = (time.time() - overall_start) / 60
        print(
            f" ✓ Done: {successes} succeeded, {failures} failed in {total_min:.1f}min "
            f"({totals['players']} players, {totals['rounds']} rounds, {totals['decklists']} decklists)"
        )

    except OperationalError as e:
        print(f" ✖ Connection error in load_players_data(): {e}")
        raise


@dg.asset(deps=[load_players_data], kinds={"Soda"}, name="data_quality_checks_on_comp_players")
def data_quality_checks_on_comp_players() -> None:
    _run_soda_scan("checks_comp_players.yml")


@dg.asset(deps=[load_players_data], kinds={"Soda"}, name="data_quality_checks_on_comp_rounds")
def data_quality_checks_on_comp_rounds() -> None:
    _run_soda_scan("checks_comp_rounds.yml")


@dg.asset(deps=[load_players_data], kinds={"Soda"}, name="data_quality_checks_on_comp_vg_decklists")
def data_quality_checks_on_comp_vg_decklists() -> None:
    _run_soda_scan("checks_comp_vg_decklists.yml")


@dg.asset(deps=[load_players_data], kinds={"Soda"}, name="data_quality_checks_on_comp_tcg_decklists")
def data_quality_checks_on_comp_tcg_decklists() -> None:
    _run_soda_scan("checks_comp_tcg_decklists.yml")
