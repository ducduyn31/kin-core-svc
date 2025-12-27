include "root" {
  path = find_in_parent_folders("root.hcl")
}

include "envcommon" {
  path   = "${dirname(find_in_parent_folders("root.hcl"))}/_envcommon/eks.hcl"
  expose = true
}

# Dependencies
dependency "vpc" {
  config_path = "../vpc"

  mock_outputs = {
    vpc_id          = "vpc-mock"
    private_subnets = ["subnet-1", "subnet-2", "subnet-3"]
    intra_subnets   = ["subnet-4", "subnet-5", "subnet-6"]
  }
}

dependency "vpc_endpoints" {
  config_path = "../vpc-endpoints"

  mock_outputs = {
    s3_endpoint_id = "vpce-mock"
  }
}

dependency "iam_roles" {
  config_path = "../iam-roles"

  mock_outputs = {
    roles = {
      "kin-production-deployment" = {
        arn  = "arn:aws:iam::123456789012:role/kin-production-deployment"
        name = "kin-production-deployment"
      }
    }
    role_arns = {
      "kin-production-deployment" = "arn:aws:iam::123456789012:role/kin-production-deployment"
    }
  }
}

# Additional inputs specific to this environment
inputs = {
  # Network configuration
  vpc_id                   = dependency.vpc.outputs.vpc_id
  subnet_ids               = dependency.vpc.outputs.private_subnets
  control_plane_subnet_ids = dependency.vpc.outputs.private_subnets
}
