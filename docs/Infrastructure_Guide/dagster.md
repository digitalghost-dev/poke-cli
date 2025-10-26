---
weight: 7
---

# 7. Dagster

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

## 

## Automating Startup with `systemd`
In order to save on costs, the EC2 and RDS instances are scheduled to start and stop once each day. To automate the
starting of the Dagster webservice, `systemd`, along with a couple of shell scripts, will be used to create this
automation.

This project, under `card_data/infrastructure/` has a file for starting the Dagster web service called `start-dagster.sh`
and is set up to be used.


```shell
aws rds describe-db-instances \
--region us-west-2 \
--query 'DBInstances[*].[DBInstanceIdentifier,Endpoint.Address,Endpoint.Port]' \
--output table
```

```shell
aws secretsmanager create-secret \
    --name dagster/supabase-creds \
    --secret-string '{"password":"your_password","user":"your_user"}' \
    --region us-west-2
```

```shell
nano /home/ubuntu/wait-for-rds.sh

nano /home/ubuntu/start-dagster.sh

sudo nano /etc/systemd/system/dagster.service
```

```shell
chmod +x /home/ubuntu/start-dagster.sh

chmod +x /home/ubuntu/wait-for-rds.sh
```

create `