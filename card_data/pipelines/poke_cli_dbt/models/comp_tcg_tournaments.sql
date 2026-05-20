-- depends_on: {{ ref('comp_events') }}

{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

SELECT
    event_id,
    name,
    count,
    rounds,
    last_updated::timestamp
FROM
    {{ source('staging', 'comp_tcg_tournaments') }}
