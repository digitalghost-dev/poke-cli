{{ config(
    materialized='incremental',
    unique_key='id',
    on_schema_change='append_new_columns'
) }}

SELECT
    id,
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
{% if is_incremental() %}
WHERE last_updated::timestamp > (
    SELECT COALESCE(MAX(last_updated), '1900-01-01'::timestamp) FROM {{ this }}
)
{% endif %}
