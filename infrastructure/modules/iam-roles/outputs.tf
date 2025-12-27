output "roles" {
  description = "Map of created IAM roles with their ARNs and names"
  value = {
    for name, role in aws_iam_role.this : name => {
      arn  = role.arn
      name = role.name
    }
  }
}

output "role_arns" {
  description = "Map of role names to ARNs"
  value       = { for name, role in aws_iam_role.this : name => role.arn }
}
