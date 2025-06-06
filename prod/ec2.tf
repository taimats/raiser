# ----------------------
# Key Pair
# ----------------------
resource "aws_key_pair" "app" {
  key_name   = "${var.project}-${var.env}-key-pair-app"
  public_key = file("./src/bhapi-key-pair-app.pub")
}

# ----------------------
# インスタンス
# ----------------------
resource "aws_instance" "app" {
  ami                         = "ami-0287977b9e404b70f"
  instance_type               = "t2.micro"
  subnet_id                   = aws_subnet.public_01.id
  associate_public_ip_address = true
  vpc_security_group_ids      = [aws_security_group.app.id]
  key_name                    = aws_key_pair.app.key_name
  user_data                   = file("./src/init.sh")

  tags = {
    Name = "${var.project}-${var.env}-ec2-app"
  }
}