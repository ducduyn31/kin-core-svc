terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.27"
    }
  }
}

data "aws_region" "current" {}

locals {
  region = data.aws_region.current.id
}

# -----------------------------------------------------------------------------
# Security Group for Interface Endpoints
# -----------------------------------------------------------------------------
resource "aws_security_group" "vpc_endpoints" {
  name        = "${var.project}-${var.environment}-vpc-endpoints"
  description = "Security group for VPC interface endpoints"
  vpc_id      = var.vpc_id

  ingress {
    description = "HTTPS from VPC"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [var.vpc_cidr]
  }

  tags = {
    Name        = "${var.project}-${var.environment}-vpc-endpoints"
    Environment = var.environment
    Project     = var.project
  }
}

# -----------------------------------------------------------------------------
# S3 Gateway Endpoint (FREE)
# -----------------------------------------------------------------------------
resource "aws_vpc_endpoint" "s3" {
  vpc_id            = var.vpc_id
  service_name      = "com.amazonaws.${local.region}.s3"
  vpc_endpoint_type = "Gateway"
  route_table_ids   = var.private_route_table_ids

  tags = {
    Name        = "${var.project}-${var.environment}-s3"
    Environment = var.environment
    Project     = var.project
  }
}

# -----------------------------------------------------------------------------
# EC2 Endpoint (Required for EKS nodes)
# -----------------------------------------------------------------------------
resource "aws_vpc_endpoint" "ec2" {
  vpc_id              = var.vpc_id
  service_name        = "com.amazonaws.${local.region}.ec2"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = var.private_subnet_ids
  security_group_ids  = [aws_security_group.vpc_endpoints.id]
  private_dns_enabled = true

  tags = {
    Name        = "${var.project}-${var.environment}-ec2"
    Environment = var.environment
    Project     = var.project
  }
}

# -----------------------------------------------------------------------------
# ECR API Endpoint
# -----------------------------------------------------------------------------
resource "aws_vpc_endpoint" "ecr_api" {
  vpc_id              = var.vpc_id
  service_name        = "com.amazonaws.${local.region}.ecr.api"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = var.private_subnet_ids
  security_group_ids  = [aws_security_group.vpc_endpoints.id]
  private_dns_enabled = true

  tags = {
    Name        = "${var.project}-${var.environment}-ecr-api"
    Environment = var.environment
    Project     = var.project
  }
}

# -----------------------------------------------------------------------------
# ECR DKR Endpoint (Docker Registry)
# -----------------------------------------------------------------------------
resource "aws_vpc_endpoint" "ecr_dkr" {
  vpc_id              = var.vpc_id
  service_name        = "com.amazonaws.${local.region}.ecr.dkr"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = var.private_subnet_ids
  security_group_ids  = [aws_security_group.vpc_endpoints.id]
  private_dns_enabled = true

  tags = {
    Name        = "${var.project}-${var.environment}-ecr-dkr"
    Environment = var.environment
    Project     = var.project
  }
}

# -----------------------------------------------------------------------------
# Secrets Manager Endpoint
# -----------------------------------------------------------------------------
resource "aws_vpc_endpoint" "secretsmanager" {
  vpc_id              = var.vpc_id
  service_name        = "com.amazonaws.${local.region}.secretsmanager"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = var.private_subnet_ids
  security_group_ids  = [aws_security_group.vpc_endpoints.id]
  private_dns_enabled = true

  tags = {
    Name        = "${var.project}-${var.environment}-secretsmanager"
    Environment = var.environment
    Project     = var.project
  }
}

# -----------------------------------------------------------------------------
# STS Endpoint (Required for IRSA)
# -----------------------------------------------------------------------------
resource "aws_vpc_endpoint" "sts" {
  vpc_id              = var.vpc_id
  service_name        = "com.amazonaws.${local.region}.sts"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = var.private_subnet_ids
  security_group_ids  = [aws_security_group.vpc_endpoints.id]
  private_dns_enabled = true

  tags = {
    Name        = "${var.project}-${var.environment}-sts"
    Environment = var.environment
    Project     = var.project
  }
}

# -----------------------------------------------------------------------------
# CloudWatch Logs Endpoint
# -----------------------------------------------------------------------------
resource "aws_vpc_endpoint" "logs" {
  vpc_id              = var.vpc_id
  service_name        = "com.amazonaws.${local.region}.logs"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = var.private_subnet_ids
  security_group_ids  = [aws_security_group.vpc_endpoints.id]
  private_dns_enabled = true

  tags = {
    Name        = "${var.project}-${var.environment}-logs"
    Environment = var.environment
    Project     = var.project
  }
}

# -----------------------------------------------------------------------------
# X-Ray Endpoint
# -----------------------------------------------------------------------------
resource "aws_vpc_endpoint" "xray" {
  vpc_id              = var.vpc_id
  service_name        = "com.amazonaws.${local.region}.xray"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = var.private_subnet_ids
  security_group_ids  = [aws_security_group.vpc_endpoints.id]
  private_dns_enabled = true

  tags = {
    Name        = "${var.project}-${var.environment}-xray"
    Environment = var.environment
    Project     = var.project
  }
}
