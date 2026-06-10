-- depends_on: {{ ref('comp_players') }}
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
    round_number,
    opponent_name,
    result,
    NULLIF(TRIM(table_number), '')::int AS table_number
FROM {{ source('staging', 'comp_rounds') }}
{% if is_incremental() %}
WHERE id > (SELECT COALESCE(MAX(id), 0) FROM {{ this }})
{% endif %}
