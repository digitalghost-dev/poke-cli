from typing import Optional

import dagster as dg
import polars as pl
from pydantic import BaseModel, HttpUrl, ValidationError
from termcolor import colored

from ....utils.json_retriever import fetch_json


class Set(BaseModel):
    series_id: str
    set_id: str
    set_name: str
    official_card_count: int | None
    total_card_count: int | None
    logo: Optional[str] = None
    symbol: Optional[str] = None


@dg.asset(kinds={"API", "Polars", "Pydantic"}, name="extract_sets_data")
def extract_sets_data() -> pl.DataFrame:
    url_list = [
        "https://api.tcgdex.net/v2/en/series/me",
        "https://api.tcgdex.net/v2/en/series/sv",
        "https://api.tcgdex.net/v2/en/series/swsh",
    ]

    flat: list[dict] = []

    for url in url_list:
        data: dict = fetch_json(url)
        series_id = data.get("id")

        for s in data.get("sets", []):
            entry = {
                "series_id": series_id,
                "set_id": s.get("id"),
                "set_name": s.get("name"),
                "official_card_count": s.get("cardCount", {}).get("official"),
                "total_card_count": s.get("cardCount", {}).get("total"),
                "logo": s.get("logo"),
                "symbol": s.get("symbol"),
            }
            flat.append(entry)

    # Pydantic validation
    try:
        validated: list[Set] = [Set(**item) for item in flat]
        print(colored(" ✓", "green"), "Pydantic validation passed for all set entries.")
    except ValidationError as e:
        print(colored(" ✖", "red"), "Pydantic validation failed.")
        print(e)
        raise

    return pl.DataFrame([s.model_dump(mode="json") for s in validated])
