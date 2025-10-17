resource "aws_vpc" "tfer--vpc-0f9d6a37031fa6597" {
  assign_generated_ipv6_cidr_block     = "false"
  cidr_block                           = "10.0.0.0/20"
  enable_dns_hostnames                 = "true"
  enable_dns_support                   = "true"
  enable_network_address_usage_metrics = "false"
  instance_tenancy                     = "default"
  region                               = "us-west-2"

  tags = {
    Name    = "poke-cli-vpc"
    project = "poke-cli"
  }

  tags_all = {
    Name    = "poke-cli-vpc"
    project = "poke-cli"
  }
}
