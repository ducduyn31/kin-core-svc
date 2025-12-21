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
