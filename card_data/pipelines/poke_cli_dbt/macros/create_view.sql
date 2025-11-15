{% macro create_view() %}
    CREATE OR REPLACE VIEW public.card_pricing_view
    WITH (security_invoker = true) AS
        WITH cards_cte AS (
            SELECT
                set_id,
                name,
                "localId",
                "set_cardCount_official",
                CONCAT(name, ' - ', "localId", '/', "set_cardCount_official") AS card_combined_name,
                set_name
            FROM public.cards
        ),

             cards_pricing_cte AS (
                 SELECT
                     product_id,
                     market_price,
                     CONCAT(name, ' - ', card_number) AS card_combined_name,
                     card_number
                 FROM public.pricing_data
             )

        SELECT
            c.set_id,
            c.name,
            CONCAT(p.card_number, ' - ', c.name) AS number_plus_name,
            c.set_name,
            c."localId",
            p."market_price",
            p."card_number"
        FROM
            cards_cte AS c
                INNER JOIN
            cards_pricing_cte AS p
            ON c.card_combined_name = p.card_combined_name
        ORDER BY c."localId"
{% endmacro %}