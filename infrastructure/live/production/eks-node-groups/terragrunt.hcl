include "root" {
  path = find_in_parent_folders("root.hcl")
}

include "envcommon" {
  path   = "${dirname(find_in_parent_folders("root.hcl"))}/_envcommon/eks-node-groups.hcl"
  expose = true
}

# Dependencies - ensures addons are created before node groups
dependency "eks" {
  config_path = "../eks"

  mock_outputs = {
    cluster_name                       = "kin-production"
    cluster_endpoint                   = "https://mock.eks.amazonaws.com"
    cluster_certificate_authority_data = "bW9jaw=="
    cluster_primary_security_group_id  = "sg-mock"
    node_security_group_id             = "sg-mock-node"
    cluster_service_cidr               = "172.20.0.0/16"
  }
}

dependency "vpc" {
  config_path = "../vpc"

  mock_outputs = {
    vpc_id          = "vpc-mock"
    private_subnets = ["subnet-1", "subnet-2", "subnet-3"]
  }
}

inputs = {
  # Cluster configuration from eks module
  cluster_name    = dependency.eks.outputs.cluster_name
  cluster_version = "1.34"

  # Network configuration
  subnet_ids = dependency.vpc.outputs.private_subnets

  # Security groups
  cluster_primary_security_group_id = dependency.eks.outputs.cluster_primary_security_group_id

  # Service CIDR for node bootstrap
  cluster_service_cidr = dependency.eks.outputs.cluster_service_cidr
}
