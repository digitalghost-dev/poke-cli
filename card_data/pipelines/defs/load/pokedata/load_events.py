import subprocess  # nosec
from pathlib import Path

import dagster as dg
import polars as pl
from dagster import RetryPolicy, Backoff
from sqlalchemy import create_engine, text
from sqlalchemy.exc import OperationalError

from ....utils.secret_retriever import fetch_secret

_CREATE_TABLE = """
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

_UPSERT = """
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


@dg.asset(
    kinds={"Supabase", "Postgres"},
    retry_policy=RetryPolicy(max_retries=3, delay=2, backoff=Backoff.EXPONENTIAL),
)
def load_events_data(create_events_dataframe: pl.DataFrame) -> None:
    database_url: str = fetch_secret()

    try:
        engine = create_engine(database_url)
        with engine.begin() as conn:
            conn.execute(text(_CREATE_TABLE))
            conn.execute(text(_UPSERT), create_events_dataframe.to_dicts())
        print(" ✓ Data upserted into staging.comp_events")
    except OperationalError as e:
        print(f" ✖ Connection error in load_events_data(): {e}")
        raise


@dg.asset(
    deps=[load_events_data],
    kinds={"Soda"},
    name="data_quality_checks_on_comp_events",
)
def data_quality_checks_on_comp_events() -> None:
    current_file_dir = Path(__file__).parent

    result = subprocess.run(  # nosec
        [
            "soda",
            "scan",
            "-d",
            "supabase",
            "-c",
            "../../../soda/configuration.yml",
            "../../../soda/checks_comp_events.yml",
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
        raise Exception(f"Soda data quality checks failed with return code {result.returncode}")
