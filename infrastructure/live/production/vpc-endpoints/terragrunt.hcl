terraform {
  source = "${dirname(find_in_parent_folders("root.hcl"))}/modules/vpc-endpoints"
}

include "root" {
  path = find_in_parent_folders("root.hcl")
}

dependency "vpc" {
  config_path = "../vpc"

  mock_outputs = {
    vpc_id                  = "vpc-mock"
    vpc_cidr_block          = "10.0.0.0/16"
    private_subnets         = ["subnet-1", "subnet-2", "subnet-3"]
    private_route_table_ids = ["rtb-1", "rtb-2", "rtb-3"]
  }
}

locals {
  env_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))

  environment = local.env_vars.locals.environment
  project     = local.env_vars.locals.project
}

inputs = {
  vpc_id                  = dependency.vpc.outputs.vpc_id
  vpc_cidr                = dependency.vpc.outputs.vpc_cidr_block
  private_subnet_ids      = dependency.vpc.outputs.private_subnets
  private_route_table_ids = dependency.vpc.outputs.private_route_table_ids

  environment = local.environment
  project     = local.project
}
