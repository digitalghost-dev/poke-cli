---
weight: 1
---

# 1 // Overview

This section serves as a knowledge base for the project’s backend data infrastructure. It was created for a few purposes:

1.	To document how I built everything, so I can easily reference it later.
2.	To help others learn how to build something similar.
3.	To provide a clear understanding of how all the components work together.

This infrastructure guide will be constantly updated as changes or updates to resources occur.

The infrastructure is mainly in support of TCG data since there are multiple data sources that are being used to support the CLI/TUI.
The VGC data simply calls one API.

!!! warning

    All of the commands ran for this project in the terminal are based on macOS (i.e., Homebrew to install packages). 
    If building on a different operating system, please find the equivalent command. Links will be provided for install 
    guides for all operating systems when possible.

## Data Infrastructure Diagram
![data_infrastructure_diagram](../assets/data_infrastructure_diagram.svg)

1. TCGPlayer pricing data and TCGDex card data are called and processed through a data pipeline orchestrated by Dagster and hosted on AWS.
    - Dagster runs on an EC2 instance.
    - Dagster metadata is stored separately in RDS.
    - The pricing pipeline is scheduled with cron: `0 14 * * *`.
    - Tournament standings data is also pulled from Limitless.

2. When the pipeline starts, Pydantic validates the incoming API data against a pre-defined schema, ensuring the data types match the expected structure.
    - Invalid or unexpected payloads fail early before data is loaded downstream.

3. Polars is used to create DataFrames.
    - DataFrames are used to clean, normalize, and prepare records for database loading.

4. The data is loaded into a Supabase staging schema.
    - The staging schema acts as the raw/validated landing area before production tables are built.

5. Soda data quality checks are performed.
    - Checks validate expectations such as row counts, required columns, missing values, duplicate keys, and URL formats.

6. dbt runs tests and builds the final tables in a Supabase production schema.
    - dbt transforms staged data into the final public-facing models.
    - The production schema powers TCG/card/tournament queries.

7. Users are then able to query the `pokeapi.co` or `supabase` APIs for either video game or trading card data, respectively.
    - The CLI uses PokéAPI for video game data.
    - The CLI and Streamlit web app use Supabase for TCG data.
    - Dagster run status is sent through an n8n webhook for Discord notifications.

## Tools & Services

Below is a list of all the tools and services used in this project's infrastructure:

- AWS
    - RDS
    - S3
    - VPC
    - EC2
- Dagster
- Polars
- Supabase
- Terraform

## Project Layout

```text
.
├── infrastructure/
│   ├── aws/
│   │   ├── .terraform
│   │   ├── ec2/
│   │   ├── rds/
│   │   └── vpc/
│   ├── dagster.server
│   ├── start-dagster.sh
│   └── wait-for-rds.sh
├── pipelines/
│   ├── defs/
│   │   ├── extract/
│   │   │   ├── tcgcsv/
│   │   │   │   └── extract_pricing.py
│   │   │   └── tcgdex/
│   │   │       ├── extract_cards.py
│   │   │       ├── extract_series.py
│   │   │       └── extract_sets.py
│   │   ├── load/
│   │   │   ├── tcgcsv/
│   │   │   │   └── load_pricing.py
│   │   │   └── tcgdex/
│   │   │       ├── load_cards.py
│   │   │       ├── load_series.py
│   │   │       └── load_sets.py
│   │   └── transform/
│   │       └── transform_data.py
│   ├── poke_cli_dbt/
│   │   ├── macros/
│   │   │   ├── create_relationships.sql
│   │   │   ├── create_rls.sql
│   │   │   └── create_view.sql
│   │   ├── models/
│   │   │   ├── cards.sql
│   │   │   ├── pricing_data.sql
│   │   │   ├── series.sql
│   │   │   ├── sets.sql
│   │   │   └── sources.yml
│   │   ├── dbt_project.yml
│   │   └── profiles.yml
│   ├── soda/
│   │   ├── checks_pricing.yml
│   │   ├── checks_series.yml
│   │   ├── checks_sets.yml
│   │   └── configuration.yml
│   ├── tests/
│   │   └── extract_series_test.py
│   └── utils/
│       ├── json_retriever.py
│       └── secret_retriever.py
├── dagster.yaml
├── pyproject.toml
└── uv.lock    
```

!!! note

    This project is a learning playground for exploring new tools, services, and programming languages.
    Some design choices are intentionally experimental or may not follow conventional patterns—
    that's part of the learning process!
    
    Feedback and suggestions are always welcome! If you spot an issue or have ideas for improvement,
    please open a [GitHub Issue](https://github.com/digitalghost-dev/poke-cli/issues).
