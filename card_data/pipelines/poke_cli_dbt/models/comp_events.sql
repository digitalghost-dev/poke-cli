{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

SELECT
    start_date::date,
    end_date::date,
    season,
    tcg_pokedata_id,
    vg_pokedata_id
FROM
    {{ source('staging', 'comp_events') }}
