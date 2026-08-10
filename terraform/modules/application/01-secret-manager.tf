# 1. สร้าง Secret Resource ใน Secret Manager
resource "google_secret_manager_secret" "db_password" {
  project   = var.project_id
  secret_id = "db-password-${var.service_name}-${var.environment}"

  replication {
    auto {}
  }

  depends_on = [google_project_service.enabled_apis]
}

# 2. ใส่ค่า Secret Data (Version)
# หมายเหตุ: ใน Production ควรรับค่าผ่านตัวแปร sensitive variable หรือดึงมาจากแหล่งปลอดภัย
resource "google_secret_manager_secret_version" "db_password_version" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = var.db_password_value # ค่ารหัสผ่านจริงที่รับมาจาก variable
}

resource "google_secret_manager_secret_iam_member" "sa_accessor" {
  project   = var.project_id
  secret_id = google_secret_manager_secret.db_password.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:758337397665-compute@developer.gserviceaccount.com"
}

# 3. กำหนดสิทธิ์ IAM ให้ Service Account ของ Cloud Run สามารถอ่าน Secret นี้ได้
resource "google_secret_manager_secret_iam_member" "cloud_run_secret_access" {
  project   = var.project_id
  secret_id = google_secret_manager_secret.db_password.secret_id
  role      = "roles/secretmanager.secretAccessor"

  member    = "serviceAccount:${google_service_account.cloud_build_sa.email}"
}