---
weight: 6
---

# 6. Dagster

!!! question "What is Dagster?"

    Dagster is an open-source data orchestration tool that helps you build, run, and monitor your data pipelines. 
    It’s designed to make working with data workflows more reliable and maintainable, giving you clear visibility 
    into each step and making it easier to catch issues before they cause problems.

    View more [about Dagster](https://dagster.io/platform-overview)

## Installation
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

