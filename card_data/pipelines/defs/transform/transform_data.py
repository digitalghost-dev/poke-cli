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
                "series": "data_quality_checks_on_series",
                "sets": "data_quality_checks_on_sets",
                "cards": "load_card_data",
                "pricing_data": "data_quality_checks_on_pricing",
                "standings": "load_standings_data",
                "comp_events": "data_quality_checks_on_comp_events",
                "comp_players": "data_quality_checks_on_comp_players",
                "comp_rounds": "data_quality_checks_on_comp_rounds",
                "comp_vg_decklists": "data_quality_checks_on_comp_vg_decklists",
                "comp_tcg_decklists": "data_quality_checks_on_comp_tcg_decklists",
                "pokemon": "load_pokemon",
                "vg_types": "load_vg_types",
                "vg_stats": "load_vg_stats",
                "vg_pokemon_types": "load_vg_pokemon_types",
                "vg_pokemon_stats": "load_vg_pokemon_stats",
                "pikalytics_speed_tiers": "trigger_pikalytics_speed_tiers",
            }
            if name in source_mapping:
                return dg.AssetKey([source_mapping[name]])

        # For models, use default behavior
        return super().get_asset_key(dbt_resource_props)


@dbt_assets(
    manifest=DBT_PROJECT_PATH / "target" / "manifest.json",
    dagster_dbt_translator=CustomDbtTranslator(),
)
def dbt_build_assets(context: dg.AssetExecutionContext, dbt: DbtCliResource):
    """
    dbt assets that transform staging data into final models.
    """
    yield from dbt.cli(["build"], context=context).stream()


dbt_resource = DbtCliResource(project_dir=DBT_PROJECT_PATH)
defs = dg.Definitions(assets=[dbt_build_assets], resources={"dbt": dbt_resource})
