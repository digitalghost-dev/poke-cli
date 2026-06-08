import polars as pl
from sqlalchemy.exc import OperationalError

from ...utils.secret_retriever import fetch_secret

RESOURCE_NAME: str = "types"


def create_dataframe(url: str) -> pl.DataFrame:
    df = pl.read_csv(url)

    return df


def upload_dataframe():
    database_url: str = fetch_secret()
    table_name: str = f"staging.vg_{RESOURCE_NAME}"
    df = create_dataframe(f"https://raw.githubusercontent.com/PokeAPI/pokeapi/master/data/v2/csv/{RESOURCE_NAME}.csv")

    try:
        df.write_database(
            table_name=table_name, connection=database_url, if_table_exists="replace"
        )
        print(f" ✓ Data loaded into {table_name}")
    except OperationalError as e:
        print(f" ✖ Connection error in load_pricing_data():", e)
        raise


if __name__ == "__main__":
    upload_dataframe()
