{% macro enable_rls(role='PUBLIC', policy_name='Enable Read Access for All Users') %}
    ALTER TABLE {{ this }} ENABLE ROW LEVEL SECURITY;
    DROP POLICY IF EXISTS "{{ policy_name }}" ON {{ this }};
    CREATE POLICY "{{ policy_name }}"
        ON {{ this }}
        AS PERMISSIVE
        FOR SELECT
        TO {{ role }}
        USING (true);
{% endmacro %}