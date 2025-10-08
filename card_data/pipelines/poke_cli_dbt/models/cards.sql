{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

SELECT id, set_id, image, name, "localId", category, hp
FROM {{ source('staging', 'cards') }}