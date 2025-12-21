terraform_binary = "tofu"

locals {
  account_vars = read_terragrunt_config(find_in_parent_folders("account.hcl"))

  region_vars = read_terragrunt_config(find_in_parent_folders("region.hcl"))

  env_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))

  account_id  = local.account_vars.locals.account_id
  aws_region  = local.region_vars.locals.aws_region
  environment = local.env_vars.locals.environment
  project     = local.env_vars.locals.project
}

generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.82"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.35"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.17"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.7"
    }
  }
}

provider "aws" {
  region = "${local.aws_region}"

  default_tags {
    tags = {
      Project     = "${local.project}"
      Environment = "${local.environment}"
      ManagedBy   = "terragrunt"
    }
  }
}
EOF
}

remote_state {
  backend = "s3"
  config = {
    encrypt        = true
    bucket         = "kin-tofu-state-${local.account_id}"
    key            = "${path_relative_to_include()}/terraform.tfstate"
    region         = local.aws_region
    dynamodb_table = "kin-tofu-locks"
  }
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
}

inputs = {
  aws_region  = local.aws_region
  environment = local.environment
  project     = local.project

  common_tags = {
    Project     = local.project
    Environment = local.environment
    ManagedBy   = "terragrunt"
  }
}
