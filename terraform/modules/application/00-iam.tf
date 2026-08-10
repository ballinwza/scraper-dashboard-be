# Service Account สำหรับ Cloud Build
resource "google_service_account" "cloud_build_sa" {
  project      = var.project_id
  account_id   = "sab-${var.service_name}-${var.environment}"
  display_name = "Service Account for Cloud Build (${var.service_name}-${var.environment})"
}

# สิทธิ์ของ Cloud Build ในการจัดการ Artifact Registry, Cloud Run และ Logs
resource "google_project_iam_member" "cloud_build_roles" {
  for_each = toset([
    "roles/artifactregistry.writer",
    "roles/run.developer",
    "roles/iam.serviceAccountUser",
    "roles/logging.logWriter",
    "roles/secretmanager.secretAccessor"
  ])

  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.cloud_build_sa.email}"
}