{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

-- Replaces champions_speed_tiers. n8n lands RAW rows (format, rank, pokemon, base_spe)
-- in staging.pikalytics_speed_tiers (full replace each run); all tier math lives here
-- in SQL (moved out of the n8n Code node) so the derivations are version-controlled.

WITH tiers AS (
    SELECT
        format,
        rank,
        pokemon,
        base_spe,
        base_spe + 20                        AS neutral_0_sp,
        base_spe + 52                        AS neutral_32_sp,
        FLOOR((base_spe + 52) * 1.1)::int    AS max_speed
    FROM {{ source('staging', 'pikalytics_speed_tiers') }}
)

SELECT
    format,
    rank,
    pokemon,
    TRIM(BOTH '-' FROM LOWER(REGEXP_REPLACE(pokemon, '[^a-zA-Z0-9]+', '-', 'g'))) AS pokemon_slug,
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
