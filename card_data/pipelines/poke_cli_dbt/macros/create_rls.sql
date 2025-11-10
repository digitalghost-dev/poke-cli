{% macro enable_rls() %}
    ALTER TABLE {{ this }} ENABLE ROW LEVEL SECURITY;
    CREATE POLICY "Enable Read Access for All Users"
        ON {{ this }}
        AS PERMISSIVE
        FOR SELECT
        TO PUBLIC
        USING (true);
{% endmacro %}