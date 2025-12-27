# -----------------------------------------------------------------------------
# Shared Account - Provider and State Resources
# -----------------------------------------------------------------------------

provider "aws" {
  alias  = "shared"
  region = var.aws_region

  assume_role {
    role_arn = "arn:aws:iam::${aws_organizations_account.shared.id}:role/OrganizationAccountAccessRole"
  }

  default_tags {
    tags = {
      ManagedBy = "tofu-org-bootstrap"
      Project   = "kin"
    }
  }
}

# State Bucket
resource "aws_s3_bucket" "shared_tfstate" {
  provider = aws.shared
  bucket   = "kin-shared-tfstate"

  lifecycle {
    prevent_destroy = true
  }

  tags = {
    Name = "OpenTofu State - Shared Account"
  }
}

resource "aws_s3_bucket_versioning" "shared_tfstate" {
  provider = aws.shared
  bucket   = aws_s3_bucket.shared_tfstate.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "shared_tfstate" {
  provider = aws.shared
  bucket   = aws_s3_bucket.shared_tfstate.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "shared_tfstate" {
  provider = aws.shared
  bucket   = aws_s3_bucket.shared_tfstate.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Lock Table
resource "aws_dynamodb_table" "shared_tf_locks" {
  provider     = aws.shared
  name         = "kin-shared-tf-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = {
    Name = "OpenTofu State Locks - Shared Account"
  }
}
