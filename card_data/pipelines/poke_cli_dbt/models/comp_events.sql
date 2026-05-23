{{ config(
    materialized='incremental',
    unique_key='id',
    on_schema_change='sync_all_columns',
    post_hook=[
        "ALTER TABLE {{ this }} DROP CONSTRAINT IF EXISTS comp_events_pkey",
        "ALTER TABLE {{ this }} DROP CONSTRAINT IF EXISTS pk_comp_events",
        "ALTER TABLE {{ this }} ADD CONSTRAINT pk_comp_events PRIMARY KEY (id)",
        "ALTER TABLE {{ this }} DROP CONSTRAINT IF EXISTS uq_comp_events_pokedata_game",
        "ALTER TABLE {{ this }} ADD CONSTRAINT uq_comp_events_pokedata_game UNIQUE (pokedata_id, game_type)",
        "{{ enable_rls() }}"
    ]
) }}

SELECT
    id,
    pokedata_id,
    game_type,
    name,
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
