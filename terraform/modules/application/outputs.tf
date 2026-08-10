output "artifact_repository_url" {
  description = "URL ของ Artifact Registry Repository"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.repo.name}"
}

output "cloud_build_service_account" {
  description = "Service Account Email ของ Cloud Build"
  value       = google_service_account.cloud_build_sa.email
}

output "trigger_id" {
  description = "ID ของ Cloud Build Trigger"
  value       = google_cloudbuild_trigger.filename_trigger.trigger_id
}