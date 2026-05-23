-- depends_on: {{ ref('comp_players') }}
{{ config(
    materialized='incremental',
    unique_key='id',
    on_schema_change='sync_all_columns',
    post_hook=[
        "ALTER TABLE {{ this }} DROP CONSTRAINT IF EXISTS comp_rounds_pkey",
        "ALTER TABLE {{ this }} DROP CONSTRAINT IF EXISTS pk_comp_rounds",
        "ALTER TABLE {{ this }} ADD CONSTRAINT pk_comp_rounds PRIMARY KEY (id)",
        "ALTER TABLE {{ this }} DROP CONSTRAINT IF EXISTS fk_comp_rounds_players",
        "ALTER TABLE {{ this }} ADD CONSTRAINT fk_comp_rounds_players FOREIGN KEY (pokedata_id, game_type, player_name) REFERENCES {{ ref('comp_players') }} (pokedata_id, game_type, player_name) ON DELETE CASCADE",
        "{{ enable_rls() }}"
    ]
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
