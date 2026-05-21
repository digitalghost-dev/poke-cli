{{ config(
    materialized='table',
    post_hook=[
        "ALTER TABLE {{ this }} ADD PRIMARY KEY (id)",
        "{{ enable_rls() }}"
    ]
) }}

SELECT
    id,
    pokedata_id,
    game_type,
    name,
    start_date::date,
    end_date::date,
    season,
    count,
    rounds,
    last_updated::timestamp
FROM {{ source('staging', 'comp_events') }}
