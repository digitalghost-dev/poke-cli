{% macro create_relationships() %}
    ALTER TABLE {{ ref('series') }} ADD CONSTRAINT pk_series PRIMARY KEY (id);
    ALTER TABLE {{ ref('sets') }} ADD CONSTRAINT pk_sets PRIMARY KEY (set_id);
    ALTER TABLE {{ ref('cards') }} ADD CONSTRAINT pk_cards PRIMARY KEY (id);
{% endmacro %}