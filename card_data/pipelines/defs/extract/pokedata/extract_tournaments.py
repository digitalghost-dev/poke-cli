import dagster as dg
import polars as pl

from .extract_events import call_api, WORLDS_TCG_ID, WORLDS_VG_ID, VALID_TYPES


def _build_dataframe(tournaments: list[dict]) -> pl.DataFrame:
    rows = [
        {
            "event_id":     t["id"],
            "name":         t["name"],
            "count":        int(t["players"]["masters"]),
            "rounds":       t["roundNumbers"]["masters"],
            "last_updated": t["lastUpdated"],
        }
        for t in tournaments
    ]
    return pl.DataFrame(rows)


@dg.asset(
    deps=[dg.AssetKey(["data_quality_checks_on_comp_events"])],
    kinds={"API", "Polars"},
    name="create_tcg_tournaments_dataframe",
)
def create_tcg_tournaments_dataframe() -> pl.DataFrame:
    data = call_api("https://www.pokedata.ovh/apiv2/tournaments")
    tournaments = [
        t for t in data["tcg"]["data"]
        if t["id"] > WORLDS_TCG_ID and any(v in t["name"] for v in VALID_TYPES)
    ]
    return _build_dataframe(tournaments)


@dg.asset(
    deps=[dg.AssetKey(["data_quality_checks_on_comp_events"])],
    kinds={"API", "Polars"},
    name="create_vg_tournaments_dataframe",
)
def create_vg_tournaments_dataframe() -> pl.DataFrame:
    data = call_api("https://www.pokedata.ovh/apiv2/tournaments")
    tournaments = [
        t for t in data["vg"]["data"]
        if t["id"] > WORLDS_VG_ID and any(v in t["name"] for v in VALID_TYPES)
    ]
    return _build_dataframe(tournaments)
