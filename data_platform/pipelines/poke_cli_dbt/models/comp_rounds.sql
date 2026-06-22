-- depends_on: {{ ref('comp_players') }}
{{ config(
    materialized='incremental',
    unique_key=['pokedata_id', 'game_type', 'player_name', 'round_number'],
    incremental_strategy='merge',
    on_schema_change='append_new_columns'
) }}

SELECT
    pokedata_id,
    game_type,
    player_name,
    round_number,
    opponent_name,
    result,
    NULLIF(TRIM(table_number), '')::int AS table_number
FROM {{ source('staging', 'comp_rounds') }}
