{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

SELECT id, set_id, image, name, "localId", category, hp, "set_cardCount_official", set_name, illustrator, "regulationMark" AS regulation_mark
FROM {{ source('staging', 'cards') }}
WHERE "localId" ~ '^[0-9]+$'