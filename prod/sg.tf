# 自分のIPアドレスを動的に取得
data "http" "my_ip" {
  url = "https://api.ipify.org?format=json"

  request_headers = {
    Accept = "application/json"
  }
}

locals {
  my_ip = "${jsondecode(data.http.my_ip.response_body)["ip"]}/32"
}

# ----------------------
# インスタンス用
# ----------------------
resource "aws_security_group" "app" {
  name   = "${var.project}-${var.env}-sg-app"
  vpc_id = aws_vpc.main.id

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "ingress_app_ssh" {
  security_group_id = aws_security_group.app.id
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 22
  to_port           = 22
  cidr_blocks       = [local.my_ip]
}

resource "aws_security_group_rule" "egress_app_https" {
  security_group_id = aws_security_group.app.id
  type              = "egress"
  protocol          = "tcp"
  from_port         = 443
  to_port           = 443
  cidr_blocks       = ["0.0.0.0/0"]
}