resource "google_cloud_run_v2_service" "game_service" {
  name        = "game-service"
  description = "B-SAM Game Service"
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
        cpu_idle          = true
      }

      startup_probe {
        failure_threshold     = 5
        initial_delay_seconds = 10
        timeout_seconds       = 3
        period_seconds        = 3

        http_get {
          path = "/healthz"
        }
      }
    }
  }
}

resource "google_cloud_run_service_iam_binding" "game_service" {
  location = google_cloud_run_v2_service.game_service.location
  service  = google_cloud_run_v2_service.game_service.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}

resource "google_cloud_run_domain_mapping" "game_service" {
  location = google_cloud_run_v2_service.game_service.location
  name     = "${var.environment}.game.${var.domain_name}"

  metadata {
    namespace = var.project
  }


  spec {
    route_name = google_cloud_run_v2_service.game_service.name
  }
}
