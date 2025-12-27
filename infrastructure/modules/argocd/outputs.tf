output "argocd_namespace" {
  description = "ArgoCD namespace"
  value       = kubernetes_namespace.argocd.metadata[0].name
}

output "argocd_admin_secret_arn" {
  description = "ARN of the ArgoCD admin password secret"
  value       = aws_secretsmanager_secret.argocd_admin.arn
}

output "argocd_server_url" {
  description = "ArgoCD server URL"
  value       = "https://argocd.${var.project}.internal"
}
