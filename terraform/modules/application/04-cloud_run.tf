resource "google_project_service" "cloud_run_api" {
  project = var.project_id
  service            = "artifactregistry.googleapis.com"
  disable_on_destroy = false
}

# 2. สร้าง Cloud Run Service (ตัวแอป)
resource "google_cloud_run_v2_service" "app" {
  project  = var.project_id
  name     = "${var.service_name}-${var.environment}"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
    }
  }

  depends_on = [
    google_project_service.cloud_run_api
  ]
}

# 3. Allow unauthenticated access (อ้างอิงจากตัว Service ด้านบน)
resource "google_cloud_run_v2_service_iam_member" "public_access" {
  project  = var.project_id
  location = google_cloud_run_v2_service.app.location
  name     = google_cloud_run_v2_service.app.name
  role     = "roles/run.invoker"
  member   = "allUsers"

  depends_on = [
    google_cloud_run_v2_service.app
  ]
}