variable "roles" {
  description = "Map of IAM roles to create"
  type = map(object({
    description              = optional(string, "")
    assume_role_arns         = optional(list(string), [])         # Exact ARNs
    assume_role_arn_patterns = optional(list(string), [])         # ARN patterns with wildcards (uses ArnLike condition)
    assume_role_account_id   = optional(string, "")               # Account ID for pattern matching (required if using patterns)
    managed_policies         = optional(list(string), [])
    inline_policies          = optional(map(string), {})          # Map of policy name to policy JSON
  }))
}

variable "account_id" {
  description = "AWS account ID (used as default for assume_role_account_id)"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
