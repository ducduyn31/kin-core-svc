# AWS Organizations Bootstrap
# This creates the AWS Organizations structure with member accounts.
#
# IMPORTANT: Run this from the MANAGEMENT account (root account).
# This is a one-time setup that creates/manages:
# - AWS Organizations
# - Organizational Units (OUs)
# - Shared Resources account
# - Kin Production account
# - Service Control Policies (SCPs)
#
# Prerequisites (ClickOps - must be created manually first):
# 1. S3 bucket: kin-mgmt-tfstate (versioning enabled, encrypted)
# 2. DynamoDB table: kin-mgmt-tf-locks (partition key: LockID)
# 3. You must be authenticated as the management account root or admin
#
# Run:
#   cd infrastructure/bootstrap/org
#   tofu init
#
#   # If AWS Organization already exists, import it first:
#   tofu import aws_organizations_organization.org $(aws organizations describe-organization --query 'Organization.Id' --output text)
#
#   tofu plan
#   tofu apply

terraform {
  required_version = ">= 1.0"

  backend "s3" {
    bucket         = "kin-mgmt-tfstate"
    key            = "bootstrap/org/terraform.tfstate"
    region         = "ap-southeast-2"
    dynamodb_table = "kin-mgmt-tf-locks"
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
      ManagedBy = "tofu-org-bootstrap"
      Project   = "kin"
    }
  }
}

# -----------------------------------------------------------------------------
# AWS Organizations
# -----------------------------------------------------------------------------
resource "aws_organizations_organization" "org" {
  aws_service_access_principals = [
    "cloudtrail.amazonaws.com",
    "config.amazonaws.com",
    "sso.amazonaws.com",
    "tagpolicies.tag.amazonaws.com",
    "account.amazonaws.com",
  ]

  enabled_policy_types = [
    "SERVICE_CONTROL_POLICY",
    "TAG_POLICY",
  ]

  feature_set = "ALL"
}

# -----------------------------------------------------------------------------
# Organizational Units
# -----------------------------------------------------------------------------
resource "aws_organizations_organizational_unit" "infrastructure" {
  name      = "Infrastructure"
  parent_id = aws_organizations_organization.org.roots[0].id
}

resource "aws_organizations_organizational_unit" "workloads" {
  name      = "Workloads"
  parent_id = aws_organizations_organization.org.roots[0].id
}

resource "aws_organizations_organizational_unit" "workloads_production" {
  name      = "Production"
  parent_id = aws_organizations_organizational_unit.workloads.id
}

# -----------------------------------------------------------------------------
# Member Accounts
# -----------------------------------------------------------------------------

# Shared Resources Account
# Used for: ECR, shared networking, DNS, etc.
resource "aws_organizations_account" "shared" {
  name      = "kin-shared"
  email     = var.shared_account_email
  parent_id = aws_organizations_organizational_unit.infrastructure.id

  role_name = "OrganizationAccountAccessRole"

  # Prevent accidental deletion
  lifecycle {
    prevent_destroy = true
  }

  tags = {
    Name        = "Kin Shared Resources"
    Environment = "shared"
  }
}

# Kin Production Account
resource "aws_organizations_account" "kin_production" {
  name      = "kin-production"
  email     = var.production_account_email
  parent_id = aws_organizations_organizational_unit.workloads_production.id

  role_name = "OrganizationAccountAccessRole"

  # Prevent accidental deletion
  lifecycle {
    prevent_destroy = true
  }

  tags = {
    Name        = "Kin Production"
    Environment = "production"
  }
}

# -----------------------------------------------------------------------------
# Service Control Policies
# -----------------------------------------------------------------------------

# Deny leaving the organization
resource "aws_organizations_policy" "deny_leave_org" {
  name        = "DenyLeaveOrganization"
  description = "Prevent accounts from leaving the organization"
  type        = "SERVICE_CONTROL_POLICY"

  content = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "DenyLeaveOrg"
        Effect    = "Deny"
        Action    = "organizations:LeaveOrganization"
        Resource  = "*"
      }
    ]
  })
}

# Deny root user actions (except essential security setup)
resource "aws_organizations_policy" "deny_root_user" {
  name        = "DenyRootUserActions"
  description = "Prevent root user from performing actions (except essential ones)"
  type        = "SERVICE_CONTROL_POLICY"

  content = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "DenyRootUser"
        Effect = "Deny"
        NotAction = [
          # Allow MFA setup
          "iam:CreateVirtualMFADevice",
          "iam:EnableMFADevice",
          "iam:GetUser",
          "iam:ListMFADevices",
          "iam:ListVirtualMFADevices",
          "iam:ResyncMFADevice",
          "iam:DeactivateMFADevice",
          "iam:DeleteVirtualMFADevice",
          # Allow password change
          "iam:GetAccountPasswordPolicy",
          "iam:ChangePassword",
          # Allow viewing account info
          "iam:GetAccountSummary",
          "iam:ListAccountAliases",
        ]
        Resource = "*"
        Condition = {
          StringLike = {
            "aws:PrincipalArn" = "arn:aws:iam::*:root"
          }
        }
      }
    ]
  })
}

# Require IMDSv2 for EC2 instances
resource "aws_organizations_policy" "require_imdsv2" {
  name        = "RequireIMDSv2"
  description = "Require IMDSv2 for EC2 instances"
  type        = "SERVICE_CONTROL_POLICY"

  content = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "RequireIMDSv2"
        Effect    = "Deny"
        Action    = "ec2:RunInstances"
        Resource  = "arn:aws:ec2:*:*:instance/*"
        Condition = {
          StringNotEquals = {
            "ec2:MetadataHttpTokens" = "required"
          }
        }
      }
    ]
  })
}

# Deny non-approved regions
resource "aws_organizations_policy" "deny_regions" {
  name        = "DenyNonApprovedRegions"
  description = "Deny actions in non-approved regions"
  type        = "SERVICE_CONTROL_POLICY"

  content = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid      = "DenyNonApprovedRegions"
        Effect   = "Deny"
        Action   = "*"
        Resource = "*"
        Condition = {
          StringNotEquals = {
            "aws:RequestedRegion" = var.allowed_regions
          }
          # Exclude global services
          "ForAllValues:StringNotLike" = {
            "aws:PrincipalArn" = [
              "arn:aws:iam::*:role/OrganizationAccountAccessRole"
            ]
          }
        }
      }
    ]
  })
}

# Attach SCPs to workloads OU
resource "aws_organizations_policy_attachment" "workloads_deny_leave" {
  policy_id = aws_organizations_policy.deny_leave_org.id
  target_id = aws_organizations_organizational_unit.workloads.id
}

resource "aws_organizations_policy_attachment" "workloads_deny_root" {
  policy_id = aws_organizations_policy.deny_root_user.id
  target_id = aws_organizations_organizational_unit.workloads.id
}

resource "aws_organizations_policy_attachment" "workloads_require_imdsv2" {
  policy_id = aws_organizations_policy.require_imdsv2.id
  target_id = aws_organizations_organizational_unit.workloads.id
}

resource "aws_organizations_policy_attachment" "workloads_deny_regions" {
  policy_id = aws_organizations_policy.deny_regions.id
  target_id = aws_organizations_organizational_unit.workloads.id
}

# Attach SCPs to infrastructure OU
resource "aws_organizations_policy_attachment" "infra_deny_leave" {
  policy_id = aws_organizations_policy.deny_leave_org.id
  target_id = aws_organizations_organizational_unit.infrastructure.id
}

resource "aws_organizations_policy_attachment" "infra_deny_regions" {
  policy_id = aws_organizations_policy.deny_regions.id
  target_id = aws_organizations_organizational_unit.infrastructure.id
}
