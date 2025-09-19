from pathlib import Path

from dagster import definitions, load_from_defs_folder
from dagster_dbt import DbtCliResource, DbtProject, dbt_assets

import dagster as dg


@definitions
def defs():
    return load_from_defs_folder(project_root=Path(__file__).parent.parent)

dbt_project_directory = Path(__file__).absolute().parent / "poke_cli_dbt"
dbt_project = DbtProject(project_dir=dbt_project_directory)

dbt_resource = DbtCliResource(project_dir=dbt_project)

# Compiles the dbt project & allow Dagster to build an asset graph
dbt_project.prepare_if_dev()

# Yields Dagster events streamed from the dbt CLI
@dbt_assets(manifest=dbt_project.manifest_path)
def dbt_models(context: dg.AssetExecutionContext, dbt: DbtCliResource):
    yield from dbt.cli(["build"], context=context).stream()

# Dagster object that contains the dbt assets and resource
defs_dbt = dg.Definitions(assets=[dbt_models], resources={"dbt": dbt_resource})