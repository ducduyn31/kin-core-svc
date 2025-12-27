output "external_secrets_role_arn" {
  description = "ARN of the IAM role for External Secrets Operator"
  value       = aws_iam_role.external_secrets.arn
}

output "alb_controller_role_arn" {
  description = "ARN of the IAM role for AWS Load Balancer Controller"
  value       = aws_iam_role.alb_controller.arn
}

output "otel_role_arn" {
  description = "ARN of the IAM role for OTEL Collector"
  value       = aws_iam_role.otel.arn
}

output "external_secrets_installed" {
  description = "Flag indicating External Secrets Operator is installed"
  value       = true
}
