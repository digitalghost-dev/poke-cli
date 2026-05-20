import subprocess  # nosec
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
def load_tcg_tournaments_data(create_tcg_tournaments_dataframe: pl.DataFrame) -> None:
    database_url: str = fetch_secret()
    table_name: str = "staging.comp_tcg_tournaments"

    try:
        create_tcg_tournaments_dataframe.write_database(
            table_name=table_name, connection=database_url, if_table_exists="replace"
        )
        print(f" ✓ Data loaded into {table_name}")
    except OperationalError as e:
        print(f" ✖ Connection error in load_tcg_tournaments_data(): {e}")
        raise


@dg.asset(
    kinds={"Supabase", "Postgres"},
    retry_policy=RetryPolicy(max_retries=3, delay=2, backoff=Backoff.EXPONENTIAL),
)
def load_vg_tournaments_data(create_vg_tournaments_dataframe: pl.DataFrame) -> None:
    database_url: str = fetch_secret()
    table_name: str = "staging.comp_vg_tournaments"

    try:
        create_vg_tournaments_dataframe.write_database(
            table_name=table_name, connection=database_url, if_table_exists="replace"
        )
        print(f" ✓ Data loaded into {table_name}")
    except OperationalError as e:
        print(f" ✖ Connection error in load_vg_tournaments_data(): {e}")
        raise


@dg.asset(
    deps=[load_tcg_tournaments_data],
    kinds={"Soda"},
    name="data_quality_checks_on_comp_tcg_tournaments",
)
def data_quality_checks_on_comp_tcg_tournaments() -> None:
    current_file_dir = Path(__file__).parent
    result = subprocess.run(  # nosec
        [
            "soda",
            "scan",
            "-d",
            "supabase",
            "-c",
            "../../../soda/configuration.yml",
            "../../../soda/checks_comp_tcg_tournaments.yml",
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


@dg.asset(
    deps=[load_vg_tournaments_data],
    kinds={"Soda"},
    name="data_quality_checks_on_comp_vg_tournaments",
)
def data_quality_checks_on_comp_vg_tournaments() -> None:
    current_file_dir = Path(__file__).parent
    result = subprocess.run(  # nosec
        [
            "soda",
            "scan",
            "-d",
            "supabase",
            "-c",
            "../../../soda/configuration.yml",
            "../../../soda/checks_comp_vg_tournaments.yml",
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
