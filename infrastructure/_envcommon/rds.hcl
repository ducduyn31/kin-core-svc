terraform {
  source = "tfr:///terraform-aws-modules/rds/aws?version=6.13.1"
}

locals {
  env_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))

  environment = local.env_vars.locals.environment
  project     = local.env_vars.locals.project

  identifier = "${local.project}-${local.environment}"
}

inputs = {
  identifier = local.identifier

  # Engine configuration
  engine               = "postgres"
  engine_version       = "18.1"
  family               = "postgres18"
  major_engine_version = "18"

  # Instance configuration
  instance_class        = "db.t3.medium"
  allocated_storage     = 20
  max_allocated_storage = 100
  storage_type          = "gp3"
  storage_encrypted     = true

  # Database configuration
  db_name  = local.project
  username = "postgres"
  port     = 5432

  # Network security
  # Note: db_subnet_group_name and vpc_security_group_ids must be provided
  # by environment-specific config (e.g., infrastructure/live/production/rds/terragrunt.hcl)
  publicly_accessible = false

  # IAM database authentication (used with EKS Pod Identity)
  # App generates short-lived tokens instead of using passwords
  iam_database_authentication_enabled = true

  # Keep Secrets Manager password for admin/migration tasks
  manage_master_user_password = true

  # Multi-AZ - disabled for cost, enable for production HA
  multi_az = false

  # Maintenance and backup
  maintenance_window      = "Mon:00:00-Mon:03:00"
  backup_window           = "03:00-06:00"
  backup_retention_period = 7

  # Enhanced monitoring
  monitoring_interval                   = 60
  monitoring_role_name                  = "${local.identifier}-rds-monitoring"
  create_monitoring_role                = true
  enabled_cloudwatch_logs_exports       = ["postgresql", "upgrade"]
  performance_insights_enabled          = true
  performance_insights_retention_period = 7

  # Deletion protection
  deletion_protection = true
  skip_final_snapshot = false
  final_snapshot_identifier_prefix = "${local.identifier}-final"

  # Parameter group
  create_db_parameter_group = true
  parameters = [
    {
      name  = "log_min_duration_statement"
      value = "1000" # Log queries taking more than 1s
    }
  ]

  tags = {
    Environment = local.environment
    Project     = local.project
    ManagedBy   = "terragrunt"
  }
}
