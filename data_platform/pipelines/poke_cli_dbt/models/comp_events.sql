{{ config(
    materialized='incremental',
    unique_key=['pokedata_id', 'game_type'],
    incremental_strategy='merge',
    on_schema_change='append_new_columns'
) }}

SELECT
    pokedata_id,
    game_type,
    name,
    regexp_replace(name, '^(\d{4}\s+)?(.+?)\s+Pokémon.*$', '\2') AS location,
    start_date::date,
    end_date::date,
    season,
    count,
    rounds,
    last_updated::timestamp
FROM {{ source('staging', 'comp_events') }}
