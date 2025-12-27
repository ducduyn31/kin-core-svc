variable "environment" {
  description = "Environment name"
  type        = string
}

variable "project" {
  description = "Project name"
  type        = string
}

variable "cluster_name" {
  description = "Name of the EKS cluster"
  type        = string
}

variable "cluster_endpoint" {
  description = "Endpoint of the EKS cluster"
  type        = string
}

variable "cluster_certificate_authority_data" {
  description = "Base64 encoded certificate authority data for the EKS cluster"
  type        = string
}

variable "git_repo_url" {
  description = "Git repository URL for GitOps"
  type        = string
}

variable "git_branch" {
  description = "Git branch to track"
  type        = string
  default     = "main"
}
