resource "google_artifact_registry_repository" "repo" {
  project       = var.project_id
  location      = var.region
  repository_id = "${var.service_name}-${var.environment}"
  description   = "Docker repository for ${var.service_name} ${var.environment}"
  format        = "DOCKER"

  cleanup_policies {
    id     = "delete-older-than-14d"
    action = "DELETE"
    condition {
      tag_state  = "ANY"
      older_than = "1209600s" # 14 วัน
    }
  }

  cleanup_policies {
    id     = "keep-recent-3-versions"
    action = "KEEP"
    most_recent_versions {
      keep_count = 3
    }
  }
}