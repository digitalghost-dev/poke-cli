import dagster as dg
from dagster_dbt import DbtCliResource, dbt_assets
from pathlib import Path

from ..load.load_data import (
    load_series_data,
    data_quality_check_on_series,
    load_set_data,
    load_card_data
)

DBT_PROJECT_PATH = Path(__file__).joinpath("..", "..", "..", "poke_cli_dbt").resolve()

@dg.asset(deps=[load_series_data, data_quality_check_on_series, load_set_data, load_card_data], kinds=["dbt"])
def dbt_transformation(context: dg.AssetExecutionContext):
    """Run dbt build after all extract and load operations complete"""
    import subprocess
    import os
    
    # Set environment variables for dbt
    env = os.environ.copy()
    env["SUPABASE_PASSWORD"] = os.getenv("SUPABASE_PASSWORD", "")
    
    # Run dbt build
    result = subprocess.run(
        ["dbt", "build"],
        cwd=str(DBT_PROJECT_PATH),
        env=env,
        capture_output=True,
        text=True
    )
    
    if result.returncode != 0:
        context.log.error(f"dbt build failed: {result.stderr}")
        raise Exception(f"dbt build failed: {result.stderr}")
    
    context.log.info(f"dbt build completed successfully: {result.stdout}")
    return "dbt build completed"

# Create definitions for this transformation
defs = dg.Definitions(
    assets=[dbt_transformation]
)
