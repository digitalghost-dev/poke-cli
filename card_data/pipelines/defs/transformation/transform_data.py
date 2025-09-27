import dagster as dg
from dagster_dbt import DbtCliResource, dbt_assets
from pathlib import Path

DBT_PROJECT_PATH = Path(__file__).joinpath("..", "..", "..", "poke_cli_dbt").resolve()

@dbt_assets(manifest=DBT_PROJECT_PATH / "target" / "manifest.json")
def poke_cli_dbt_assets(context: dg.AssetExecutionContext, dbt: DbtCliResource):
    """
    dbt assets that transform staging data into final models.
    """
    yield from dbt.cli(["build"], context=context).stream()

dbt_resource = DbtCliResource(project_dir=DBT_PROJECT_PATH)
defs = dg.Definitions(
    assets=[poke_cli_dbt_assets],
    resources={"dbt": dbt_resource}
)
