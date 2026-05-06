{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

WITH latest AS (
    SELECT MAX(snapshot_month) AS snapshot_month
    FROM {{ source('staging', 'champions_speed_tiers') }}
)
SELECT
    s.rank,
    s.pokemon,
    s.base_spe,
    s.neutral_0_sp,
    s.neutral_32_sp,
    s.max_speed,
    s.neg_spe_0_sp,
    s.max_scarf,
    s.neutral_32_scarf,
    s.pokemon ILIKE 'mega %' AS is_mega,
    CASE
        WHEN s.base_spe >= 100 THEN 'fast'
        WHEN s.base_spe >=  60 THEN 'mid'
        ELSE                        'slow'
    END AS speed_bucket,
    s.snapshot_month
FROM {{ source('staging', 'champions_speed_tiers') }} AS s
INNER JOIN latest USING (snapshot_month)
ORDER BY s.rank