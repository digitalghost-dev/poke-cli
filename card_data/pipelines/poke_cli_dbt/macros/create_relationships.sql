{% macro create_relationships() %}
    ALTER TABLE {{ target.schema }}.series ADD CONSTRAINT pk_series PRIMARY KEY (id);
    ALTER TABLE {{ target.schema }}.sets ADD CONSTRAINT pk_sets PRIMARY KEY (set_id);
    ALTER TABLE {{ target.schema }}.cards ADD CONSTRAINT pk_cards PRIMARY KEY (id);

    ALTER TABLE public.sets
    ADD CONSTRAINT fk_sets_series
    FOREIGN KEY (series_id)
        REFERENCES public.series (id);

    ALTER TABLE public.cards
    ADD CONSTRAINT fk_cards_sets
    FOREIGN KEY (set_id)
    REFERENCES public.sets (set_id);

{% endmacro %}