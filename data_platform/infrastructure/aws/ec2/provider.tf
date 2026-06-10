provider "aws" {
  region = "us-west-2"
}

terraform {
  cloud {
    organization = "digitalghost-dev"

    workspaces {
      project = "poke-cli"
      name = "ec2"
    }
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.7.0"
    }
  }
}