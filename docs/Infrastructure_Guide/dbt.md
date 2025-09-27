---
weight: 5
---

# 5. dbt

!!! question "What is dbt?"

    dbt (data build tool) is a command-line tool that enables data analysts and engineers to transform data 
    in warehouses using SQL. The tool allows users to write modular SQL queries as "models" that build upon 
    each other, automatically managing dependencies and enabling version control, testing, and documentation 
    of data transformations. dbt compiles SQL models into executable queries and runs them in the proper order,
    turning raw data into analysis-ready datasets.

    View more [about dbt](www.getdbt.com/product/what-is-dbt)

## Installation & Initialization

Install with `uv`:
```bash
uv add dbt
```

Initialize a `dbt` project in the `card_data` directory:
```bash
dbt init
```

Follow the prompts to finish setting up the `dbt` project.

## Models
Models are the pieces of SQL code that run when using that `dbt build` command that _build_ the
tables to the destination schema. In this project, that would the `public` schema in the PostgreSQL
database on Supabase.

The `public` schema is the public facing schema that exposes the API to the data in the tables.

## Sources
Create a `source.yml` file under the `models/` directory. More info on [sources here](https://docs.getdbt.com/docs/build/sources).

This file is used to declare and configure the raw data sources. These tables are the foundation for
the dbt models but are not managed by dbt itself.

For example:
```yaml
sources:
  - name: staging
    description: "Staging schema containing raw data loaded from extract pipeline"
    tables:
      - name: series
        description: "Pokemon card series data"
        columns:
          - name: id
            description: "Unique series identifier"
          - name: name
            description: "Series name"
          - name: logo
            description: "Series logo URL"
```

The above `yml` defines the structure for the raw `series` table from the `staging` schema.