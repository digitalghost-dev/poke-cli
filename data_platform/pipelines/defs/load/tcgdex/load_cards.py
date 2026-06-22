import dagster as dg
import polars as pl
from ....utils.secret_retriever import fetch_secret
from dagster import Backoff, RetryPolicy
from sqlalchemy.exc import OperationalError
from termcolor import colored


@dg.asset(
    kinds={"Supabase", "Postgres"},
    name="load_card_data",
    retry_policy=RetryPolicy(max_retries=3, delay=2, backoff=Backoff.EXPONENTIAL),
)
def load_card_data(create_card_dataframe: pl.DataFrame) -> None:
    database_url: str = fetch_secret()
    table_name: str = "staging.cards"

    df = create_card_dataframe
    try:
        df.write_database(
            table_name=table_name, connection=database_url, if_table_exists="append"
        )
        print(colored(" ✓", "green"), f"Data loaded into {table_name}")
    except OperationalError as e:
        print(colored(" ✖", "red"), "Connection error in load_card_data():", e)
        raise
