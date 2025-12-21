# -----------------------------------------------------------------------------
# ECR Repositories
# -----------------------------------------------------------------------------
# Add new repositories to this list
locals {
  ecr_repositories = {
    "kin-core-svc" = {
      description = "Kin Core Service"
    }
    # Add new repositories here:
    # "kin-other-svc" = {
    #   description = "Kin Other Service"
    # }
  }
}

resource "aws_ecr_repository" "repos" {
  for_each = local.ecr_repositories

  name                 = each.key
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "AES256"
  }

  tags = {
    Name = each.value.description
  }
}

resource "aws_ecr_lifecycle_policy" "repos" {
  for_each = local.ecr_repositories

  repository = aws_ecr_repository.repos[each.key].name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 30 images"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = 30
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# ECR Repository Policy - Allow cross-account pull from environment accounts
resource "aws_ecr_repository_policy" "repos" {
  for_each = local.ecr_repositories

  repository = aws_ecr_repository.repos[each.key].name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowCrossAccountPull"
        Effect = "Allow"
        Principal = {
          AWS = [for env, account_id in var.environment_account_ids : "arn:aws:iam::${account_id}:root"]
        }
        Action = [
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:BatchCheckLayerAvailability"
        ]
      }
    ]
  })
}
