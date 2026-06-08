{{ config(
    materialized='table',
    post_hook="{{ enable_rls(role='authenticated', policy_name='Enable Read Access for Authenticated Users') }}"
) }}

SELECT
    pokemon_id,
    stat_id,
    base_stat
FROM {{ source('staging', 'vg_pokemon_stats') }}
