-- depends_on: {{ ref('comp_events') }}
{{ config(
    materialized='incremental',
    unique_key='id',
    on_schema_change='sync_all_columns',
    post_hook=[
        "ALTER TABLE {{ this }} DROP CONSTRAINT IF EXISTS comp_players_pkey",
        "ALTER TABLE {{ this }} DROP CONSTRAINT IF EXISTS pk_comp_players",
        "ALTER TABLE {{ this }} ADD CONSTRAINT pk_comp_players PRIMARY KEY (id)",
        "ALTER TABLE {{ this }} DROP CONSTRAINT IF EXISTS uq_comp_players_natural",
        "ALTER TABLE {{ this }} ADD CONSTRAINT uq_comp_players_natural UNIQUE (pokedata_id, game_type, player_name)",
        "ALTER TABLE {{ this }} DROP CONSTRAINT IF EXISTS fk_comp_players_events",
        "ALTER TABLE {{ this }} ADD CONSTRAINT fk_comp_players_events FOREIGN KEY (pokedata_id, game_type) REFERENCES {{ ref('comp_events') }} (pokedata_id, game_type) ON DELETE CASCADE",
        "{{ enable_rls() }}"
    ]
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
