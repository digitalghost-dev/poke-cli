{{ config(
    materialized='table',
    post_hook="{{ enable_rls(role='authenticated', policy_name='Enable Read Access for Authenticated Users') }}"
) }}

SELECT
    id,
    identifier
FROM {{ source('staging', 'vg_stats') }}
