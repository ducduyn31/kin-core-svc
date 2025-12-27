# -----------------------------------------------------------------------------
# IAM Roles
# Creates IAM roles that can be assumed by specified principals
# -----------------------------------------------------------------------------

data "aws_iam_policy_document" "assume_role" {
  for_each = var.roles

  # Statement for exact ARNs (no wildcards)
  dynamic "statement" {
    for_each = length(each.value.assume_role_arns) > 0 ? [1] : []
    content {
      sid     = "AllowAssumeRoleExact"
      effect  = "Allow"
      actions = ["sts:AssumeRole"]

      principals {
        type        = "AWS"
        identifiers = each.value.assume_role_arns
      }
    }
  }

  # Statement for ARN patterns with wildcards (uses ArnLike condition)
  dynamic "statement" {
    for_each = length(each.value.assume_role_arn_patterns) > 0 ? [1] : []
    content {
      sid     = "AllowAssumeRolePattern"
      effect  = "Allow"
      actions = ["sts:AssumeRole"]

      # Use account root instead of "*" to avoid overly permissive warning
      principals {
        type        = "AWS"
        identifiers = ["arn:aws:iam::${coalesce(each.value.assume_role_account_id, var.account_id)}:root"]
      }

      condition {
        test     = "ArnLike"
        variable = "aws:PrincipalArn"
        values   = each.value.assume_role_arn_patterns
      }
    }
  }
}

resource "aws_iam_role" "this" {
  for_each = var.roles

  name               = each.key
  description        = each.value.description
  assume_role_policy = data.aws_iam_policy_document.assume_role[each.key].json

  tags = merge(var.tags, {
    Name = each.key
  })
}

locals {
  role_policy_attachments = flatten([
    for role_name, role_config in var.roles : [
      for policy_arn in role_config.managed_policies : {
        role_name  = role_name
        policy_arn = policy_arn
        key        = "${role_name}-${replace(policy_arn, "/[^a-zA-Z0-9]/", "-")}"
      }
    ]
  ])
}

resource "aws_iam_role_policy_attachment" "this" {
  for_each = { for item in local.role_policy_attachments : item.key => item }

  role       = aws_iam_role.this[each.value.role_name].name
  policy_arn = each.value.policy_arn
}

# Inline policies
locals {
  inline_policies = flatten([
    for role_name, role_config in var.roles : [
      for policy_name, policy_json in role_config.inline_policies : {
        role_name   = role_name
        policy_name = policy_name
        policy_json = policy_json
        key         = "${role_name}-${policy_name}"
      }
    ]
  ])
}

resource "aws_iam_role_policy" "this" {
  for_each = { for item in local.inline_policies : item.key => item }

  name   = each.value.policy_name
  role   = aws_iam_role.this[each.value.role_name].name
  policy = each.value.policy_json
}
