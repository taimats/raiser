# ----------------------
# Route Public
# ----------------------
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.project}-${var.env}-route"
  }
}

resource "aws_route_table_association" "public_01" {
  route_table_id = aws_route_table.public.id
  subnet_id      = aws_subnet.public_01.id
}

resource "aws_route" "igw" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.main.id
}