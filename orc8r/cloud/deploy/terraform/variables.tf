variable "region" {
  default = "eu-west-1"
}

data "aws_availability_zones" "available" {}

variable "vpc_name" {
  default = "orc8r-vpc"
}

variable "cluster_name" {
  default = "orc8r"
}

variable "db_password" {
  description = "Password for the DB user. You should put this value into a file NOT checked into source control!"
  type        = "string"
}

variable "key_name" {
  description = "Name of the EC2 Keypair to use for SSH access to nodes"
  type        = "string"
}

variable "map_users" {
  description = "Additional IAM users to add to the aws-auth configmap"
  type        = "list"
  default     = []

  # For e.g.:
  # [
  #   {
  #     user_arn = "arn:aws:iam::66666666666:user/user1"
  #     username = "user1"
  #     group    = "system:masters"
  #   },
  # ]
}
