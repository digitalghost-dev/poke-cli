import requests
import polars as pl

pl.Config(tbl_rows=-1)

SET_PRODUCT_MATCHING = {
    "sv01": "22873"
}

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

def pull_product_information():
    url = f"https://tcgcsv.com/tcgplayer/3/{SET_PRODUCT_MATCHING['sv01']}/products"
    r = requests.get(url)

    if r.status_code != 200:
        return

    data = r.json()

    url_prices = f"https://tcgcsv.com/tcgplayer/3/22873/prices"
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

    print(df.sort("card_number"))

pull_product_information()
