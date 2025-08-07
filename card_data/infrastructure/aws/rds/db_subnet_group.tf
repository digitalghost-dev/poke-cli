resource "aws_db_subnet_group" "tfer--poke-cli-db-subnet-group" {
  description = "Subnet group for RDS databases"
  name        = "poke-cli-db-subnet-group"
  region      = "us-west-2"
  subnet_ids  = ["subnet-08dae4b7aede93128", "subnet-0eeb519cf23a763bf", "subnet-04fe6e100221b27d4", "subnet-0be3aac807720c1d6"]
}
