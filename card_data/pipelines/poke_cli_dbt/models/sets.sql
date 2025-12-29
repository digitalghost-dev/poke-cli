{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

SELECT *
FROM {{ source('staging', 'sets') }}
WHERE set_id NOT IN ('fut2020', 'mep', 'svp', 'swshp')