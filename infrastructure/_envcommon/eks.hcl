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

  # EKS Managed Add-ons
  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
    }
    eks-pod-identity-agent = {
      most_recent = true
    }
  }

  # Managed node groups
  eks_managed_node_groups = {
    default = {
      name               = "default"
      instance_types     = ["t4g.medium"]
      kubernetes_version = "1.34"

      min_size     = 2
      max_size     = 5
      desired_size = 2

      # Use latest EKS optimized AMI
      ami_type = "AL2023_ARM_64_STANDARD"

      # IMDS configuration - hop limit must be 2 for containers
      metadata_options = {
        http_endpoint               = "enabled"
        http_tokens                 = "required"
        http_put_response_hop_limit = 2
        instance_metadata_tags      = "disabled"
      }

      # We will setup with Prometheus and Grafana for monitoring later
      enable_monitoring = false

      # Disk configuration
      block_device_mappings = {
        xvda = {
          device_name = "/dev/xvda"
          ebs = {
            volume_size           = 50
            volume_type           = "gp3"
            encrypted             = true
            delete_on_termination = true
          }
        }
      }

      labels = {
        Environment = local.environment
        NodeGroup   = "default"
      }

      tags = {
        Environment = local.environment
        Project     = local.project
      }
    }
  }

  # Cluster access
  enable_cluster_creator_admin_permissions = true

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
