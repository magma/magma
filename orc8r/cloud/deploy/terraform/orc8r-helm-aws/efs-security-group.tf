resource "aws_security_group" "efs" {
  name        = "efs-sg"
  description = "Allows inbound EFS traffic"
  vpc_id      = var.vpc_id

  ingress {
    from_port       = 2049
    to_port         = 2049
    protocol        = "tcp"
    security_groups = [var.node_security_group_id]
  }
}

resource "aws_security_group_rule" "efs" {
  description              = "Allow outbound EFS traffic"
  type                     = "egress"
  from_port                = 2049
  to_port                  = 2049
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.efs.id
  security_group_id        = var.node_security_group_id
}
