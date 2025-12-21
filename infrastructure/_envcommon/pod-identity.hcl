terraform {
  source = "${dirname(find_in_parent_folders())}//modules/pod-identity"
}

locals {
  env_vars    = read_terragrunt_config(find_in_parent_folders("env.hcl"))
  region_vars = read_terragrunt_config(find_in_parent_folders("region.hcl"))

  environment = local.env_vars.locals.environment
  project     = local.env_vars.locals.project
  aws_region  = local.region_vars.locals.aws_region
}

inputs = {
  environment = local.environment
  project     = local.project
  aws_region  = local.aws_region

  namespace            = "kin"
  service_account_name = "kin-core-svc"

  # Database user for IAM authentication
  # This user must be created in PostgreSQL with: GRANT rds_iam TO core_svc;
  rds_iam_user = "core_svc"
}
