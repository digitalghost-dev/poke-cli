-- depends_on: {{ ref('pokemon') }}
{{ config(
    materialized='incremental',
    unique_key='pokemon_id',
    incremental_strategy='merge',
    post_hook="{{ enable_rls(role='authenticated', policy_name='Enable Read Access for Authenticated Users') }}"
) }}

SELECT
    COALESCE(
        substring(gif_sprite_url FROM '/(\d+)\.gif$'),
        substring(png_sprite_url FROM '/(\d+)\.png$')
    )::int AS pokemon_id,
    gif_sprite_url,
    png_sprite_url
FROM {{ source('staging', 'vg_pokemon_sprites') }}
