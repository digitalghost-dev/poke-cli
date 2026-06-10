output "aws_db_instance_tfer--dagster-db_id" {
  value = "${aws_db_instance.tfer--dagster-db.id}"
}

output "aws_db_subnet_group_tfer--poke-cli-db-subnet-group_id" {
  value = "${aws_db_subnet_group.tfer--poke-cli-db-subnet-group.id}"
}
