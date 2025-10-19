from pathlib import Path

from dagster import definitions, load_from_defs_folder
from dagster_dbt import DbtCliResource, DbtProject, dbt_assets

import dagster as dg

from .defs.extract.extract_pricing_data import build_dataframe
from .defs.load.load_pricing_data import load_pricing_data


@definitions
def defs():
    folder_defs = load_from_defs_folder(project_root=Path(__file__).parent.parent)
    return dg.Definitions.merge(folder_defs, defs_pricing)

dbt_project_directory = Path(__file__).absolute().parent / "poke_cli_dbt"
dbt_project = DbtProject(project_dir=dbt_project_directory)

dbt_resource = DbtCliResource(project_dir=dbt_project)

# Compiles the dbt project & allow Dagster to build an asset graph
dbt_project.prepare_if_dev()

# Yields Dagster events streamed from the dbt CLI
@dbt_assets(manifest=dbt_project.manifest_path)
def dbt_models(context: dg.AssetExecutionContext, dbt: DbtCliResource):
    yield from dbt.cli(["build"], context=context).stream()

price_schedule = dg.ScheduleDefinition(
    name="price_schedule",
    cron_schedule="10 10 * * *",
    target=[build_dataframe, load_pricing_data],
    execution_timezone="America/Los_Angeles",
)

defs_pricing = dg.Definitions(
    schedules=[price_schedule],
)