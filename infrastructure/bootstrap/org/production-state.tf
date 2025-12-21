# -----------------------------------------------------------------------------
# Production Account - Provider and State Resources
# -----------------------------------------------------------------------------

provider "aws" {
  alias  = "production"
  region = var.aws_region

  assume_role {
    role_arn = "arn:aws:iam::${aws_organizations_account.kin_production.id}:role/OrganizationAccountAccessRole"
  }

  default_tags {
    tags = {
      ManagedBy = "tofu-org-bootstrap"
      Project   = "kin"
    }
  }
}

# State Bucket
resource "aws_s3_bucket" "production_tfstate" {
  provider = aws.production
  bucket   = "kin-production-tfstate"

  lifecycle {
    prevent_destroy = true
  }

  tags = {
    Name = "OpenTofu State - Production Account"
  }
}

resource "aws_s3_bucket_versioning" "production_tfstate" {
  provider = aws.production
  bucket   = aws_s3_bucket.production_tfstate.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "production_tfstate" {
  provider = aws.production
  bucket   = aws_s3_bucket.production_tfstate.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "production_tfstate" {
  provider = aws.production
  bucket   = aws_s3_bucket.production_tfstate.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Lock Table
resource "aws_dynamodb_table" "production_tf_locks" {
  provider     = aws.production
  name         = "kin-production-tf-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = {
    Name = "OpenTofu State Locks - Production Account"
  }
}
