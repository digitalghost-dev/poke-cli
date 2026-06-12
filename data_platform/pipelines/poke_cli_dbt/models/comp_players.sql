-- depends_on: {{ ref('comp_events') }}
{{ config(
    materialized='incremental',
    unique_key=['pokedata_id', 'game_type', 'player_name'],
    incremental_strategy='merge',
    on_schema_change='append_new_columns'
) }}

SELECT
    pokedata_id,
    game_type,
    player_name,
    country,
    placement,
    wins,
    losses,
    ties,
    resistance_self,
    resistance_opp,
    resistance_oppopp,
    dropped_round,
    trainer_name
FROM {{ source('staging', 'comp_players') }}
