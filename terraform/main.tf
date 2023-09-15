terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }
  required_version = ">= 1.2.0"
}

provider "aws" {
  region     = "us-east-1"
  access_key = ""
  secret_key = ""
}

resource "aws_security_group" "rollandplay_sg" {
  name = "rollandplay_sg"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1" # Allow all outbound traffic
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "app_server" {
  ami           = "ami-01c647eace872fc02"
  instance_type = "t2.micro"
  key_name      = "rollandplayapi"

  security_groups = [aws_security_group.rollandplay_sg.name]

  tags = {
    Name = "rollandplay"
  }

  user_data = <<-EOF
    #!/bin/bash

    sudo yum install -y golang

    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo 'export GOPATH=$HOME/go' >> ~/.bashrc
    echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc

    source ~/.bashrc
  EOF
}

output "instance_ip" {
  value = aws_instance.app_server.public_ip
}

