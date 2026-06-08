{{ config(
    materialized='table',
    post_hook="{{ enable_rls(role='authenticated', policy_name='Enable Read Access for Authenticated Users') }}"
) }}

SELECT
    pokemon_id,
    type_id,
    slot
FROM {{ source('staging', 'vg_pokemon_types') }}
