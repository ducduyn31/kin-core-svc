terraform {
  source = "${dirname(find_in_parent_folders("root.hcl"))}//modules/iam-roles"
}

locals {
  env_vars     = read_terragrunt_config(find_in_parent_folders("env.hcl"))
  account_vars = read_terragrunt_config(find_in_parent_folders("account.hcl"))
  region_vars  = read_terragrunt_config(find_in_parent_folders("region.hcl"))

  environment = local.env_vars.locals.environment
  project     = local.env_vars.locals.project
  account_id  = local.account_vars.locals.account_id
  region      = local.region_vars.locals.aws_region

  # SSO role ARN pattern for AdministratorAccess
  sso_admin_role_arn = "arn:aws:iam::${local.account_id}:role/aws-reserved/sso.amazonaws.com/${local.region}/AWSReservedSSO_AdministratorAccess_*"
}

inputs = {
  account_id = local.account_id

  tags = {
    Environment = local.environment
    Project     = local.project
    ManagedBy   = "terragrunt"
  }
}
