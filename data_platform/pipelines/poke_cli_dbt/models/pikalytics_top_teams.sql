{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

-- Replaces pikalytics_comp_top_teams. n8n lands RAW rows (rank, author, record,
-- tournament, archetypes, pokemon, web_url) in staging.pikalytics_top_teams (full
-- replace each run). All derivations live here in SQL (moved out of the n8n Code node):
-- wins/losses/ties parsed from record, format/source constants, and pokemon_ids —
-- a parallel array linking each team member to the public.pokemon hub via resolve_pokemon_id.

WITH staged AS (
    SELECT
        rank,
        author,
        record,
        tournament,
        archetypes,
        pokemon,
        web_url
    FROM {{ source('staging', 'pikalytics_top_teams') }}
),

record_parsed AS (
    SELECT
        s.*,
        (
            SELECT ARRAY_AGG(m[1]::int)
            FROM REGEXP_MATCHES(s.record, '[0-9]+', 'g') AS m
        ) AS record_nums
    FROM staged s
)

SELECT
    'gen9championsvgc2026regma'  AS format,
    rank,
    author,
    record,
    record_nums[1]               AS wins,
    record_nums[2]               AS losses,
    record_nums[3]               AS ties,
    tournament,
    archetypes,
    pokemon,
    (
        SELECT JSONB_AGG(
                   {{ resolve_pokemon_id("TRIM(BOTH '-' FROM LOWER(REGEXP_REPLACE(elem, '[^a-zA-Z0-9]+', '-', 'g')))") }}
                   ORDER BY ord
               )
        FROM JSONB_ARRAY_ELEMENTS_TEXT(pokemon) WITH ORDINALITY AS t(elem, ord)
        WHERE TRIM(elem) <> ''
    )                            AS pokemon_ids,
    web_url,
    'pikalytics'                 AS source
FROM record_parsed
ORDER BY rank
