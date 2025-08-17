---
weight: 3
---

# 3. Terraform

!!! question "What is Terraform?"

    Terraform is an open-source Infrastructure as Code (IaC) tool for automating cloud infrastructure using code. 
    You define what you want (servers, networks, etc.) in config files, and Terraform builds and manages it across cloud providers.

    Instead of manually setting up servers, networks, or other resources, you write code to describe the desired state of your infrastructure, 
    and Terraform takes care of creating and maintaining it.

## Overview

Terraform is used in this project to build/destroy the following resources:

- AWS
    - RDS
    - S3
    - VPC
- Supabase

---

## Terraformer

After manually creating the resources in the [AWS](aws.md) section, [terraformer](https://github.com/GoogleCloudPlatform/terraformer) can now be used to create the `terraform` files.
Terraformer acts as a "reverse Terraform" tool where it can read from created resources on AWS. 
Although the output can sometimes be a bit verbose, it's an easy way to capture everything to create reproducible builds.

### Install
* [Install Guide](https://github.com/GoogleCloudPlatform/terraformer?tab=readme-ov-file#installation)

Install Terraformer:

```bash
brew install terraformer
```


