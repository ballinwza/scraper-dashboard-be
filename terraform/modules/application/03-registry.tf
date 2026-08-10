resource "google_artifact_registry_repository" "repo" {
  project       = var.project_id
  location      = var.region
  repository_id = "${var.service_name}-${var.environment}"
  description   = "Docker repository for ${var.service_name} ${var.environment}"
  format        = "DOCKER"
}