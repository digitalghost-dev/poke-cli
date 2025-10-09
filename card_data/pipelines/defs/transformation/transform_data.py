import dagster as dg
from dagster_dbt import DbtCliResource, DagsterDbtTranslator, dbt_assets
from pathlib import Path

DBT_PROJECT_PATH = Path(__file__).joinpath("..", "..", "..", "poke_cli_dbt").resolve()

class CustomDbtTranslator(DagsterDbtTranslator):
    def get_asset_key(self, dbt_resource_props):

        resource_type = dbt_resource_props["resource_type"]
        name = dbt_resource_props["name"]

        if resource_type == "source":
            # Map staging sources to load assets
            source_mapping = {
                "series": "quality_checks_series",
                "sets": "load_set_data",
                "cards": "load_card_data"
            }
            if name in source_mapping:
                return dg.AssetKey([source_mapping[name]])

        # For models, use default behavior
        return super().get_asset_key(dbt_resource_props)

@dbt_assets(
    manifest=DBT_PROJECT_PATH / "target" / "manifest.json",
    dagster_dbt_translator=CustomDbtTranslator()
)
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
