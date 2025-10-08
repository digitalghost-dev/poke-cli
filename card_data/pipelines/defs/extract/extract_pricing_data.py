import time
import json

import dagster as dg
import polars as pl

from pydantic import BaseModel, HttpUrl, ValidationError
from typing import Optional
from termcolor import colored

import requests

SET_PRODUCT_MATCHING = {
    "sv01": "22873",
    "sv02": "23120",
}

class CardPricing(BaseModel):
    product_id: int
    name: str
    card_number: str
    market_price: float

def is_card(item):
    """Check if item has a 'Number' field in extendedData"""
    for data_field in item["extendedData"]:
        if data_field["name"] == "Number":
            return True
    return False

def get_card_number(card):
    """Get the card number from extendedData"""
    for data_field in card["extendedData"]:
        if data_field["name"] == "Number":
            return data_field["value"]
    return None

def pull_product_information(set_number: str) -> pl.DataFrame | str:
    url = f"https://tcgcsv.com/tcgplayer/3/{SET_PRODUCT_MATCHING[set_number]}/products"
    data = requests.get(url).json()
    r = requests.get(url)

    try:
        validated: list[CardPricing] = [CardPricing(**item) for item in data]
        print(
            colored(" ✓", "green"), "Pydantic validation passed for all series entries."
        )
        if r.status_code == 200:
            print(
                colored(" ✓", "green"), "Successful connection to API."
            )

            url_prices = f"https://tcgcsv.com/tcgplayer/3/{SET_PRODUCT_MATCHING[set_number]}/prices"
            r_prices = requests.get(url_prices)
            price_data = r_prices.json()

            price_dict = {price["productId"]: price["marketPrice"]
                        for price in price_data["results"]}

            product_id_list = []
            name_list = []
            card_number_list = []
            price_list = []

            for card in data["results"]:
                if not is_card(card):
                    continue

                number = get_card_number(card)
                card_number_list.append(number)

                name = card["name"].partition("-")[0].strip() if "-" in card["name"] else card["name"]
                name_list.append(name)

                product_id = card["productId"]
                product_id_list.append(product_id)

                market_price = price_dict.get(product_id)
                price_list.append(market_price)

            df = pl.DataFrame({
                "product_id": product_id_list,
                "name": name_list,
                "card_number": card_number_list,
                "market_price": price_list,
            }).with_columns(pl.col("market_price").cast(pl.Decimal(scale=2)))

            return df

        else:
            return str(colored(" ✖", "red"), f"Connection error to API: {r.status_code}")

    except ValidationError as e:
        print(colored(" ✖", "red"), "Pydantic validation failed.")
        print(e)
        raise
