include "root" {
  path = find_in_parent_folders("root.hcl")
}

include "envcommon" {
  path   = "${dirname(find_in_parent_folders("root.hcl"))}/_envcommon/iam-roles.hcl"
  expose = true
}

inputs = {
  roles = {
    "kin-production-deployment" = {
      description = "Deployment role for infrastructure provisioning and EKS deployments"
      assume_role_arn_patterns = [
        include.envcommon.locals.sso_admin_role_arn
      ]
      managed_policies = [
        # Full access for infrastructure provisioning (excludes IAM user/group management)
        "arn:aws:iam::aws:policy/PowerUserAccess",
      ]
      inline_policies = {
        # Scoped IAM access for infrastructure provisioning
        # Allows: roles, policies, instance profiles, OIDC providers, service-linked roles
        # Excludes: users, groups, access keys, account settings
        "iam-infrastructure" = jsonencode({
          Version = "2012-10-17"
          Statement = [
            {
              Sid    = "ReadIAM"
              Effect = "Allow"
              Action = [
                "iam:Get*",
                "iam:List*",
              ]
              Resource = "*"
            },
            {
              Sid    = "ManageRoles"
              Effect = "Allow"
              Action = [
                "iam:*Role",
                "iam:*RolePolicy",
                "iam:*RoleTags",
              ]
              Resource = "*"
            },
            {
              Sid    = "ManagePolicies"
              Effect = "Allow"
              Action = [
                "iam:*Policy",
                "iam:*PolicyVersion",
                "iam:*PolicyTags",
              ]
              Resource = "*"
            },
            {
              Sid    = "ManageInstanceProfiles"
              Effect = "Allow"
              Action = [
                "iam:*InstanceProfile",
                "iam:*RoleToInstanceProfile",
              ]
              Resource = "*"
            },
            {
              Sid    = "ManageOIDCProviders"
              Effect = "Allow"
              Action = [
                "iam:*OpenIDConnectProvider*",
              ]
              Resource = "*"
            },
            {
              Sid    = "ManageServiceLinkedRoles"
              Effect = "Allow"
              Action = [
                "iam:*ServiceLinkedRole*",
              ]
              Resource = "*"
            },
          ]
        })
      }
    }

    # Add more roles here as needed, e.g.:
    # "kin-production-ci" = {
    #   description = "CI/CD pipeline role"
    #   assume_role_arns = ["arn:aws:iam::SHARED_ACCOUNT_ID:role/github-actions"]
    #   managed_policies = [...]
    # }
  }
}
