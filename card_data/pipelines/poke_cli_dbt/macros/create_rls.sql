{% macro enable_rls() %}
    ALTER TABLE {{ this }} ENABLE ROW LEVEL SECURITY;
    CREATE POLICY "Enable read access for all users" ON {{ this }} TO PUBLIC USING (true);
{% endmacro %}