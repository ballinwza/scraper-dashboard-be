resource "google_cloudbuild_trigger" "filename_trigger" {
  project     = var.project_id
  name        = "trigger-${var.service_name}-${var.environment}"
  description = "CI/CD Trigger for ${var.service_name} (${var.environment}) on branch ${var.branch_name}"

  service_account = google_service_account.cloud_build_sa.id

  github {
    owner = var.github_owner
    name  = var.github_repo
    push {
      branch = var.branch_name
    }
  }

  filename = "cloudbuild.yaml"

  substitutions = {
    _REGION      = var.region 
    _ENVIRONMENT = var.environment
    _SERVICE     = "${var.service_name}-${var.environment}"
    _REPO_NAME   = google_artifact_registry_repository.repo.name
  }

  depends_on = [
    google_project_service.enabled_apis,
    google_project_iam_member.cloud_build_roles
  ]
}