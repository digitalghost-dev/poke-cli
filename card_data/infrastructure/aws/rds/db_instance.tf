resource "aws_db_instance" "tfer--dagster-db" {
  allocated_storage                     = "20"
  auto_minor_version_upgrade            = "true"
  availability_zone                     = "us-west-2a"
  backup_retention_period               = "7"
  backup_target                         = "region"
  backup_window                         = "13:05-13:35"
  ca_cert_identifier                    = "rds-ca-rsa2048-g1"
  copy_tags_to_snapshot                 = "true"
  customer_owned_ip_enabled             = "false"
  database_insights_mode                = "standard"
  db_subnet_group_name                  = "${aws_db_subnet_group.tfer--poke-cli-db-subnet-group.name}"
  dedicated_log_volume                  = "false"
  deletion_protection                   = "false"
  depends_on = [aws_db_subnet_group.tfer--poke-cli-db-subnet-group]
  enabled_cloudwatch_logs_exports       = ["iam-db-auth-error", "postgresql"]
  engine                                = "postgres"
  engine_lifecycle_support              = "open-source-rds-extended-support-disabled"
  engine_version                        = "17.4"
  iam_database_authentication_enabled   = "false"
  identifier                            = "dagster-db"
  instance_class                        = "db.t3.micro"
  kms_key_id                            = var.kms-arn-key
  license_model                         = "postgresql-license"
  maintenance_window                    = "wed:08:28-wed:08:58"
  manage_master_user_password           = true
  max_allocated_storage                 = "1000"
  monitoring_interval                   = "60"
  monitoring_role_arn                   = "arn:aws:iam::940482406130:role/rds-monitoring-role"
  multi_az                              = "false"
  network_type                          = "IPV4"
  option_group_name                     = "default:postgres-17"
  parameter_group_name                  = "default.postgres17"
  performance_insights_enabled          = "true"
  performance_insights_kms_key_id       = var.kms-arn-key
  performance_insights_retention_period = "7"
  port                                  = "5432"
  publicly_accessible                   = "true"
  region                                = "us-west-2"
  skip_final_snapshot                   = true
  storage_encrypted                     = "true"
  storage_type                          = "gp3"

  tags = {
    project = "poke-cli"
  }

  tags_all = {
    project = "poke-cli"
  }

  username               = "postgres"
  vpc_security_group_ids = ["sg-026cc204887184c98", "sg-09ff0b46e3dd7a843"]
}
