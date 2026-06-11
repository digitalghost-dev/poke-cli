import dagster as dg
import polars as pl
from dagster import RetryPolicy, Backoff
from sqlalchemy.exc import OperationalError
from termcolor import colored

from ...utils.secret_retriever import fetch_secret

RESOURCE_NAME: str = "types"


def create_dataframe(url: str) -> pl.DataFrame:
    df = pl.read_csv(url)

    return df


@dg.asset(
    kinds={"Supabase", "Postgres"},
    retry_policy=RetryPolicy(max_retries=3, delay=2, backoff=Backoff.EXPONENTIAL),
)
def load_vg_types() -> None:
    database_url: str = fetch_secret()
    table_name: str = f"staging.vg_{RESOURCE_NAME}"
    df = create_dataframe(
        f"https://raw.githubusercontent.com/PokeAPI/pokeapi/master/data/v2/csv/{RESOURCE_NAME}.csv"
    )

    try:
        df.write_database(
            table_name=table_name, connection=database_url, if_table_exists="replace"
        )
        print(colored(" ✓", "green"), f"Data loaded into {table_name}")
    except OperationalError as e:
        print(colored(" ✖", "red"), "Connection error in load_vg_types():", e)
        raise
