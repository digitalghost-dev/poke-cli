import subprocess # nosec
from pathlib import Path

import dagster as dg
import polars as pl
from dagster import RetryPolicy, Backoff
from sqlalchemy.exc import OperationalError

from ....utils.secret_retriever import fetch_secret


@dg.asset(
    kinds={"Supabase", "Postgres"},
    retry_policy=RetryPolicy(max_retries=3, delay=2, backoff=Backoff.EXPONENTIAL),
)
def load_events_data(create_events_dataframe: pl.DataFrame) -> None:
    database_url: str = fetch_secret()
    table_name: str = "staging.comp_events"

    try:
        create_events_dataframe.write_database(
            table_name=table_name, connection=database_url, if_table_exists="replace"
        )
        print(f" ✓ Data loaded into {table_name}")
    except OperationalError as e:
        print(f" ✖ Connection error in load_events_data(): {e}")
        raise


@dg.asset(
    deps=[load_events_data],
    kinds={"Soda"},
    name="data_quality_checks_on_comp_events",
)
def data_quality_checks_on_comp_events() -> None:
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
            "../../../soda/checks_comp_events.yml",
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
