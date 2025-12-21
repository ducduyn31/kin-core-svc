output "cluster_id" {
  description = "ElastiCache cluster ID"
  value       = module.elasticache.cluster_id
}

output "cluster_arn" {
  description = "ElastiCache cluster ARN"
  value       = module.elasticache.cluster_arn
}

output "cluster_cache_nodes" {
  description = "List of cache node objects including address and port"
  value       = module.elasticache.cluster_cache_nodes
}

output "cluster_address" {
  description = "DNS name of the cache cluster"
  value       = module.elasticache.cluster_address
}

output "cluster_configuration_endpoint" {
  description = "Configuration endpoint address"
  value       = module.elasticache.cluster_configuration_endpoint
}

output "auth_token_secret_arn" {
  description = "ARN of the Secrets Manager secret containing the auth token"
  value       = aws_secretsmanager_secret.auth_token.arn
}

output "auth_token_secret_name" {
  description = "Name of the Secrets Manager secret containing the auth token"
  value       = aws_secretsmanager_secret.auth_token.name
}
