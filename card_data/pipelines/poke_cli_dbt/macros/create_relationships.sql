{% macro create_relationships() %}
    {{ print("Dropping existing constraints...") }}

    -- Drop existing constraints if they exist (in reverse dependency order)
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".cards DROP CONSTRAINT IF EXISTS fk_cards_sets") %}
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".sets DROP CONSTRAINT IF EXISTS fk_sets_series") %}
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".cards DROP CONSTRAINT IF EXISTS pk_cards") %}
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".sets DROP CONSTRAINT IF EXISTS pk_sets") %}
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".series DROP CONSTRAINT IF EXISTS pk_series") %}

    {{ print("Adding primary keys...") }}

    -- Add primary keys
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".series ADD CONSTRAINT pk_series PRIMARY KEY (id)") %}
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".sets ADD CONSTRAINT pk_sets PRIMARY KEY (set_id)") %}
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".cards ADD CONSTRAINT pk_cards PRIMARY KEY (id)") %}

    {{ print("Adding foreign keys...") }}

    -- Add foreign keys
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".sets ADD CONSTRAINT fk_sets_series FOREIGN KEY (series_id) REFERENCES " ~ target.schema ~ ".series (id)") %}
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".cards ADD CONSTRAINT fk_cards_sets FOREIGN KEY (set_id) REFERENCES " ~ target.schema ~ ".sets (set_id)") %}

    {{ print("Relationships created successfully") }}

    {% do return('') %}
{% endmacro %}