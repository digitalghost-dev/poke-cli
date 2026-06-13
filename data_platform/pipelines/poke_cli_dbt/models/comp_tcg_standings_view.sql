{{ config(
    materialized='view',
    post_hook=[
        "ALTER VIEW {{ this }} SET (security_invoker = on)",
        "GRANT SELECT ON {{ this }} TO anon, authenticated"
    ]
) }}

WITH lookup AS (  -- pokedata_id -> Limitless tournament_id (TCG only)
    SELECT ce.pokedata_id, COALESCE(loc.tid, sole.tid) AS tournament_id
    FROM (
        SELECT pokedata_id, name, start_date
        FROM {{ ref('comp_events') }}
        WHERE game_type = 'TCG'
    ) ce
    LEFT JOIN LATERAL (  -- tier 1: same date + city name appears in the event name
        SELECT t.tournament_id AS tid
        FROM {{ source('staging', 'tournaments') }} t
        WHERE t.start_date = ce.start_date
          AND ce.name ILIKE '%' || t.location || '%'
        LIMIT 1
    ) loc ON true
    LEFT JOIN LATERAL (  -- tier 2: fallback to the sole tournament on that date
        SELECT t.tournament_id AS tid
        FROM {{ source('staging', 'tournaments') }} t
        JOIN (
            SELECT start_date
            FROM {{ source('staging', 'tournaments') }}
            GROUP BY start_date
            HAVING count(*) = 1
        ) u ON u.start_date = t.start_date
        WHERE t.start_date = ce.start_date
        LIMIT 1
    ) sole ON true
)

SELECT
    p.placement                                           AS rank,
    p.player_name                                         AS name,
    (p.wins * 3 + p.ties)                                 AS points,
    (p.wins || ' - ' || p.losses || ' - ' || p.ties)     AS record,
    to_char(p.resistance_opp * 100, 'FM990.00') || '%'    AS opp_win_percent,
    to_char(p.resistance_oppopp * 100, 'FM990.00') || '%' AS opp_opp_win_percent,
    s.deck,
    s.decklist,
    cc.country_name                                       AS player_country,
    lower(p.country)                                      AS country_code,
    e.location,
    e.start_date,
    e.end_date,
    CASE
        WHEN e.start_date = e.end_date
            THEN to_char(e.start_date, 'FMMonth FMDD') || ', ' || to_char(e.start_date, 'YYYY')
        WHEN extract(month FROM e.start_date) = extract(month FROM e.end_date)
            THEN to_char(e.start_date, 'FMMonth FMDD') || '–' || to_char(e.end_date, 'FMDD') || ', ' || to_char(e.end_date, 'YYYY')
        ELSE to_char(e.start_date, 'FMMonth FMDD') || '–' || to_char(e.end_date, 'FMMonth FMDD') || ', ' || to_char(e.end_date, 'YYYY')
    END                                                   AS text_date,
    CASE
        WHEN e.name ILIKE '%International%' THEN 'International'
        WHEN e.name ILIKE '%Regional%'     THEN 'Regional'
        WHEN e.name ILIKE '%World%'        THEN 'Worlds'
    END                                                   AS type,
    e.count                                               AS player_quantity
FROM {{ ref('comp_players') }} p
JOIN {{ ref('comp_events') }} e
    ON e.pokedata_id = p.pokedata_id AND e.game_type = p.game_type
LEFT JOIN lookup lk ON lk.pokedata_id = p.pokedata_id
LEFT JOIN LATERAL (  -- one deck/decklist per player; rank=placement disambiguates same-name players
    SELECT s.deck, s.decklist
    FROM {{ source('staging', 'standings') }} s
    WHERE s.tournament_id = lk.tournament_id
      AND s.name = p.player_name
    ORDER BY (s.rank = p.placement) DESC, s.rank
    LIMIT 1
) s ON true
LEFT JOIN {{ source('staging', 'country_codes') }} cc ON cc.code = p.country
WHERE p.game_type = 'TCG'
  AND e.season >= (
      CASE WHEN extract(month FROM CURRENT_DATE) >= 9
           THEN extract(year FROM CURRENT_DATE) + 1
           ELSE extract(year FROM CURRENT_DATE)
      END
  )