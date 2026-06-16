{% macro create_pikalytics_relationships() %}
    {% if not execute %}{% do return('') %}{% endif %}

    {% set targets = ['pikalytics_speed_tiers', 'pikalytics_usage', 'pikalytics_pokemon_comp_info'] %}
    {% set ran = [] %}
    {% for res in results %}
        {% if res.node.name in targets %}{% do ran.append(res.node.name) %}{% endif %}
    {% endfor %}
    {% if ran | length == 0 %}{% do return('') %}{% endif %}

    {{ print("Dropping existing Pikalytics constraints...") }}

    -- Drop existing foreign keys if they exist
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".pikalytics_speed_tiers DROP CONSTRAINT IF EXISTS fk_pikalytics_speed_tiers_pokemon") %}
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".pikalytics_usage DROP CONSTRAINT IF EXISTS fk_pikalytics_usage_pokemon") %}
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".pikalytics_pokemon_comp_info DROP CONSTRAINT IF EXISTS fk_pikalytics_pokemon_comp_info_pokemon") %}

    {{ print("Adding Pikalytics foreign keys...") }}

    -- Link pokemon_id to the central pokemon hub (NULLs allowed; top_teams omitted — its
    -- Pokémon live in a JSONB array, which Postgres cannot foreign-key into)
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".pikalytics_speed_tiers ADD CONSTRAINT fk_pikalytics_speed_tiers_pokemon FOREIGN KEY (pokemon_id) REFERENCES " ~ target.schema ~ ".pokemon (id)") %}
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".pikalytics_usage ADD CONSTRAINT fk_pikalytics_usage_pokemon FOREIGN KEY (pokemon_id) REFERENCES " ~ target.schema ~ ".pokemon (id)") %}
    {% do run_query("ALTER TABLE " ~ target.schema ~ ".pikalytics_pokemon_comp_info ADD CONSTRAINT fk_pikalytics_pokemon_comp_info_pokemon FOREIGN KEY (pokemon_id) REFERENCES " ~ target.schema ~ ".pokemon (id)") %}

    {{ print("Pikalytics relationships created successfully") }}

    {% do return('') %}
{% endmacro %}
