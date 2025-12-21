locals {
  environment = "production"
  project     = "kin"

  domain_name = "api.kin.coffeewithegg.com"
  cors_origin = "https://kin.coffeewithegg.com"

  tags = {
    Environment = "production"
    Project     = "kin"
    ManagedBy   = "terragrunt"
  }
}
