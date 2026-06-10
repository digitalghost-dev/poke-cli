import shutil
import subprocess  # nosec
from pathlib import Path

import dagster as dg
import polars as pl
from dagster import Backoff, RetryPolicy
from sqlalchemy.exc import OperationalError
from termcolor import colored

from ....utils.secret_retriever import fetch_secret

SODA_PATH = shutil.which("soda") or "soda"


@dg.asset(
    kinds={"Supabase", "Postgres"},
    name="load_series_data",
    retry_policy=RetryPolicy(max_retries=3, delay=2, backoff=Backoff.EXPONENTIAL),
)
def load_series_data(extract_series_data: pl.DataFrame) -> None:
    database_url: str = fetch_secret()
    table_name: str = "staging.series"

    df = extract_series_data
    try:
        df.write_database(
            table_name=table_name, connection=database_url, if_table_exists="replace"
        )
        print(colored(" ✓", "green"), f"Data loaded into {table_name}")
    except OperationalError as e:
        print(colored(" ✖", "red"), "Connection error in load_series_data():", e)
        raise


@dg.asset(deps=[load_series_data], kinds={"Soda"}, name="data_quality_checks_on_series")
def data_quality_check_on_series() -> None:
    current_file_dir = Path(__file__).parent
    print(f"Setting cwd to: {current_file_dir}")

    result = subprocess.run(  # nosec B603
        [
            SODA_PATH,
            "scan",
            "-d",
            "supabase",
            "-c",
            "../../../soda/configuration.yml",
            "../../../soda/checks_series.yml",
        ],
        capture_output=True,
        text=True,
        cwd=current_file_dir,
    )

    # Print output for debugging
    if result.stdout:
        print(result.stdout)
    if result.stderr:
        print(result.stderr)

    if result.returncode != 0:
        raise Exception(
            f"Soda data quality checks failed with return code {result.returncode}"
        )
