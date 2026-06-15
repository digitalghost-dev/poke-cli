SELECT
    pokedata_id,
    game_type,
    name,
    location
FROM {{ ref('comp_events') }}
WHERE location LIKE '%Championships%'
