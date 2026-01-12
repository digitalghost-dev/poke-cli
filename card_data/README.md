# Card Data

This directory stores all the code for all backend data processing related to Pokémon TCG data.

Instead of calling directly to the PokéAPI for data from the video game, I took this a step further
and decided to process all the data myself, load it into Supabase, and read from that API.

## Data Architecture
Runs at 2:00PM PST daily.
![data_diagram](https://poke-cli-s3-bucket.s3.us-west-2.amazonaws.com/data_infrastructure_v2.png)


1. TCGPlayer pricing data and TCGDex card data are called and processed through a data pipeline orchestrated by Dagster and hosted on AWS.

2. When the pipeline starts, Pydantic validates the incoming API data against a pre-defined schema, ensuring the data types match the expected structure.

3. Polars is used to create DataFrames.

4. The data is loaded into a Supabase staging schema.

5. Soda data quality checks are performed.

6. `dbt` runs and builds the final tables in a Supabase production schema.

7. Users are then able to query the `pokeapi.co` or supabase APIs for either video game or trading card data, respectively.
