import time
import json

import dagster as dg
import polars as pl

from pydantic import BaseModel, HttpUrl, ValidationError
from typing import Optional
from termcolor import colored

import requests


class Series(BaseModel):
    id: str
    name: str
    logo: Optional[HttpUrl] = None


class Set(BaseModel):
    series_id: str
    set_id: str
    set_name: str
    official_card_count: int | None
    total_card_count: int | None
    logo: Optional[str] = None
    symbol: Optional[str] = None


@dg.asset(kinds={"API", "Polars", "Pydantic"})
def extract_series_data() -> pl.DataFrame:
    url: str = "https://api.tcgdex.net/v2/en/series"
    data = requests.get(url).json()

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


@dg.asset(kinds={"API", "Polars", "Pydantic"}, name="extract_set_data")
def extract_set_data() -> pl.DataFrame:
    url_list = [
        "https://api.tcgdex.net/v2/en/series/me",
        "https://api.tcgdex.net/v2/en/series/sv",
        "https://api.tcgdex.net/v2/en/series/swsh",
    ]

    flat: list[dict] = []

    for url in url_list:
        data = requests.get(url).json()
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


@dg.asset(kinds={"API"}, name="extract_card_url_from_set_data")
def extract_card_url_from_set() -> list:
    urls = ["https://api.tcgdex.net/v2/en/sets/me02"]

    all_card_urls = []

    for url in urls:
        try:
            r = requests.get(url)
            r.raise_for_status()

            data = r.json()["cards"]

            set_card_urls = [
                f"https://api.tcgdex.net/v2/en/cards/{card['id']}"
                for card in data
                if "-TG" not in card["id"]
            ]
            all_card_urls.extend(set_card_urls)

            time.sleep(0.1)

        except requests.RequestException as e:
            print(f"Failed to fetch set {url}: {e}")

    return all_card_urls


@dg.asset(kinds={"API"}, name="extract_card_info")
def extract_card_info(extract_card_url_from_set_data: list) -> list:
    card_url_list = extract_card_url_from_set_data
    cards_list = []

    for url in card_url_list:
        try:
            r = requests.get(url)
            r.raise_for_status()
            data = r.json()
            cards_list.append(data)
            print(f"Retrieved card: {data['id']} - {data.get('name', 'Unknown')}")
            time.sleep(0.1)
        except requests.RequestException as e:
            print(f"Failed to fetch {url}: {e}")

    return cards_list


@dg.asset(kinds={"Polars"}, name="create_card_dataframe")
def create_card_dataframe(extract_card_info: list) -> pl.DataFrame:
    cards_list = extract_card_info

    all_flat_cards = []

    for card in cards_list:
        flat = {}

        # Copy top-level scalar values
        scalar_keys = [
            "category",
            "hp",
            "id",
            "illustrator",
            "image",
            "localId",
            "name",
            "rarity",
            "regulationMark",
            "retreat",
            "stage",
        ]
        for key in scalar_keys:
            flat[key] = card.get(key)

        # Flatten nested dicts with prefixes
        for key, value in card.get("legal", {}).items():
            flat[f"legal_{key}"] = value

        for key, value in card.get("set", {}).items():
            if isinstance(value, dict):
                for sub_key, sub_val in value.items():
                    flat[f"set_{key}_{sub_key}"] = sub_val
            else:
                flat[f"set_{key}"] = value

        # Flatten types (list of strings)
        flat["types"] = ", ".join(card.get("types", []))

        flat["attacks_json"] = json.dumps(card.get("attacks", []), ensure_ascii=False)

        attacks = card.get("attacks", [])
        for i, atk in enumerate(attacks):
            prefix = f"attack_{i + 1}"
            flat[f"{prefix}_name"] = atk.get("name")
            flat[f"{prefix}_damage"] = atk.get("damage")
            flat[f"{prefix}_effect"] = atk.get("effect")
            flat[f"{prefix}_cost"] = ", ".join(atk.get("cost", []))

        all_flat_cards.append(flat)

    df = pl.DataFrame(all_flat_cards)

    return df
