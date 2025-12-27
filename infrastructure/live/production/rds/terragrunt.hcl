include "root" {
  path = find_in_parent_folders("root.hcl")
}

include "envcommon" {
  path   = "${dirname(find_in_parent_folders("root.hcl"))}/_envcommon/rds.hcl"
  expose = true
}

# Dependencies
dependency "vpc" {
  config_path = "../vpc"

  mock_outputs = {
    vpc_id                   = "vpc-mock"
    database_subnet_group_name = "mock-db-subnet-group"
    private_subnets          = ["subnet-1", "subnet-2", "subnet-3"]
    default_security_group_id = "sg-mock"
  }
}

inputs = {
  # Network configuration
  db_subnet_group_name   = dependency.vpc.outputs.database_subnet_group_name
  vpc_security_group_ids = [dependency.vpc.outputs.default_security_group_id]

  # Create security group for RDS
  create_db_subnet_group = false # Use the one from VPC module
}
