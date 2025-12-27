output "organization_id" {
  description = "AWS Organizations ID"
  value       = aws_organizations_organization.org.id
}

output "organization_arn" {
  description = "AWS Organizations ARN"
  value       = aws_organizations_organization.org.arn
}

output "management_account_id" {
  description = "Management account ID"
  value       = aws_organizations_organization.org.master_account_id
}

output "shared_account_id" {
  description = "Shared Resources account ID"
  value       = aws_organizations_account.shared.id
}

output "shared_account_arn" {
  description = "Shared Resources account ARN"
  value       = aws_organizations_account.shared.arn
}

output "production_account_id" {
  description = "Kin Production account ID"
  value       = aws_organizations_account.kin_production.id
}

output "production_account_arn" {
  description = "Kin Production account ARN"
  value       = aws_organizations_account.kin_production.arn
}

output "infrastructure_ou_id" {
  description = "Infrastructure OU ID"
  value       = aws_organizations_organizational_unit.infrastructure.id
}

output "workloads_ou_id" {
  description = "Workloads OU ID"
  value       = aws_organizations_organizational_unit.workloads.id
}

output "production_ou_id" {
  description = "Production OU ID"
  value       = aws_organizations_organizational_unit.workloads_production.id
}

# Cross-account role ARNs for assuming roles
output "shared_account_role_arn" {
  description = "Role ARN to assume in Shared account"
  value       = "arn:aws:iam::${aws_organizations_account.shared.id}:role/OrganizationAccountAccessRole"
}

output "production_account_role_arn" {
  description = "Role ARN to assume in Production account"
  value       = "arn:aws:iam::${aws_organizations_account.kin_production.id}:role/OrganizationAccountAccessRole"
}

# State bucket and lock table outputs
output "shared_tfstate_bucket" {
  description = "S3 bucket for Shared account OpenTofu state"
  value       = aws_s3_bucket.shared_tfstate.id
}

output "shared_tf_locks_table" {
  description = "DynamoDB table for Shared account OpenTofu state locking"
  value       = aws_dynamodb_table.shared_tf_locks.id
}

output "production_tfstate_bucket" {
  description = "S3 bucket for Production account OpenTofu state"
  value       = aws_s3_bucket.production_tfstate.id
}

output "production_tf_locks_table" {
  description = "DynamoDB table for Production account OpenTofu state locking"
  value       = aws_dynamodb_table.production_tf_locks.id
}

# Identity Center outputs
output "identity_center_instance_arn" {
  description = "IAM Identity Center instance ARN"
  value       = local.sso_instance_arn
}

output "identity_store_id" {
  description = "IAM Identity Store ID"
  value       = local.identity_store_id
}

output "admins_group_id" {
  description = "Admins group ID in Identity Store"
  value       = aws_identitystore_group.admins.group_id
}

output "permission_set_administrator_access_arn" {
  description = "AdministratorAccess permission set ARN"
  value       = aws_ssoadmin_permission_set.administrator_access.arn
}

output "permission_set_view_only_access_arn" {
  description = "ViewOnlyAccess permission set ARN"
  value       = aws_ssoadmin_permission_set.view_only_access.arn
}
