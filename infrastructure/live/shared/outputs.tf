output "ecr_repository_urls" {
  description = "ECR repository URLs for container images"
  value       = { for name, repo in aws_ecr_repository.repos : name => repo.repository_url }
}

output "ecr_repository_arns" {
  description = "ECR repository ARNs"
  value       = { for name, repo in aws_ecr_repository.repos : name => repo.arn }
}

output "github_actions_role_arn" {
  description = "IAM role ARN for GitHub Actions to push to ECR"
  value       = aws_iam_role.github_actions.arn
}

output "aws_account_id" {
  description = "Shared account ID"
  value       = data.aws_caller_identity.current.account_id
}
