terraform {
  source = "tfr:///terraform-aws-modules/s3-bucket/aws?version=5.9.1"
}

locals {
  env_vars     = read_terragrunt_config(find_in_parent_folders("env.hcl"))
  account_vars = read_terragrunt_config(find_in_parent_folders("account.hcl"))

  environment = local.env_vars.locals.environment
  project     = local.env_vars.locals.project
  account_id  = local.account_vars.locals.account_id

  bucket_name = "${local.project}-media-${local.environment}-${local.account_id}"
}

inputs = {
  bucket = local.bucket_name

  # Block public access
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true

  # Versioning
  versioning = {
    enabled = true
  }

  # Server-side encryption
  server_side_encryption_configuration = {
    rule = {
      apply_server_side_encryption_by_default = {
        sse_algorithm = "AES256"
      }
    }
  }

  # CORS configuration for browser uploads
  cors_rule = [
    {
      allowed_headers = ["*"]
      allowed_methods = ["GET", "PUT", "POST"]
      allowed_origins = ["https://kin.coffeewithegg.com"]
      expose_headers  = ["ETag"]
      max_age_seconds = 3600
    }
  ]

  # Lifecycle rules
  lifecycle_rule = [
    {
      id      = "cleanup-incomplete-uploads"
      enabled = true

      abort_incomplete_multipart_upload_days = 7
    }
  ]

  tags = {
    Environment = local.environment
    Project     = local.project
    ManagedBy   = "terragrunt"
  }
}
