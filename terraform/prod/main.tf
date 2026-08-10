terraform {
  required_version = ">= 1.15.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 7.0"
    }
  }
  backend "gcs" {
    bucket = "scraper-dashboard-be-tfstate-prod" # ใช้ Bucket ของ Prod เท่านั้น
    prefix = "terraform/state"
  }
}

module "application" {
  source            = "../modules/application"
  project_id    = var.project_id
  region        = var.region
  environment  = "prod"
  service_name = var.service_name
  github_owner          = var.github_owner
  github_repo       = var.github_repo
  branch_name = var.branch_name
  db_password_value = var.db_password_value
}