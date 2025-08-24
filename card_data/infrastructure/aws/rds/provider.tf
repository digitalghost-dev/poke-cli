provider "aws" {
  profile = "terraform-user"
  region = "us-west-2"
}

terraform {
  cloud {
    organization = "digitalghost-dev"

    workspaces {
      project = "poke-cli"
      name = "poke-cli"
    }
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.7.0"
    }
  }
}