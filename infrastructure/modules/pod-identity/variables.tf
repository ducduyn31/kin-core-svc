variable "environment" {
  description = "Environment name"
  type        = string
}

variable "project" {
  description = "Project name"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "cluster_name" {
  description = "Name of the EKS cluster"
  type        = string
}

variable "namespace" {
  description = "Kubernetes namespace for the service account"
  type        = string
  default     = "kin"
}

variable "service_account_name" {
  description = "Name of the Kubernetes service account"
  type        = string
  default     = "kin-core-svc"
}

variable "s3_bucket_name" {
  description = "Name of the S3 bucket for media storage"
  type        = string
}

variable "rds_resource_id" {
  description = "The resource ID of the RDS instance (e.g., db-ABCDEFGHIJKL)"
  type        = string
}

variable "rds_iam_user" {
  description = "The database username for IAM authentication"
  type        = string
  default     = "core_svc"
}
