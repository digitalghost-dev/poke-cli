{% macro create_view() %}
    CREATE OR REPLACE VIEW public.card_pricing_view
    WITH (security_invoker = true) AS
        WITH cards_cte AS (
            SELECT
                set_id,
                name,
                image,
                illustrator,
                "localId",
                "set_cardCount_official",
                CONCAT(name, ' - ', LPAD("localId", 3, '0'), '/', LPAD("set_cardCount_official"::text, 3, '0')) AS card_combined_name,
                set_name
            FROM public.cards
        ),
         cards_pricing_cte AS (
             SELECT
                 product_id,
                 market_price,
                 CONCAT(REPLACE(name, ' (Secret)', ''), ' - ', card_number) AS card_combined_name,
                 card_number
             FROM public.pricing_data
         )
        SELECT
            c.set_id,
            c.name,
            CONCAT(COALESCE(p.card_number, LPAD(c."localId", 3, '0')), ' - ', c.name) AS number_plus_name,
            CONCAT(c.image, '/high.png') AS image_url,
            c.set_name,
            LPAD(c."localId", 3, '0') AS "localId",
            p."market_price",
            COALESCE(p."card_number", LPAD(c."localId", 3, '0')) AS card_number,
            c.illustrator
        FROM
            cards_cte AS c
                LEFT JOIN
            cards_pricing_cte AS p
            ON c.card_combined_name = p.card_combined_name
        ORDER BY c."localId"::integer
{% endmacro %}