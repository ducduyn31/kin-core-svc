terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.27"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.17"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.7"
    }
  }
}

provider "kubernetes" {
  host                   = var.cluster_endpoint
  cluster_ca_certificate = base64decode(var.cluster_certificate_authority_data)

  exec {
    api_version = "client.authentication.k8s.io/v1"
    command     = "aws"
    args        = ["eks", "get-token", "--cluster-name", var.cluster_name]
  }
}

provider "helm" {
  kubernetes {
    host                   = var.cluster_endpoint
    cluster_ca_certificate = base64decode(var.cluster_certificate_authority_data)

    exec {
      api_version = "client.authentication.k8s.io/v1"
      command     = "aws"
      args        = ["eks", "get-token", "--cluster-name", var.cluster_name]
    }
  }
}

resource "kubernetes_namespace" "argocd" {
  metadata {
    name = "argocd"
  }
}

resource "random_password" "argocd_admin" {
  length  = 32
  special = false
}

resource "terraform_data" "argocd_admin_bcrypt" {
  input = bcrypt(random_password.argocd_admin.result)

  triggers_replace = [
    random_password.argocd_admin.id
  ]

  lifecycle {
    ignore_changes = [input]
  }
}

resource "aws_secretsmanager_secret" "argocd_admin" {
  name        = "${var.project}/${var.environment}/argocd-admin"
  description = "ArgoCD admin password"

  tags = {
    Environment = var.environment
    Project     = var.project
  }
}

resource "aws_secretsmanager_secret_version" "argocd_admin" {
  secret_id     = aws_secretsmanager_secret.argocd_admin.id
  secret_string = random_password.argocd_admin.result
}

resource "helm_release" "argocd" {
  name       = "argocd"
  namespace  = kubernetes_namespace.argocd.metadata[0].name
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "9.1.9"

  values = [<<-EOT
    global:
      domain: argocd.${var.project}.internal

    configs:
      params:
        server.insecure: true  # TLS terminated at ALB

      cm:
        url: https://argocd.${var.project}.internal
        application.resourceTrackingMethod: annotation

    server:
      replicas: 2

      ingress:
        enabled: true
        ingressClassName: alb
        annotations:
          alb.ingress.kubernetes.io/scheme: internal
          alb.ingress.kubernetes.io/target-type: ip
          alb.ingress.kubernetes.io/backend-protocol: HTTP
          alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'
          alb.ingress.kubernetes.io/ssl-redirect: '443'
        hosts:
          - argocd.${var.project}.internal

      resources:
        requests:
          cpu: 100m
          memory: 256Mi
        limits:
          cpu: 500m
          memory: 512Mi

    controller:
      replicas: 1
      resources:
        requests:
          cpu: 250m
          memory: 512Mi
        limits:
          cpu: 1000m
          memory: 1Gi

    repoServer:
      replicas: 2
      resources:
        requests:
          cpu: 100m
          memory: 256Mi
        limits:
          cpu: 500m
          memory: 512Mi

    applicationSet:
      replicas: 1
      resources:
        requests:
          cpu: 50m
          memory: 128Mi
        limits:
          cpu: 200m
          memory: 256Mi

    notifications:
      enabled: false

    dex:
      enabled: false
  EOT
  ]

  set_sensitive {
    name  = "configs.secret.argocdServerAdminPassword"
    value = terraform_data.argocd_admin_bcrypt.output
  }
}

resource "kubernetes_manifest" "kin_project" {
  depends_on = [helm_release.argocd]

  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "AppProject"
    metadata = {
      name      = var.project
      namespace = "argocd"
    }
    spec = {
      description = "Kin Core Service Project"
      sourceRepos = [var.git_repo_url]
      destinations = [
        {
          namespace = "kin"
          server    = "https://kubernetes.default.svc"
        },
        {
          namespace = "argocd"
          server    = "https://kubernetes.default.svc"
        }
      ]
      clusterResourceWhitelist = [
        {
          group = ""
          kind  = "Namespace"
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "app_of_apps" {
  depends_on = [kubernetes_manifest.kin_project]

  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "Application"
    metadata = {
      name      = "app-of-apps"
      namespace = "argocd"
    }
    spec = {
      project = var.project
      source = {
        repoURL        = var.git_repo_url
        targetRevision = var.git_branch
        path           = "argocd/applications"
      }
      destination = {
        server    = "https://kubernetes.default.svc"
        namespace = "argocd"
      }
      syncPolicy = {
        automated = {
          prune    = true
          selfHeal = true
        }
      }
    }
  }
}
