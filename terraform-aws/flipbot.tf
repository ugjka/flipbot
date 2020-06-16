provider "aws" {
  profile = "default"
  region  = "us-east-1"
}

resource "aws_vpc" "main" {
    cidr_block = "10.1.0.0/16"
    tags = {
        Name = "The FLIPBOT VPC"
    }
}

resource "aws_internet_gateway" "main" {
    vpc_id = aws_vpc.main.id
}

resource "aws_subnet" "main" {
    vpc_id = aws_vpc.main.id
    cidr_block = "10.1.1.0/24"
    availability_zone = "us-east-1a"
}

resource "aws_route_table" "default" {
    vpc_id = aws_vpc.main.id
    route {
        cidr_block = "0.0.0.0/0"
        gateway_id = aws_internet_gateway.main.id
    }
}

resource "aws_route_table_association" "main" {
    subnet_id = aws_subnet.main.id
    route_table_id = aws_route_table.default.id
}

resource "aws_network_acl" "allowall" {
    vpc_id = aws_vpc.main.id

    egress {
        protocol = "-1"
        rule_no = 100
        action = "allow"
        cidr_block = "0.0.0.0/0"
        from_port = 0
        to_port = 0
    }

    ingress {
        protocol = "-1"
        rule_no = 200
        action = "allow"
        cidr_block = "0.0.0.0/0"
        from_port = 0
        to_port = 0
    }
}

resource "aws_security_group" "flipbot" {
    name = "Flipbot Security Group"
    description = "Security rules for flipbot vpc"
    vpc_id = aws_vpc.main.id
}

resource "aws_security_group_rule" "ssh" {
    type = "ingress"
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    security_group_id = aws_security_group.flipbot.id
}

resource "aws_security_group_rule" "mosh" {
    type = "ingress"
    from_port = 60000
    to_port = 61000
    protocol = "udp"
    cidr_blocks = ["0.0.0.0/0"]
    security_group_id = aws_security_group.flipbot.id
}

resource "aws_security_group_rule" "icmp" {
    type = "ingress"
    protocol = "icmp"
    from_port = -1
    to_port = -1
    cidr_blocks = ["0.0.0.0/0"]
    security_group_id = aws_security_group.flipbot.id

}

resource "aws_security_group_rule" "egress" {
    type = "egress"
    from_port = 0
    to_port = 65535
    protocol = "all"
    cidr_blocks = ["0.0.0.0/0"]
    security_group_id = aws_security_group.flipbot.id

}

resource "aws_eip" "flipbot" {
    instance = aws_instance.flipbot.id
    vpc = true
    depends_on = [aws_internet_gateway.main]
}

resource "aws_key_pair" "default" {
    key_name = "flibpot"
    public_key = file("~/.ssh/id_rsa.pub")
}

resource "aws_ebs_volume" "flipbot_data" {
    availability_zone = "us-east-1a"
    size = 3
}

resource "aws_volume_attachment" "fl_data_att" {
    device_name = "/dev/xvdb"
    volume_id = aws_ebs_volume.flipbot_data.id
    instance_id = aws_instance.flipbot.id
}

resource "aws_instance" "flipbot" {
  ami           = "ami-080a54665e2d19667"
  instance_type = "t2.micro"
  availability_zone = "us-east-1a"
  key_name = aws_key_pair.default.key_name
  vpc_security_group_ids = [aws_security_group.flipbot.id]
  subnet_id = aws_subnet.main.id
}

output "public_ip" {
    value = aws_eip.flipbot.public_ip
}