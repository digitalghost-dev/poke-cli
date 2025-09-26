{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

SELECT *
FROM {{ source('staging', 'sets') }}