{% macro resolve_pokemon_id(slug_column) %}
  COALESCE(
    (SELECT p.id FROM {{ ref('pokemon') }} p WHERE p.identifier = {{ slug_column }}),
    (SELECT p.id FROM {{ ref('pokemon') }} p
       WHERE p.is_default = 1 AND p.identifier LIKE {{ slug_column }} || '-%' LIMIT 1)
  )
{% endmacro %}
