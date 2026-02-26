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
    t.location,
    t.start_date,
    t.end_date,
    t.type,
    t.player_quantity
FROM
    {{ source('staging', 'standings') }} AS s
INNER JOIN {{ source('staging', 'tournaments') }} AS t
   ON s.tournament_id = t.tournament_id
INNER JOIN {{ source('staging', 'country_codes') }} AS c
   ON s.country = c.code