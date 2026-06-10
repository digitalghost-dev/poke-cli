from pathlib import Path

from dagster import definitions, load_from_defs_folder

import dagster as dg

from .defs.extract.limitless.extract_standings import create_standings_dataframe
from .defs.competitive import (
    create_events_dataframe,
    load_events_data,
    data_quality_checks_on_comp_events,
    load_players_data,
    data_quality_checks_on_comp_players,
    data_quality_checks_on_comp_rounds,
    data_quality_checks_on_comp_vg_decklists,
    data_quality_checks_on_comp_tcg_decklists,
)
from .defs.extract.tcgcsv.extract_pricing import build_dataframe
from .defs.extract.tcgdex.extract_sets import extract_sets_data
from .defs.extract.tcgdex.extract_series import extract_series_data
from .defs.load.limitless.load_standings import load_standings_data
from .defs.load.tcgcsv.load_pricing import (
    load_pricing_data,
    data_quality_checks_on_pricing,
)
from .defs.load.tcgdex.load_sets import load_sets_data, data_quality_check_on_sets
from .defs.load.tcgdex.load_series import load_series_data, data_quality_check_on_series
from .defs.pokeapi.pokemon import load_pokemon
from .defs.pokeapi.types import load_vg_types
from .defs.pokeapi.stats import load_vg_stats
from .defs.pokeapi.pokemon_types import load_vg_pokemon_types
from .defs.pokeapi.pokemon_stats import load_vg_pokemon_stats
from .defs.pikalytics.speed_tiers import trigger_pikalytics_speed_tiers
from .defs.pikalytics.usage import trigger_pikalytics_usage
from .defs.pikalytics.top_teams import trigger_pikalytics_top_teams
from .sensors import discord_success_sensor, discord_failure_sensor


@definitions
def defs() -> dg.Definitions:
    # load_from_defs_folder discovers dbt assets from transform_data.py
    folder_defs: dg.Definitions = load_from_defs_folder(
        project_root=Path(__file__).parent.parent
    )
    return dg.Definitions.merge(
        folder_defs,
        defs_discord_sensors,
        defs_pricing,
        defs_sets,
        defs_series,
        defs_standings,
        defs_events,
        defs_pokeapi,
        defs_pikalytics,
    )


defs_discord_sensors: dg.Definitions = dg.Definitions(
    sensors=[discord_success_sensor, discord_failure_sensor],
)

# Pricing pipeline job
pricing_pipeline = dg.define_asset_job(
    name="pricing_pipeline_job",
    selection=dg.AssetSelection.assets(build_dataframe).downstream(include_self=True),
)

price_schedule: dg.ScheduleDefinition = dg.ScheduleDefinition(
    name="price_schedule",
    cron_schedule="0 14 * * *",
    target=pricing_pipeline,
    execution_timezone="America/Los_Angeles",
)

defs_pricing: dg.Definitions = dg.Definitions(
    assets=[build_dataframe, load_pricing_data, data_quality_checks_on_pricing],
    jobs=[pricing_pipeline],
    schedules=[price_schedule],
)

# Series pipeline job
series_pipeline = dg.define_asset_job(
    name="series_pipeline_job",
    selection=dg.AssetSelection.assets(extract_series_data).downstream(
        include_self=True
    ),
)

defs_series: dg.Definitions = dg.Definitions(
    assets=[extract_series_data, load_series_data, data_quality_check_on_series],
    jobs=[series_pipeline],
)

# Sets pipeline job
sets_pipeline = dg.define_asset_job(
    name="sets_pipeline_job",
    selection=dg.AssetSelection.assets(extract_sets_data).downstream(include_self=True),
)

defs_sets: dg.Definitions = dg.Definitions(
    assets=[extract_sets_data, load_sets_data, data_quality_check_on_sets],
    jobs=[sets_pipeline],
)

# Standings pipeline job
standings_pipeline = dg.define_asset_job(
    name="standings_pipeline_job",
    selection=dg.AssetSelection.assets(create_standings_dataframe).downstream(
        include_self=True
    ),
)

defs_standings: dg.Definitions = dg.Definitions(
    assets=[create_standings_dataframe, load_standings_data],
    jobs=[standings_pipeline],
)

# Competitive events + tournaments pipeline job (single job, branches into TCG and VG)
events_pipeline = dg.define_asset_job(
    name="comp_pipeline_job",
    selection=dg.AssetSelection.assets(create_events_dataframe).downstream(include_self=True),
)

defs_events: dg.Definitions = dg.Definitions(
    assets=[
        create_events_dataframe,
        load_events_data,
        data_quality_checks_on_comp_events,
        load_players_data,
        data_quality_checks_on_comp_players,
        data_quality_checks_on_comp_rounds,
        data_quality_checks_on_comp_vg_decklists,
        data_quality_checks_on_comp_tcg_decklists,
    ],
    jobs=[events_pipeline],
)


# PokéAPI video-game data pipeline job (5 staging loads + their downstream dbt models)
pokeapi_pipeline = dg.define_asset_job(
    name="pokeapi_pipeline_job",
    selection=dg.AssetSelection.assets(
        load_pokemon,
        load_vg_types,
        load_vg_stats,
        load_vg_pokemon_types,
        load_vg_pokemon_stats,
    ).downstream(include_self=True),
)

# Runs on the 1st and 15th of each month at 14:00 LA time
pokeapi_schedule: dg.ScheduleDefinition = dg.ScheduleDefinition(
    name="pokeapi_schedule",
    cron_schedule="0 14 1,15 * *",
    target=pokeapi_pipeline,
    execution_timezone="America/Los_Angeles",
)

defs_pokeapi: dg.Definitions = dg.Definitions(
    assets=[
        load_pokemon,
        load_vg_types,
        load_vg_stats,
        load_vg_pokemon_types,
        load_vg_pokemon_stats,
    ],
    jobs=[pokeapi_pipeline],
    schedules=[pokeapi_schedule],
)


# Pikalytics pipeline job (Dagster-first: triggers the n8n scrape, then dbt builds public).
# Scaffolded for fan-out — add the other pikalytics trigger assets to the selection as they migrate.
pikalytics_pipeline = dg.define_asset_job(
    name="pikalytics_pipeline_job",
    selection=dg.AssetSelection.assets(
        trigger_pikalytics_speed_tiers,
        trigger_pikalytics_usage,
        trigger_pikalytics_top_teams,
    ).downstream(include_self=True),
)

# Weekly, Mondays at 08:00 LA time (mirrors the old n8n speed-tiers cadence)
pikalytics_schedule: dg.ScheduleDefinition = dg.ScheduleDefinition(
    name="pikalytics_schedule",
    cron_schedule="0 8 * * 1",
    target=pikalytics_pipeline,
    execution_timezone="America/Los_Angeles",
)

defs_pikalytics: dg.Definitions = dg.Definitions(
    assets=[
        trigger_pikalytics_speed_tiers,
        trigger_pikalytics_usage,
        trigger_pikalytics_top_teams,
    ],
    jobs=[pikalytics_pipeline],
    schedules=[pikalytics_schedule],
)
