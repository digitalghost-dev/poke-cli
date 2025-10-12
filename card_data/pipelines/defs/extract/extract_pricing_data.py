from typing import Optional

import dagster as dg
import polars as pl
import requests
from pydantic import BaseModel, ValidationError
from termcolor import colored


SET_PRODUCT_MATCHING = {
    "sv01": "22873",
    "sv02": "23120",
}


class CardPricing(BaseModel):
    product_id: int
    name: str
    card_number: str
    market_price: Optional[float] = None


def is_card(item: dict) -> bool:
    """Check if item has a 'Number' field in extendedData"""
    return any(
        data_field.get("name") == "Number"
        for data_field in item.get("extendedData", [])
    )


def get_card_number(card: dict) -> Optional[str]:
    """Get the card number from extendedData"""
    for data_field in card.get("extendedData", []):
        if data_field.get("name") == "Number":
            return data_field.get("value")
    return None


def extract_card_name(full_name: str) -> str:
    """Extract clean card name, removing variant information after dash"""
    return full_name.partition("-")[0].strip() if "-" in full_name else full_name


def pull_product_information(set_number: str) -> pl.DataFrame:
    """Pull product and pricing information for a given set number."""

    print(colored(" →", "blue"), f"Processing set: {set_number}")

    product_id = SET_PRODUCT_MATCHING[set_number]

    # Fetch product data
    products_url = (f"https://tcgcsv.com/tcgplayer/3/{product_id}/products")
    products_data = requests.get(products_url, timeout=30).json()

    # Fetch pricing data
    prices_url = (f"https://tcgcsv.com/tcgplayer/3/{product_id}/prices")
    prices_data = requests.get(prices_url, timeout=30).json()

    price_dict = {
        price["productId"]: price.get("marketPrice")
        for price in prices_data.get("results", [])
    }

    cards_data = []
    for card in products_data.get("results", []):
        if not is_card(card):
            continue

        card_info = {
            "product_id": card["productId"],
            "name": extract_card_name(card["name"]),
            "card_number": get_card_number(card),
            "market_price": price_dict.get(card["productId"]),
        }
        cards_data.append(card_info)

    # Pydantic validation
    try:
        validated: list[CardPricing] = [CardPricing(**card) for card in cards_data]
        print(
            colored(" ✓", "green"),
            f"Pydantic validation passed for {len(validated)} cards.",
        )
    except ValidationError as e:
        print(colored(" ✖", "red"), "Pydantic validation failed.")
        print(e)
        raise

    df_data = [card.model_dump(mode="json") for card in validated]
    return pl.DataFrame(df_data)


@dg.asset(kinds={"API", "Polars", "Pydantic"}, name="build_pricing_dataframe")
def build_dataframe() -> pl.DataFrame:
    all_cards = []
    for set_number in SET_PRODUCT_MATCHING.keys():
        df = pull_product_information(set_number)

        # Raise error if any DataFrame is empty
        if df is None or df.shape[1] == 0 or df.is_empty():
            error_msg = f"Empty DataFrame returned for set '{set_number}'. " \
                       f"Cannot proceed with drop+replace operation to avoid data loss."
            print(colored(" ✖", "red"), error_msg)
            raise ValueError(error_msg)

        all_cards.append(df)

    concatenated = pl.concat(all_cards)
    print(concatenated)
    return concatenated
