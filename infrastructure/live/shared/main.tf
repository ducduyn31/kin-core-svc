# Shared Resources Account
# This creates shared resources used across all environments:
# - ECR repositories (container images)
# - GitHub Actions OIDC for CI/CD
# - Cross-account access policies
#
# Prerequisites:
# - Run org bootstrap first (creates this account and its state bucket)
#
# Run:
#   cd infrastructure/live/shared
#
#   # Assume role into the shared account
#   export AWS_PROFILE=kin-shared  # or use aws sts assume-role
#
#   cp terraform.tfvars.example terraform.tfvars
#   # Edit terraform.tfvars with your account IDs
#
#   tofu init
#   tofu apply

terraform {
  required_version = ">= 1.0"

  backend "s3" {
    bucket         = "kin-shared-tfstate"
    key            = "shared/terraform.tfstate"
    region         = "ap-southeast-2"
    dynamodb_table = "kin-shared-tf-locks"
    encrypt        = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.82"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      ManagedBy   = "opentofu"
      Project     = "kin"
      Environment = "shared"
    }
  }
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

locals {
  account_id = data.aws_caller_identity.current.account_id
  region     = data.aws_region.current.name
}


# -----------------------------------------------------------------------------
# GitHub Actions OIDC Provider (for pushing to ECR)
# -----------------------------------------------------------------------------
resource "aws_iam_openid_connect_provider" "github_actions" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = ["sts.amazonaws.com"]

  thumbprint_list = [
    "6938fd4d98bab03faadb97b34396831e3780aea1",
    "1c58a3a8518e8759bf075b76b750d4f2df264fcd"
  ]

  tags = {
    Name = "GitHub Actions OIDC"
  }
}

# IAM Role for GitHub Actions to push to ECR
data "aws_iam_policy_document" "github_actions_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github_actions.arn]
    }

    actions = ["sts:AssumeRoleWithWebIdentity"]

    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }

    condition {
      test     = "StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values   = ["repo:${var.github_org}/${var.github_repo}:*"]
    }
  }
}

resource "aws_iam_role" "github_actions" {
  name               = "github-actions-ecr"
  assume_role_policy = data.aws_iam_policy_document.github_actions_assume_role.json

  tags = {
    Name = "GitHub Actions ECR Role"
  }
}

data "aws_iam_policy_document" "github_actions_ecr" {
  statement {
    effect    = "Allow"
    actions   = ["ecr:GetAuthorizationToken"]
    resources = ["*"]
  }

  statement {
    effect = "Allow"
    actions = [
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage",
      "ecr:PutImage",
      "ecr:InitiateLayerUpload",
      "ecr:UploadLayerPart",
      "ecr:CompleteLayerUpload"
    ]
    resources = [for repo in aws_ecr_repository.repos : repo.arn]
  }
}

resource "aws_iam_role_policy" "github_actions_ecr" {
  name   = "ecr-push"
  role   = aws_iam_role.github_actions.id
  policy = data.aws_iam_policy_document.github_actions_ecr.json
}
