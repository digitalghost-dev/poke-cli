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

## Installation

## Sources
* [docs](https://docs.getdbt.com/docs/build/sources)
Create a `source.yml` file under the `models/` directory
