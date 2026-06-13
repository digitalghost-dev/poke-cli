{{ config(
    materialized='view',
    post_hook=[
        "ALTER VIEW {{ this }} SET (security_invoker = on)",
        "GRANT SELECT ON {{ this }} TO anon, authenticated"
    ]
) }}

SELECT
    p.placement                                           AS rank,
    p.player_name                                         AS name,
    (p.wins * 3 + p.ties)                                 AS points,
    (p.wins || ' - ' || p.losses || ' - ' || p.ties)     AS record,
    to_char(p.resistance_opp * 100, 'FM990.00') || '%'    AS opp_win_percent,
    to_char(p.resistance_oppopp * 100, 'FM990.00') || '%' AS opp_opp_win_percent,
    d.decklist                                            AS team,
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
        WHEN e.name ILIKE '%Special%'      THEN 'Special'
        WHEN e.name ILIKE '%World%'        THEN 'Worlds'
    END                                                   AS type,
    e.count                                               AS player_quantity
FROM {{ ref('comp_players') }} p
JOIN {{ ref('comp_events') }} e
    ON e.pokedata_id = p.pokedata_id AND e.game_type = p.game_type
LEFT JOIN {{ ref('comp_vg_decklists') }} d
    ON d.pokedata_id = p.pokedata_id
   AND d.game_type = p.game_type
   AND d.player_name = p.player_name
LEFT JOIN {{ source('staging', 'country_codes') }} cc ON cc.code = p.country
WHERE p.game_type = 'VGC'
  AND e.season >= (
      CASE WHEN extract(month FROM CURRENT_DATE) >= 9
           THEN extract(year FROM CURRENT_DATE) + 1
           ELSE extract(year FROM CURRENT_DATE)
      END
  )
