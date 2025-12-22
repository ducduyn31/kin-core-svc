# EKS Add-ons configuration for production
# Installs External Secrets Operator, AWS Load Balancer Controller, OTEL Collector
# Uses EKS Pod Identity (not IRSA) for IAM permissions

terraform {
  source = "${dirname(find_in_parent_folders("root.hcl"))}/modules/eks-addons"
}

include "root" {
  path = find_in_parent_folders("root.hcl")
}

# Dependencies
dependency "eks" {
  config_path = "../eks"

  mock_outputs = {
    cluster_name                       = "kin-production"
    cluster_endpoint                   = "https://mock.eks.amazonaws.com"
    cluster_certificate_authority_data = "bW9jaw=="
  }
}

dependency "vpc" {
  config_path = "../vpc"

  mock_outputs = {
    vpc_id = "vpc-mock"
  }
}

locals {
  env_vars     = read_terragrunt_config(find_in_parent_folders("env.hcl"))
  account_vars = read_terragrunt_config(find_in_parent_folders("account.hcl"))
  region_vars  = read_terragrunt_config(find_in_parent_folders("region.hcl"))

  environment = local.env_vars.locals.environment
  project     = local.env_vars.locals.project
  account_id  = local.account_vars.locals.account_id
  aws_region  = local.region_vars.locals.aws_region
}

inputs = {
  cluster_name                       = dependency.eks.outputs.cluster_name
  cluster_endpoint                   = dependency.eks.outputs.cluster_endpoint
  cluster_certificate_authority_data = dependency.eks.outputs.cluster_certificate_authority_data
  vpc_id                             = dependency.vpc.outputs.vpc_id

  environment = local.environment
  project     = local.project
  aws_region  = local.aws_region
}
