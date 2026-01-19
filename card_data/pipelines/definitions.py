from pathlib import Path

from dagster import definitions, load_from_defs_folder

import dagster as dg

from .defs.extract.tcgcsv.extract_pricing import build_dataframe
from .defs.extract.tcgdex.extract_sets import extract_sets_data
from .defs.extract.tcgdex.extract_series import extract_series_data
from .defs.load.tcgcsv.load_pricing import load_pricing_data, data_quality_checks_on_pricing
from .defs.load.tcgdex.load_sets import load_sets_data, data_quality_check_on_sets
from .defs.load.tcgdex.load_series import load_series_data, data_quality_check_on_series
from .sensors import discord_success_sensor, discord_failure_sensor


@definitions
def defs() -> dg.Definitions:
    # load_from_defs_folder discovers dbt assets from transform_data.py
    folder_defs: dg.Definitions = load_from_defs_folder(project_root=Path(__file__).parent.parent)
    return dg.Definitions.merge(folder_defs, defs_pricing, defs_sets, defs_series)

# Pricing pipeline
pricing_pipeline_job = dg.define_asset_job(
    name="pricing_pipeline_job",
    selection=dg.AssetSelection.assets(build_dataframe).downstream(include_self=True),
)

price_schedule: dg.ScheduleDefinition = dg.ScheduleDefinition(
    name="price_schedule",
    cron_schedule="0 14 * * *",
    target=pricing_pipeline_job,
    execution_timezone="America/Los_Angeles",
)

defs_pricing: dg.Definitions = dg.Definitions(
    assets=[build_dataframe, load_pricing_data, data_quality_checks_on_pricing],
    jobs=[pricing_pipeline_job],
    schedules=[price_schedule],
    sensors=[discord_success_sensor, discord_failure_sensor]
)

# Series pipeline
series_pipeline_job = dg.define_asset_job(
    name="series_pipeline_job",
    selection=dg.AssetSelection.assets(extract_series_data).downstream(include_self=True),
)

defs_series: dg.Definitions = dg.Definitions(
    assets=[extract_series_data, load_series_data, data_quality_check_on_series],
    jobs=[series_pipeline_job],
    sensors=[discord_success_sensor, discord_failure_sensor]
)

# Sets pipeline
sets_pipeline_job = dg.define_asset_job(
    name="sets_pipeline_job",
    selection=dg.AssetSelection.assets(extract_sets_data).downstream(include_self=True),
)

defs_sets: dg.Definitions = dg.Definitions(
    assets=[extract_sets_data, load_sets_data, data_quality_check_on_sets],
    jobs=[sets_pipeline_job],
    sensors=[discord_success_sensor, discord_failure_sensor]
)