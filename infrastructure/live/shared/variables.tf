variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-southeast-2"
}

variable "environment_account_ids" {
  description = "AWS account IDs for environments that can pull from ECR"
  type        = map(string)
  # Example:
  # {
  #   production = "123456789012"
  #   staging    = "234567890123"
  # }
}

variable "github_org" {
  description = "GitHub organization or username"
  type        = string
  default     = "ducduyn31"
}

variable "github_repo" {
  description = "GitHub repository name"
  type        = string
  default     = "kin-core-svc"
}
