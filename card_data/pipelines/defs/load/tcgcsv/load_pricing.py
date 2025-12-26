import subprocess # nosec
from pathlib import Path

import dagster as dg
import polars as pl
from dagster import RetryPolicy, Backoff
from sqlalchemy.exc import OperationalError
from termcolor import colored

from ....utils.secret_retriever import fetch_secret


@dg.asset(
    kinds={"Supabase", "Postgres"},
    retry_policy=RetryPolicy(max_retries=3, delay=2, backoff=Backoff.EXPONENTIAL),
)
def load_pricing_data(build_pricing_dataframe: pl.DataFrame) -> None:
    database_url: str = fetch_secret()
    table_name: str = "staging.pricing_data"

    try:
        build_pricing_dataframe.write_database(
            table_name=table_name, connection=database_url, if_table_exists="replace"
        )
        print(colored(" ✓", "green"), f"Data loaded into {table_name}")
    except OperationalError as e:
        print(colored(" ✖", "red"), "Connection error in load_pricing_data():", e)
        raise


@dg.asset(
    deps=[load_pricing_data],
    kinds={"Soda"},
    name="data_quality_checks_on_pricing",
)
def data_quality_checks_on_pricing() -> None:
    current_file_dir = Path(__file__).parent
    print(f"Setting cwd to: {current_file_dir}")

    result = subprocess.run( # nosec
        [
            "soda",
            "scan",
            "-d",
            "supabase",
            "-c",
            "../../../soda/configuration.yml",
            "../../../soda/checks_pricing.yml",
        ],
        capture_output=True,
        text=True,
        cwd=current_file_dir,
    )

    if result.stdout:
        print(result.stdout)
    if result.stderr:
        print(result.stderr)

    if result.returncode != 0:
        raise Exception(f"Soda data quality checks failed with return code {result.returncode}")
