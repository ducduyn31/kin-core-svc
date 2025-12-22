terraform {
  source = "${dirname(find_in_parent_folders("root.hcl"))}/modules/elasticache"
}

locals {
  env_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))

  environment = local.env_vars.locals.environment
  project     = local.env_vars.locals.project

  cluster_id = "${local.project}-${local.environment}"
}

inputs = {
  # Required for auth token secret naming
  project     = local.project
  environment = local.environment

  cluster_id = local.cluster_id

  # Engine configuration
  engine         = "valkey"
  engine_version = "8.2"
  node_type      = "cache.t3.micro"

  # Cluster configuration - single node for cost optimization
  num_cache_nodes = 1

  # Parameter group
  parameter_group_family = "valkey8"
  parameters = [
    {
      name  = "maxmemory-policy"
      value = "allkeys-lru"
    }
  ]

  # Security
  at_rest_encryption_enabled = true
  transit_encryption_enabled = true

  # Maintenance
  maintenance_window = "sun:05:00-sun:09:00"
  snapshot_window    = "00:00-04:00"
  snapshot_retention_limit = 7

  # Auto minor version upgrade
  auto_minor_version_upgrade = true

  tags = {
    Environment = local.environment
    Project     = local.project
    ManagedBy   = "terragrunt"
  }
}
