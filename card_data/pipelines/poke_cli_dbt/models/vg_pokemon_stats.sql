-- depends_on: {{ ref('pokemon') }}
-- depends_on: {{ ref('vg_stats') }}
{{ config(
    materialized='incremental',
    unique_key=['pokemon_id', 'stat_id'],
    incremental_strategy='merge',
    post_hook="{{ enable_rls(role='authenticated', policy_name='Enable Read Access for Authenticated Users') }}"
) }}

SELECT
    pokemon_id,
    stat_id,
    base_stat
FROM {{ source('staging', 'vg_pokemon_stats') }}
