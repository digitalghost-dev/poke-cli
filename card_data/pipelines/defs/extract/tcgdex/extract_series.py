from typing import Optional

import dagster as dg
import polars as pl
from pydantic import BaseModel, HttpUrl, ValidationError
from termcolor import colored

from ....utils.json_retriever import fetch_json


class Series(BaseModel):
    id: str
    name: str
    logo: Optional[HttpUrl] = None


@dg.asset(kinds={"API", "Polars", "Pydantic"})
def extract_series_data() -> pl.DataFrame:
    url: str = "https://api.tcgdex.net/v2/en/series"
    data: dict = fetch_json(url)

    # Pydantic validation
    try:
        validated: list[Series] = [Series(**item) for item in data]
        print(
            colored(" ✓", "green"), "Pydantic validation passed for all series entries."
        )
    except ValidationError as e:
        print(colored(" ✖", "red"), "Pydantic validation failed.")
        print(e)
        raise

    filtered = [
        s.model_dump(mode="json") for s in validated if s.id in ["swsh", "sv", "me"]
    ]
    return pl.DataFrame(filtered)
