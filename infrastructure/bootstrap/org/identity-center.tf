# -----------------------------------------------------------------------------
# IAM Identity Center (AWS SSO)
# IAM Identity Center is enabled in us-east-1
# -----------------------------------------------------------------------------

# Get the existing IAM Identity Center instance
data "aws_ssoadmin_instances" "this" {
  provider = aws.us_east_1
}

locals {
  identity_store_id = (
    length(data.aws_ssoadmin_instances.this.identity_store_ids) > 0
    ? tolist(data.aws_ssoadmin_instances.this.identity_store_ids)[0]
    : error("AWS Identity Center not enabled or no instances found")
  )
  sso_instance_arn = (
    length(data.aws_ssoadmin_instances.this.arns) > 0
    ? tolist(data.aws_ssoadmin_instances.this.arns)[0]
    : error("AWS Identity Center not enabled or no instances found")
  )
}

# -----------------------------------------------------------------------------
# Groups
# -----------------------------------------------------------------------------

resource "aws_identitystore_group" "admins" {
  provider = aws.us_east_1

  identity_store_id = local.identity_store_id
  display_name      = "Admins"
  description       = "Administrators with full access to all accounts"
}

# -----------------------------------------------------------------------------
# Account Assignments
# -----------------------------------------------------------------------------

# Admins group - AdministratorAccess on Shared account
resource "aws_ssoadmin_account_assignment" "admins_admin_shared" {
  provider = aws.us_east_1

  instance_arn       = local.sso_instance_arn
  permission_set_arn = aws_ssoadmin_permission_set.administrator_access.arn

  principal_id   = aws_identitystore_group.admins.group_id
  principal_type = "GROUP"

  target_id   = aws_organizations_account.shared.id
  target_type = "AWS_ACCOUNT"
}

# Admins group - AdministratorAccess on Production account
resource "aws_ssoadmin_account_assignment" "admins_admin_production" {
  provider = aws.us_east_1

  instance_arn       = local.sso_instance_arn
  permission_set_arn = aws_ssoadmin_permission_set.administrator_access.arn

  principal_id   = aws_identitystore_group.admins.group_id
  principal_type = "GROUP"

  target_id   = aws_organizations_account.kin_production.id
  target_type = "AWS_ACCOUNT"
}

# Admins group - ViewOnlyAccess on Shared account
resource "aws_ssoadmin_account_assignment" "admins_viewonly_shared" {
  provider = aws.us_east_1

  instance_arn       = local.sso_instance_arn
  permission_set_arn = aws_ssoadmin_permission_set.view_only_access.arn

  principal_id   = aws_identitystore_group.admins.group_id
  principal_type = "GROUP"

  target_id   = aws_organizations_account.shared.id
  target_type = "AWS_ACCOUNT"
}

# Admins group - ViewOnlyAccess on Production account
resource "aws_ssoadmin_account_assignment" "admins_viewonly_production" {
  provider = aws.us_east_1

  instance_arn       = local.sso_instance_arn
  permission_set_arn = aws_ssoadmin_permission_set.view_only_access.arn

  principal_id   = aws_identitystore_group.admins.group_id
  principal_type = "GROUP"

  target_id   = aws_organizations_account.kin_production.id
  target_type = "AWS_ACCOUNT"
}
