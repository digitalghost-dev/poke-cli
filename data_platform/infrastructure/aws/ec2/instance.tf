resource "aws_instance" "tfer--i-01dbf82e21c0da38f_dagster-webserver" {
  ami                    = "ami-0326baaa98cf958ed"
  instance_type          = "t4g.small"
  key_name               = "dagster-vm-key-pair"
  subnet_id              = "subnet-04fe6e100221b27d4"
  vpc_security_group_ids = ["sg-04c2a30cb05044ad6", "sg-026cc204887184c98"]
  ebs_optimized          = true
  monitoring             = false

  root_block_device {
    delete_on_termination = true
  }

  tags = {
    Name = "dagster-webserver"
  }

  tags_all = {
    Name = "dagster-webserver"
  }
}
