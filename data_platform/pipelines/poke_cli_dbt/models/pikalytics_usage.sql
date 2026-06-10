{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

WITH staged AS (
    SELECT
        rank,
        pokemon,
        usage_percent,
        web_url,
        TRIM(BOTH '-' FROM LOWER(REGEXP_REPLACE(pokemon, '[^a-zA-Z0-9]+', '-', 'g'))) AS pokemon_slug
    FROM {{ source('staging', 'pikalytics_usage') }}
)

SELECT
    'gen9championsvgc2026regma'              AS format,
    rank,
    pokemon,
    pokemon_slug,
    {{ resolve_pokemon_id('pokemon_slug') }} AS pokemon_id,
    usage_percent,
    web_url,
    'pikalytics'                             AS source
FROM staged
ORDER BY rank
