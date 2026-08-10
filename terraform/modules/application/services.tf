resource "google_project_service" "enabled_apis" {
  for_each = toset([
    "secretmanager.googleapis.com", # สำหรับ Secret Manager
    "cloudbuild.googleapis.com",    # สำหรับ Cloud Build
    "artifactregistry.googleapis.com", # สำหรับ Artifact Registry
    "run.googleapis.com"            # สำหรับ Cloud Run
  ])

  project                    = var.project_id
  service                    = each.key
  disable_on_destroy         = false
  disable_dependent_services = false
}

resource "google_project_service" "cloudbuild_api" {
  project                    = var.project_id
  service                    = "cloudbuild.googleapis.com"
  disable_on_destroy         = false
  disable_dependent_services = false
}