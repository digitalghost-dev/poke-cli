import dagster as dg
from dagster import RetryPolicy, Backoff
from sqlalchemy.exc import OperationalError
from ..extract.extract_pricing_data import build_dataframe
from ...utils.secret_retriever import fetch_secret
from termcolor import colored


@dg.asset(
    deps=[build_dataframe],
    kinds={"Supabase", "Postgres"},
    retry_policy=RetryPolicy(max_retries=3, delay=2, backoff=Backoff.EXPONENTIAL),
)
def load_pricing_data() -> None:
    database_url: str = fetch_secret()
    table_name: str = "staging.pricing_data"

    df = build_dataframe()
    try:
        df.write_database(
            table_name=table_name, connection=database_url, if_table_exists="replace"
        )
        print(colored(" ✓", "green"), f"Data loaded into {table_name}")
    except OperationalError as e:
        print(colored(" ✖", "red"), "Connection error in load_series_data():", e)
        raise
