terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.27"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.7"
    }
  }
}

# -----------------------------------------------------------------------------
# Auth Token Generation
# -----------------------------------------------------------------------------
resource "random_password" "auth_token" {
  length  = 32
  special = false # ElastiCache auth tokens don't support all special characters
}

# -----------------------------------------------------------------------------
# Auth Token Secret in Secrets Manager
# -----------------------------------------------------------------------------
resource "aws_secretsmanager_secret" "auth_token" {
  name        = "${var.project}/${var.environment}/elasticache-auth"
  description = "ElastiCache auth token for ${var.project} ${var.environment}"

  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "auth_token" {
  secret_id     = aws_secretsmanager_secret.auth_token.id
  secret_string = random_password.auth_token.result
}

# -----------------------------------------------------------------------------
# ElastiCache Cluster
# -----------------------------------------------------------------------------
module "elasticache" {
  source  = "terraform-aws-modules/elasticache/aws"
  version = "1.10.3"

  # Use replication group (module default) - more flexible than cluster
  replication_group_id = var.cluster_id

  # Engine configuration
  engine         = var.engine
  engine_version = var.engine_version
  node_type      = var.node_type

  # Replication group configuration
  num_cache_clusters = var.num_cache_nodes

  # Parameter group
  parameter_group_family = var.parameter_group_family
  parameters             = var.parameters

  # Network
  subnet_ids            = var.subnet_ids
  security_group_ids    = var.security_group_ids
  create_security_group = false  # Use provided security groups, don't create new one

  # Security - enable auth token with transit encryption
  at_rest_encryption_enabled = var.at_rest_encryption_enabled
  transit_encryption_enabled = var.transit_encryption_enabled
  auth_token                 = random_password.auth_token.result

  # Maintenance
  maintenance_window       = var.maintenance_window
  snapshot_window          = var.snapshot_window
  snapshot_retention_limit = var.snapshot_retention_limit

  # Auto minor version upgrade
  auto_minor_version_upgrade = var.auto_minor_version_upgrade

  tags = var.tags
}
