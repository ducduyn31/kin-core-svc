# ElastiCache Redis configuration for production
# Creates managed Redis cluster in private subnets

include "root" {
  path = find_in_parent_folders("root.hcl")
}

include "envcommon" {
  path   = "${dirname(find_in_parent_folders("root.hcl"))}/_envcommon/elasticache.hcl"
  expose = true
}

# Dependencies
dependency "vpc" {
  config_path = "../vpc"

  mock_outputs = {
    vpc_id                    = "vpc-mock"
    private_subnets           = ["subnet-1", "subnet-2", "subnet-3"]
    default_security_group_id = "sg-mock"
    elasticache_subnet_group_name = "mock-cache-subnet-group"
  }
}

# Additional inputs specific to this environment
inputs = {
  # Network configuration - use private subnets
  subnet_ids         = dependency.vpc.outputs.private_subnets
  security_group_ids = [dependency.vpc.outputs.default_security_group_id]
}
