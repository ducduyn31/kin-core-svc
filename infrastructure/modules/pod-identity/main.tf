terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.27"
    }
  }
}

data "aws_caller_identity" "current" {}

locals {
  account_id = data.aws_caller_identity.current.account_id
  role_name  = "${var.project}-${var.environment}-kin-core-svc"
}

# -----------------------------------------------------------------------------
# IAM Role for Pod Identity
# -----------------------------------------------------------------------------
data "aws_iam_policy_document" "pod_identity_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["pods.eks.amazonaws.com"]
    }

    actions = [
      "sts:AssumeRole",
      "sts:TagSession"
    ]
  }
}

resource "aws_iam_role" "kin_core_svc" {
  name               = local.role_name
  assume_role_policy = data.aws_iam_policy_document.pod_identity_assume_role.json

  tags = {
    Name        = local.role_name
    Environment = var.environment
    Project     = var.project
  }
}

# -----------------------------------------------------------------------------
# S3 Access Policy
# -----------------------------------------------------------------------------
data "aws_iam_policy_document" "s3_access" {
  statement {
    effect = "Allow"
    actions = [
      "s3:PutObject",
      "s3:GetObject",
      "s3:DeleteObject",
      "s3:ListBucket"
    ]
    resources = [
      "arn:aws:s3:::${var.s3_bucket_name}",
      "arn:aws:s3:::${var.s3_bucket_name}/*"
    ]
  }
}

resource "aws_iam_role_policy" "s3_access" {
  name   = "s3-access"
  role   = aws_iam_role.kin_core_svc.id
  policy = data.aws_iam_policy_document.s3_access.json
}

# -----------------------------------------------------------------------------
# Secrets Manager Access Policy
# -----------------------------------------------------------------------------
data "aws_iam_policy_document" "secrets_read" {
  statement {
    effect = "Allow"
    actions = [
      "secretsmanager:GetSecretValue",
      "secretsmanager:DescribeSecret"
    ]
    resources = [
      "arn:aws:secretsmanager:${var.aws_region}:${local.account_id}:secret:${var.project}/${var.environment}/*"
    ]
  }
}

resource "aws_iam_role_policy" "secrets_read" {
  name   = "secrets-read"
  role   = aws_iam_role.kin_core_svc.id
  policy = data.aws_iam_policy_document.secrets_read.json
}

# -----------------------------------------------------------------------------
# RDS IAM Authentication Policy
# -----------------------------------------------------------------------------
data "aws_iam_policy_document" "rds_connect" {
  statement {
    effect  = "Allow"
    actions = ["rds-db:connect"]
    resources = [
      "arn:aws:rds-db:${var.aws_region}:${local.account_id}:dbuser:${var.rds_resource_id}/${var.rds_iam_user}"
    ]
  }
}

resource "aws_iam_role_policy" "rds_connect" {
  name   = "rds-connect"
  role   = aws_iam_role.kin_core_svc.id
  policy = data.aws_iam_policy_document.rds_connect.json
}

# -----------------------------------------------------------------------------
# Pod Identity Association
# -----------------------------------------------------------------------------
resource "aws_eks_pod_identity_association" "kin_core_svc" {
  cluster_name    = var.cluster_name
  namespace       = var.namespace
  service_account = var.service_account_name
  role_arn        = aws_iam_role.kin_core_svc.arn

  tags = {
    Name        = "${local.role_name}-pod-identity"
    Environment = var.environment
    Project     = var.project
  }
}
