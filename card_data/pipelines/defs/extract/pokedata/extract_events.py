import re
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


def normalize_name(name: str) -> str:
    name = re.sub(r'\bTCG\b', '', name, flags=re.IGNORECASE)
    name = re.sub(r'\bVGC?\b', '', name, flags=re.IGNORECASE)
    name = re.sub(r'\b20\d{2}\b', '', name)
    return re.sub(r'\s+', ' ', name).strip(' --')


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

    tcg_by_name = {normalize_name(t["name"]): t for t in tcg_events}
    vg_by_name  = {normalize_name(t["name"]): t for t in vg_events}

    all_names = set(tcg_by_name) | set(vg_by_name)

    rows = []
    for name in sorted(all_names):
        tcg = tcg_by_name.get(name)
        vg  = vg_by_name.get(name)
        source = tcg or vg
        rows.append({
            "start_date":      source["date"]["start"],
            "end_date":        source["date"]["end"],
            "season":          infer_season(source["date"]["start"]),
            "tcg_pokedata_id": tcg["id"] if tcg else None,
            "vg_pokedata_id":  vg["id"]  if vg  else None,
        })

    return pl.DataFrame(rows)


