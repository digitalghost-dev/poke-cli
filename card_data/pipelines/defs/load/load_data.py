import dagster as dg
from sqlalchemy.exc import OperationalError
from ..extract.extract_data import (
    extract_series_data,
    extract_set_data,
    create_card_dataframe,
)
from ...utils.secret_retriever import fetch_secret
from termcolor import colored
import subprocess
from pathlib import Path


@dg.asset(deps=[extract_series_data], kinds={"Supabase", "Postgres"})
def load_series_data() -> None:
    database_url: str = fetch_secret()
    table_name: str = "staging.series"

    df = extract_series_data()
    try:
        df.write_database(
            table_name=table_name, connection=database_url, if_table_exists="replace"
        )
        print(colored(" ✓", "green"), f"Data loaded into {table_name}")
    except OperationalError as e:
        print(colored(" ✖", "red"), "Error:", e)


@dg.asset(deps=[load_series_data], kinds={"Soda"})
def data_quality_check() -> None:
    # Set working directory to where this file is located
    current_file_dir = Path(__file__).parent
    print(f"Setting cwd to: {current_file_dir}")

    result = subprocess.run(
        [
            "soda",
            "scan",
            "-d",
            "supabase",
            "-c",
            "../../soda/configuration.yml",
            "../../soda/checks.yml",
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


@dg.asset(deps=[extract_set_data], kinds={"Supabase", "Postgres"})
def load_set_data() -> None:
    database_url: str = fetch_secret()
    table_name: str = "staging.sets"

    df = extract_set_data()
    try:
        df.write_database(
            table_name=table_name, connection=database_url, if_table_exists="replace"
        )
        print(colored(" ✓", "green"), f"Data loaded into {table_name}")
    except OperationalError as e:
        print(colored(" ✖", "red"), "Error:", e)


@dg.asset(deps=[create_card_dataframe], kinds={"Supabase", "Postgres"})
def load_card_data() -> None:
    database_url: str = fetch_secret()
    table_name: str = "staging.cards"

    df = create_card_dataframe()
    try:
        df.write_database(
            table_name=table_name, connection=database_url, if_table_exists="append"
        )
        print(colored(" ✓", "green"), f"Data loaded into {table_name}")
    except OperationalError as e:
        print(colored(" ✖", "red"), "Error:", e)
