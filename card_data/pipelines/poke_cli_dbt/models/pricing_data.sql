{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

SELECT product_id, name, card_number, market_price
FROM {{ source('staging', 'pricing_data') }}