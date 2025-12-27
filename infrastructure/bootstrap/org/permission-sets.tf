# -----------------------------------------------------------------------------
# IAM Identity Center Permission Sets
# IAM Identity Center is enabled in us-east-1
# -----------------------------------------------------------------------------

# AdministratorAccess permission set
resource "aws_ssoadmin_permission_set" "administrator_access" {
  provider = aws.us_east_1

  name             = "AdministratorAccess"
  description      = "Full administrator access using AWS managed policy"
  instance_arn     = local.sso_instance_arn
  session_duration = "PT8H"

  tags = {
    Name = "AdministratorAccess"
  }
}

resource "aws_ssoadmin_managed_policy_attachment" "administrator_access" {
  provider = aws.us_east_1

  instance_arn       = local.sso_instance_arn
  managed_policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
  permission_set_arn = aws_ssoadmin_permission_set.administrator_access.arn
}

# ViewOnlyAccess permission set
resource "aws_ssoadmin_permission_set" "view_only_access" {
  provider = aws.us_east_1

  name             = "ViewOnlyAccess"
  description      = "View-only access using AWS managed policy"
  instance_arn     = local.sso_instance_arn
  session_duration = "PT8H"

  tags = {
    Name = "ViewOnlyAccess"
  }
}

resource "aws_ssoadmin_managed_policy_attachment" "view_only_access" {
  provider = aws.us_east_1

  instance_arn       = local.sso_instance_arn
  managed_policy_arn = "arn:aws:iam::aws:policy/job-function/ViewOnlyAccess"
  permission_set_arn = aws_ssoadmin_permission_set.view_only_access.arn
}
