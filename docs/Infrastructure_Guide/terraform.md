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

### Profile
Before installing Terraformer, a profile from AWS will need to be configured under `~/.aws/credentials` which is an INI-like file (AWS credentials format).

Configure the `terraform-user` profile from [1. AWS](aws.md#iam) to look like this:

```ini
[terraform-user]
aws_access_key_id = <ACCESS_KEY_ID>
aws_secret_access_key = <SECRET_ACCESS_KEY>
```

To retrieve these options, head to IAM > Users > terraform-user > Security Credentials tab

### Install
* [Install Guide](https://github.com/GoogleCloudPlatform/terraformer?tab=readme-ov-file#installation)

Install Terraformer on macOS:

```bash
brew install terraformer
```

### Run

Once Terraformer is installed, run the tool to start generated `.tf` files.

For example, to retrieve the build for an RDS instance, run:

```bash
terraformer import aws --regions us-west-2 --resources rds --profile terraform-user
```

Where:
* `--regions `is where the RDS instance is located.
* `--resources` is the resources under RDS
* `--profile` is the `terraform-user` profile from `~/.aws/credentials`

After running this command, Terraformer will create the `.tf` files under `~/generated`.
The path to where these will be created can be changed with the `--path-output` flag.