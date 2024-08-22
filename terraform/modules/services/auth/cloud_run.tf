resource "google_cloud_run_v2_service" "auth_service" {
  name        = "auth-service"
  description = "B-SAM AUth Service"
  location    = var.location
  ingress     = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = var.image

      env {
        name  = "GOOGLE_CLOUD_PROJECT_ID"
        value = var.project
      }
      env {
        name = "JWT_SECRET_KEY"
        value_source {
          secret_key_ref {
            secret  = "jwt-secret-key"
            version = "1"
          }
        }
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
        startup_cpu_boost = true
      }
    }
  }

  # 未認証のアクセスを許可する設定を追加
  labels = {
    "cloud.googleapis.com/allow-unauthenticated" = "true"
  }
}

resource "google_cloud_run_service_iam_member" "auth_service_public" {
  location = google_cloud_run_v2_service.auth_service.location
  project  = google_cloud_run_v2_service.auth_service.project
  service  = google_cloud_run_v2_service.auth_service.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
