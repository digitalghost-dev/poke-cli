provider "aws" {
	profile = "terraform-user"
	region = "us-west-2"
}

terraform {
	required_providers {
		aws = {
	    version = "~> 6.7.0"
		}
  }
}
