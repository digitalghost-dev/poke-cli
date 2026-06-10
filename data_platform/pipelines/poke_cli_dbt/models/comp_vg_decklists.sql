-- depends_on: {{ ref('comp_players') }}
{{ config(
    materialized='incremental',
    unique_key='id',
    on_schema_change='sync_all_columns'
) }}

SELECT
    id,
    pokedata_id,
    game_type,
    player_name,
    decklist
FROM {{ source('staging', 'comp_vg_decklists') }}
{% if is_incremental() %}
WHERE id > (SELECT COALESCE(MAX(id), 0) FROM {{ this }})
{% endif %}
