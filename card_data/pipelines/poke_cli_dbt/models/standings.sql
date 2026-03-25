{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

SELECT
    s.rank,
    s.name,
    s.points,
    s.record,
    s.opp_win_percent,
    s.opp_opp_win_percent,
    s.deck,
    s.decklist,
    c.country_name AS player_country,
    LOWER(c.code) AS country_code,
    t.location,
    t.country_code AS iso_code,
    t.latitude AS tournament_latitude,
    t.longitude AS tournament_longitude,
    t.logo,
    t.start_date,
    t.end_date,
    t.text_date,
    t.type,
    t.player_quantity
FROM
    {{ source('staging', 'standings') }} AS s
INNER JOIN {{ source('staging', 'tournaments') }} AS t
   ON s.tournament_id = t.tournament_id
LEFT JOIN {{ source('staging', 'country_codes') }} AS c
   ON s.country = c.code