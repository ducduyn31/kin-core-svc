terraform {
  source = "tfr:///terraform-aws-modules/eks/aws?version=21.8.0"
}

locals {
  env_vars     = read_terragrunt_config(find_in_parent_folders("env.hcl"))
  account_vars = read_terragrunt_config(find_in_parent_folders("account.hcl"))

  environment = local.env_vars.locals.environment
  project     = local.env_vars.locals.project
  account_id  = local.account_vars.locals.account_id

  cluster_name = "${local.project}-${local.environment}"
}

inputs = {
  name            = local.cluster_name
  cluster_version = "1.34"

  # Cluster endpoint access
  cluster_endpoint_public_access  = true
  cluster_endpoint_private_access = true

  enable_irsa = false

  # Only DaemonSet addons here - they must exist before nodes join
  # coredns is created in eks-node-groups module (needs nodes to schedule)
  addons = {
    kube-proxy = {
      most_recent    = true
      before_compute = true
    }
    vpc-cni = {
      most_recent    = true
      before_compute = true
    }
    eks-pod-identity-agent = {
      most_recent    = true
      before_compute = true
    }
  }

  eks_managed_node_groups = {}

  # Cluster access
  enable_cluster_creator_admin_permissions = true

  # Access entries for additional roles
  access_entries = {
    deployment = {
      principal_arn = "arn:aws:iam::${local.account_id}:role/${local.project}-${local.environment}-deployment"
      policy_associations = {
        cluster_admin = {
          policy_arn = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"
          access_scope = {
            type = "cluster"
          }
        }
      }
    }
  }

  # Encryption
  cluster_encryption_config = {
    resources = ["secrets"]
  }

  tags = {
    Environment = local.environment
    Project     = local.project
    ManagedBy   = "terragrunt"
  }
}
