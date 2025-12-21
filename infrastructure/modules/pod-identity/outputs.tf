output "role_arn" {
  description = "ARN of the IAM role for kin-core-svc"
  value       = aws_iam_role.kin_core_svc.arn
}

output "role_name" {
  description = "Name of the IAM role for kin-core-svc"
  value       = aws_iam_role.kin_core_svc.name
}

output "pod_identity_association_id" {
  description = "ID of the Pod Identity Association"
  value       = aws_eks_pod_identity_association.kin_core_svc.association_id
}
