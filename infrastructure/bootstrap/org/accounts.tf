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
