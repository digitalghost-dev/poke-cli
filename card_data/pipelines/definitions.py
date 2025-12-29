from pathlib import Path

from dagster import definitions, load_from_defs_folder

import dagster as dg

from .defs.extract.tcgcsv.extract_pricing import build_dataframe
from .defs.load.tcgcsv.load_pricing import load_pricing_data, data_quality_checks_on_pricing


@definitions
def defs() -> dg.Definitions:
    # load_from_defs_folder discovers dbt assets from transform_data.py
    folder_defs: dg.Definitions = load_from_defs_folder(project_root=Path(__file__).parent.parent)
    return dg.Definitions.merge(folder_defs, defs_pricing)

# Define the pricing pipeline job that materializes the assets and downstream dbt model
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
)