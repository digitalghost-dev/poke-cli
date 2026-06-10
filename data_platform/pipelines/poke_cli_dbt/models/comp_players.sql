-- depends_on: {{ ref('comp_events') }}
{{ config(
    materialized='incremental',
    unique_key='id',
    on_schema_change='append_new_columns'
) }}

SELECT
    id,
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
{% if is_incremental() %}
WHERE id > (SELECT COALESCE(MAX(id), 0) FROM {{ this }})
{% endif %}
