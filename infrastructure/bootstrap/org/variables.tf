variable "aws_region" {
  description = "AWS region for the management account"
  type        = string
  default     = "ap-southeast-2"
}

variable "shared_account_email" {
  description = "Email address for the Shared Resources account (must be unique)"
  type        = string
}

variable "production_account_email" {
  description = "Email address for the Kin Production account (must be unique)"
  type        = string
}

variable "allowed_regions" {
  description = "List of allowed AWS regions"
  type        = list(string)
  default = [
    "ap-southeast-2", # Sydney (primary)
    "us-east-1",      # N. Virginia (for global services like CloudFront, IAM)
  ]
}
