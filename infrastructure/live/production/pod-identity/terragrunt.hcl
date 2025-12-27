include "root" {
  path = find_in_parent_folders("root.hcl")
}

include "envcommon" {
  path   = "${dirname(find_in_parent_folders("root.hcl"))}/_envcommon/pod-identity.hcl"
  expose = true
}

# Dependencies
dependency "eks" {
  config_path = "../eks"

  mock_outputs = {
    cluster_name = "kin-production"
  }
}

dependency "rds" {
  config_path = "../rds"

  mock_outputs = {
    db_instance_resource_id = "db-MOCKRESOURCEID"
  }
}

dependency "s3" {
  config_path = "../s3"

  mock_outputs = {
    s3_bucket_id = "kin-media-production-123456789012"
  }
}

inputs = {
  cluster_name    = dependency.eks.outputs.cluster_name
  rds_resource_id = dependency.rds.outputs.db_instance_resource_id
  s3_bucket_name  = dependency.s3.outputs.s3_bucket_id
}
