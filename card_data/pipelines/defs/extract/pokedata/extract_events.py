import dagster as dg
import requests
import polars as pl

WORLDS_TCG_ID = "0000129"
WORLDS_VG_ID  = "0000115"

VALID_TYPES = ["Regional", "International", "Worlds"]


def call_api(url: str) -> dict:
    r = requests.get(url, timeout=30)
    r.raise_for_status()
    return r.json()


def infer_season(start_date: str) -> int:
    year, month = int(start_date[:4]), int(start_date[5:7])

    if month >= 9:
        return year + 1
    else:
        return year


@dg.asset(kinds={"API", "Polars"}, name="create_events_dataframe")
def create_events_dataframe() -> pl.DataFrame:
    data = call_api("https://www.pokedata.ovh/apiv2/tournaments")
    return build_events_dataframe(data)


def build_events_dataframe(data: dict) -> pl.DataFrame:
    tcg_events = [
        t for t in data["tcg"]["data"]
        if t["id"] > WORLDS_TCG_ID and any(v in t["name"] for v in VALID_TYPES)
    ]
    vg_events = [
        t for t in data["vg"]["data"]
        if t["id"] > WORLDS_VG_ID and any(v in t["name"] for v in VALID_TYPES)
    ]

    rows = []
    for t in tcg_events:
        rows.append({
            "pokedata_id":  t["id"],
            "game_type":    "TCG",
            "name":         t["name"],
            "start_date":   t["date"]["start"],
            "end_date":     t["date"]["end"],
            "season":       infer_season(t["date"]["start"]),
            "count":        int(t["players"]["masters"]),
            "rounds":       t["roundNumbers"]["masters"],
            "last_updated": t["lastUpdated"],
        })
    for t in vg_events:
        rows.append({
            "pokedata_id":  t["id"],
            "game_type":    "VGC",
            "name":         t["name"],
            "start_date":   t["date"]["start"],
            "end_date":     t["date"]["end"],
            "season":       infer_season(t["date"]["start"]),
            "count":        int(t["players"]["masters"]),
            "rounds":       t["roundNumbers"]["masters"],
            "last_updated": t["lastUpdated"],
        })

    return pl.DataFrame(rows)
