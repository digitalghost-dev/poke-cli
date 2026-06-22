{{ config(
    materialized='view',
    post_hook=[
        "ALTER VIEW {{ this }} SET (security_invoker = on)",
        "GRANT SELECT ON {{ this }} TO anon, authenticated"
    ]
) }}

SELECT
    p.id AS pokemon_id,
    p.id AS dex_number,
    p.identifier AS slug,
    INITCAP(REPLACE(p.identifier, '-', ' ')) AS name,
    s.gif_sprite_url,
    s.png_sprite_url,
    COALESCE(
        JSONB_AGG(
            JSONB_BUILD_OBJECT(
                'slot', pt.slot,
                'name', t.identifier
            )
            ORDER BY pt.slot
        ) FILTER (WHERE t.identifier IS NOT NULL),
        '[]'::JSONB
    ) AS types
FROM {{ ref('pokemon') }} p
LEFT JOIN {{ ref('vg_pokemon_sprites') }} s
    ON s.pokemon_id = p.id
LEFT JOIN {{ ref('vg_pokemon_types') }} pt
    ON pt.pokemon_id = p.id
LEFT JOIN {{ ref('vg_types') }} t
    ON t.id = pt.type_id
GROUP BY
    p.id,
    p.identifier,
    s.gif_sprite_url,
    s.png_sprite_url
