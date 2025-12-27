# ArgoCD configuration for production
# Installs ArgoCD and configures App-of-Apps pattern

terraform {
  source = "${dirname(find_in_parent_folders("root.hcl"))}/modules/argocd"
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

dependency "eks_addons" {
  config_path = "../eks-addons"

  mock_outputs = {
    external_secrets_installed = true
  }
}

locals {
  env_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))
  environment = local.env_vars.locals.environment
  project     = local.env_vars.locals.project
}

inputs = {
  cluster_name                       = dependency.eks.outputs.cluster_name
  cluster_endpoint                   = dependency.eks.outputs.cluster_endpoint
  cluster_certificate_authority_data = dependency.eks.outputs.cluster_certificate_authority_data

  environment = local.environment
  project     = local.project

  # GitHub repository for GitOps
  git_repo_url = "https://github.com/danielng/kin-core-svc"
  git_branch   = "main"
}
