import dagster as dg
import polars as pl
from dagster import Backoff, RetryPolicy
from sqlalchemy.exc import OperationalError
from termcolor import colored

from ....utils.secret_retriever import fetch_secret


@dg.asset(
    kinds={"Supabase", "Postgres"},
    name="load_standings_data",
    retry_policy=RetryPolicy(max_retries=3, delay=2, backoff=Backoff.EXPONENTIAL),
)
def load_standings_data(create_standings_dataframe: pl.DataFrame) -> None:
    database_url: str = fetch_secret()
    table_name: str = "staging.standings"

    df = create_standings_dataframe
    try:
        df.write_database(
            table_name=table_name, connection=database_url, if_table_exists="replace"
        )
        print(colored(" ✓", "green"), f"Data loaded into {table_name}")
    except OperationalError as e:
        print(colored(" ✖", "red"), "Connection error in load_card_data():", e)
        raise
