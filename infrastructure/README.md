# Kin Infrastructure

This directory contains the IaC for deploying Kin Core Service on AWS using Terragrunt, EKS, and ArgoCD.

## Multi-Account Architecture

```mermaid
graph TD
    subgraph org["AWS Organizations"]
        mgmt["Management Account<br/>Organizations, SCPs, Billing"]

        subgraph infra_ou["Infrastructure OU"]
            shared["Shared Account<br/>• ECR<br/>• CI/CD Role"]
        end

        subgraph workloads_ou["Workloads OU"]
            subgraph prod_ou["Production OU"]
                prod["Kin Production Account<br/>• EKS, RDS, ElastiCache, S3"]
            end
        end

        mgmt --> infra_ou
        mgmt --> workloads_ou
    end
```

## Directory Structure

```
infrastructure/
├── README.md
├── terragrunt.hcl              # Root Terragrunt config
├── bootstrap/
│   └── org/                    # AWS Organizations (run from management account)
├── _envcommon/                 # Shared module configs
├── modules/                    # Custom OpenTofu modules
└── live/
    ├── shared/                 # Shared account (ECR, GitHub OIDC) - plain OpenTofu
    └── production/             # Production environment - Terragrunt
```

## Prerequisites

1. **AWS Management Account** with Organizations enabled
2. **Tools** (installed via Nix flake):
   - OpenTofu
   - Terragrunt
   - AWS CLI v2
   - kubectl, Helm, ArgoCD CLI

Enter the Nix development shell:

```bash
nix develop
```

## Initial Setup

### Step 0: Manual Prerequisites (ClickOps)

Before running any OpenTofu commands, you must manually create the following resources in the **management account** via the AWS Console:

#### 1. Create S3 Bucket for State

1. Go to **S3** in AWS Console
2. Click **Create bucket**
3. Configure:
   - **Bucket name**: `kin-mgmt-tfstate`
   - **Region**: `ap-southeast-2`
   - **Object Ownership**: ACLs disabled (recommended)
   - **Block Public Access**: Block all (enabled)
   - **Bucket Versioning**: Enable
   - **Default encryption**: SSE-S3 or SSE-KMS
4. Click **Create bucket**

#### 2. Create DynamoDB Table for State Locking

1. Go to **DynamoDB** in AWS Console
2. Click **Create table**
3. Configure:
   - **Table name**: `kin-mgmt-tf-locks`
   - **Partition key**: `LockID` (String)
   - **Table settings**: Default settings (On-demand capacity)
4. Click **Create table**

#### Why ClickOps for State Backend?

The state bucket and lock table must exist before OpenTofu can initialize. This is a chicken-and-egg problem - you can't use OpenTofu to create the resources that OpenTofu needs to store its state. These resources are created once per account and rarely change.

---

### Step 1: Bootstrap AWS Organizations (Management Account)

Run from your **management account**:

```bash
cd infrastructure/bootstrap/org

# Copy and edit the variables
cp terraform.tfvars.example terraform.tfvars

tofu init

# If AWS Organization already exists, import it first:
ORG_ID=$(aws organizations describe-organization --query 'Organization.Id' --output text)
tofu import aws_organizations_organization.org $ORG_ID

# Review changes and apply
tofu plan
tofu apply
```

> **Note:** The import step is required if your AWS account already has an Organization.
> If you get `AlreadyInOrganizationException`, run the import command above.

This creates:

- AWS Organizations with OUs
- Shared Resources account (`kin-shared`)
- Production account (`kin-production`)
- Service Control Policies
- State buckets and lock tables for each member account (`kin-shared-tfstate`, `kin-production-tfstate`)

Note the account IDs from the output:

```bash
tofu output
```

#### Secure New Member Accounts

After creating the accounts, you'll receive emails for each new account:

1. Go to [AWS Console](https://console.aws.amazon.com/) → Root user → enter account email → Next → Forgot password
2. Set a strong root password for each account
3. (Optional) Set up MFA/Passkey for each account root user
4. Store credentials securely

> **Note:** You'll use `OrganizationAccountAccessRole` for day-to-day access, not root. The `DenyRootUserActions` SCP blocks most root actions. MFA on member account roots is optional but recommended for compliance.
>
> **To set up MFA:** The `DenyRootUserActions` SCP must be temporarily disabled. Comment out `workloads_deny_root` attachment in `bootstrap/org/main.tf`, run `tofu apply`, set up MFA, then re-enable and apply again.

### Step 2: Deploy Shared Account Resources

> **Note:** The shared state bucket (`kin-shared-tfstate`) was already created by the org bootstrap.

Assume role into the **shared account**:

```bash
# Option 1: Use AWS CLI profiles
aws configure --profile kin-shared
# Set role_arn to the OrganizationAccountAccessRole

# Option 2: Assume role directly
eval $(aws sts assume-role \
  --role-arn arn:aws:iam::SHARED_ACCOUNT_ID:role/OrganizationAccountAccessRole \
  --role-session-name shared \
  --query 'Credentials.[AccessKeyId,SecretAccessKey,SessionToken]' \
  --output text | awk '{print "export AWS_ACCESS_KEY_ID="$1" AWS_SECRET_ACCESS_KEY="$2" AWS_SESSION_TOKEN="$3}')
```

Deploy shared account resources:

```bash
cd infrastructure/live/shared

# Copy and configure variables
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your environment account IDs

tofu init
tofu apply
```

This creates:

- ECR repositories (shared across environments)
- GitHub Actions OIDC role for ECR push
- Cross-account pull policies for all configured environments

Note the `github_actions_role_arn` and `ecr_repository_urls` from the output.

### Step 3: Configure Account Settings

Update `live/production/account.hcl` with your production account ID:

```hcl
locals {
  account_id   = "YOUR_PRODUCTION_ACCOUNT_ID"
  account_name = "kin-production"
}
```

### Step 4: Deploy Production Infrastructure

From the production account:

```bash
cd infrastructure/live/production

# Deploy all infrastructure
terragrunt run-all apply
```

### Step 5: Add GitHub Secrets

Add to your GitHub repository:

| Secret                | Value                                  | Source               |
| --------------------- | -------------------------------------- | -------------------- |
| `AWS_SHARED_ROLE_ARN` | Shared account GitHub Actions role ARN | `live/shared` output |

## Account Access

### AWS CLI Profiles

Configure your `~/.aws/config`:

```ini
[profile kin-management]
region = ap-southeast-2

[profile kin-shared]
role_arn = arn:aws:iam::SHARED_ACCOUNT_ID:role/OrganizationAccountAccessRole
source_profile = kin-management
region = ap-southeast-2

[profile kin-production]
role_arn = arn:aws:iam::PRODUCTION_ACCOUNT_ID:role/OrganizationAccountAccessRole
source_profile = kin-management
region = ap-southeast-2
```

Then use:

```bash
aws --profile kin-production eks update-kubeconfig --name kin-production
```

## Deployment Order

```
1. Organizations (management account)      → Creates accounts, SCPs, and state buckets
2. Shared account (live/shared)            → ECR, GitHub OIDC
3. Production VPC                          → Network infrastructure
4. Production VPC Endpoints                → Private AWS service connectivity
5. Production RDS, ElastiCache, S3         → Data layer (parallel)
6. Production Secrets                      → Secrets Manager
7. Production EKS                          → Kubernetes cluster
8. Production Pod Identity                 → IAM roles + Pod Identity associations
9. Production EKS-addons                   → ESO, ALB Controller, OTEL
10. Production ArgoCD                      → GitOps (auto-deploys app)
```

## Security

- **Account Isolation**: Workloads separated from shared resources
- **SCPs**: Enforce region restrictions, deny root user, require IMDSv2
- **Cross-Account ECR**: Production pulls images from Shared account
- **Pod Identity**: Pods use EKS Pod Identity for AWS access (replaces IRSA)
- **RDS IAM Auth**: Database authentication via IAM tokens (no passwords)
- **Secrets**: AWS Secrets Manager for admin credentials

### RDS IAM Authentication

The application uses IAM database authentication instead of passwords. This provides:

- No password rotation needed - tokens are generated on-demand (valid 15 min)
- Existing connections stay open even after token expiry
- Audit trail via CloudTrail

**One-time database setup** (after RDS is created):

```sql
-- Connect as postgres using Secrets Manager password
CREATE USER core_svc;
GRANT rds_iam TO core_svc;
GRANT ALL PRIVILEGES ON DATABASE kin TO core_svc;
```

**Application connection**:

```go
import "github.com/aws/aws-sdk-go-v2/feature/rds/auth"

// Generate token for each new connection
token, err := auth.BuildAuthToken(ctx, dbEndpoint, region, "core_svc", cfg.Credentials)
// Use token as password
```
