{{ config(
    materialized='incremental',
    pre_hook="truncate table {{ this }}"
) }}

WITH staged AS (
    SELECT
        format,
        rank,
        pokemon,
        base_spe,
        TRIM(BOTH '-' FROM LOWER(REGEXP_REPLACE(pokemon, '[^a-zA-Z0-9]+', '-', 'g'))) AS pokemon_slug
    FROM {{ source('staging', 'pikalytics_speed_tiers') }}
),

tiers AS (
    SELECT
        format,
        rank,
        pokemon,
        pokemon_slug,
        -- internal only (not selected): normalize the naive slug to the hub form to resolve the id
        CASE
            WHEN pokemon_slug ~ '^mega-.*-[xy]$' THEN REGEXP_REPLACE(pokemon_slug, '^mega-(.*)-([xy])$', '\1-mega-\2')
            WHEN pokemon_slug LIKE 'mega-%'      THEN REGEXP_REPLACE(pokemon_slug, '^mega-(.*)$', '\1-mega')
            WHEN pokemon_slug = 'basculegion-f'  THEN 'basculegion-female'
            ELSE pokemon_slug
        END AS hub_slug,
        base_spe,
        base_spe + 20                     AS neutral_0_sp,
        base_spe + 52                     AS neutral_32_sp,
        FLOOR((base_spe + 52) * 1.1)::int AS max_speed
    FROM staged
)

SELECT
    format,
    rank,
    pokemon,
    pokemon_slug,
    {{ resolve_pokemon_id('hub_slug') }} AS pokemon_id,
    base_spe,
    neutral_0_sp,
    neutral_32_sp,
    FLOOR(neutral_0_sp * 0.9)::int   AS neg_spe_0_sp,
    max_speed,
    FLOOR(max_speed * 1.5)::int      AS max_scarf,
    FLOOR(neutral_32_sp * 1.5)::int  AS neutral_32_scarf,
    pokemon ILIKE 'mega %'           AS is_mega,
    CASE
        WHEN base_spe >= 100 THEN 'fast'
        WHEN base_spe >=  60 THEN 'mid'
        ELSE                     'slow'
    END                              AS speed_bucket,
    'pikalytics'                     AS source
FROM tiers
ORDER BY rank
