{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

SELECT name, card_number, market_price
FROM {{ source('staging', 'pricing_data') }}