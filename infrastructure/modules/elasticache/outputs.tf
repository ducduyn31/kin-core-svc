output "replication_group_id" {
  description = "ElastiCache replication group ID"
  value       = module.elasticache.replication_group_id
}

output "replication_group_arn" {
  description = "ElastiCache replication group ARN"
  value       = module.elasticache.replication_group_arn
}

output "primary_endpoint_address" {
  description = "Primary endpoint address for the replication group"
  value       = module.elasticache.replication_group_primary_endpoint_address
}

output "reader_endpoint_address" {
  description = "Reader endpoint address for the replication group"
  value       = module.elasticache.replication_group_reader_endpoint_address
}

output "port" {
  description = "ElastiCache port"
  value       = 6379
}

output "auth_token_secret_arn" {
  description = "ARN of the Secrets Manager secret containing the auth token"
  value       = aws_secretsmanager_secret.auth_token.arn
}

output "auth_token_secret_name" {
  description = "Name of the Secrets Manager secret containing the auth token"
  value       = aws_secretsmanager_secret.auth_token.name
}
