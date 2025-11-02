---
weight: 3
---

# 3 // Local Deployment
This page explains how to set up Python, dbt, Dagster, and other tools for running the data pipelines
locally. The [4. AWS](aws.md) section will show how to deploy this solution on the cloud.

## Python

### Installing uv
_uv is the main package and project manager used in this project._

Learn more about [uv](https://docs.astral.sh/uv/).

1. Install via their [installation script](https://docs.astral.sh/uv/getting-started/installation/):
    ```bash
    curl -LsSf https://astral.sh/uv/install.sh | sh
    ```
   or brew:
    ```bash
    brew install uv
    ```
2. Install Python with `uv`:
    ```bash
    uv python install 3.12
    ```
   
## dbt

!!! question "What is dbt?"

    dbt (data build tool) is a command-line tool that enables data analysts and engineers to transform data 
    in warehouses using SQL. The tool allows users to write modular SQL queries as "models" that build upon 
    each other, automatically managing dependencies and enabling version control, testing, and documentation 
    of data transformations. dbt compiles SQL models into executable queries and runs them in the proper order,
    turning raw data into analysis-ready datasets.

    View more [about dbt](https://www.getdbt.com/product/what-is-dbt)

### Installation & Initialization

Install with `uv`:
```bash
uv add dbt
```

Initialize a `dbt` project in the `card_data` directory:
```bash
dbt init
```

Follow the prompts to finish setting up the `dbt` project.

### Models
Models are the pieces of SQL code that run when using that `dbt build` command that _build_ the
tables to the destination schema. In this project, that would the `public` schema in the PostgreSQL
database on Supabase.

The `public` schema is the public facing schema that exposes the API to the data in the tables.

### Sources
Create a `source.yml` file under the `models/` directory. 
More info on [sources here](https://docs.getdbt.com/docs/build/sources).

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

---

## Dagster

!!! question "What is Dagster?"

    Dagster is an open-source data orchestration tool that helps you build, run, and monitor your data pipelines. 
    It’s designed to make working with data workflows more reliable and maintainable, giving you clear visibility 
    into each step and making it easier to catch issues before they cause problems.

    View more [about Dagster](https://dagster.io/platform-overview)

### Installation
Dagster and its components needed for the project can be installed with `uv`: `

```bash
uv add dagster dagster-webserver dagster-dg-cli dagster-postgres>=0.27.3 dagster-dbt
```

## Project Layout
In my experience, Dagster needed a specific directory structure in order for the program to find all necessary files.
This project uses a directory named `pipelines` to store all the Dagster files:

```
.
└── pipelines/
    ├── defs/
    │   ├── extract/
    │   │   └── extract_data.py
    │   ├── load/
    │   │   └── load_data.py
    │   └── transformation/
    │       └── transform_data.py
    ├── poke_cli_dbt/
    │   ├── logs
    │   ├── macros/
    │   │   ├── create_relationships.sql
    │   │   └── create_rls.sql
    │   ├── models/
    │   │   ├── cards.sql
    │   │   ├── series.sql
    │   │   ├── sets.sql
    │   │   └── sources.yml
    │   ├── target
    │   ├── dbt_project.yml
    │   └── profiles.yml
    └── soda/
        ├── checks.yml
        └── configuration.yml
```