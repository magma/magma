resource "aws_efs_file_system" "eks" {
  tags = {
    Name = "${var.cluster_name}-EFS"
  }
}

resource "aws_efs_mount_target" "subnet_0" {
  file_system_id  = aws_efs_file_system.eks.id
  subnet_id       = var.subnets.0
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_mount_target" "subnet_1" {
  file_system_id  = aws_efs_file_system.eks.id
  subnet_id       = var.subnets.1
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_mount_target" "subnet_2" {
  file_system_id  = aws_efs_file_system.eks.id
  subnet_id       = var.subnets.2
  security_groups = [aws_security_group.efs.id]
}
