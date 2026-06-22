-- depends_on: {{ ref('comp_players') }}
{{ config(
    materialized='incremental',
    unique_key=['pokedata_id', 'game_type', 'player_name'],
    incremental_strategy='merge',
    on_schema_change='sync_all_columns'
) }}

SELECT
    pokedata_id,
    game_type,
    player_name,
    decklist
FROM {{ source('staging', 'comp_vg_decklists') }}
