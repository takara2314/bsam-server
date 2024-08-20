resource "google_artifact_registry_repository" "api_repository" {
  location      = var.location
  repository_id = "api-service"
  description   = "B-SAM API Service Repository"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "game_repository" {
  location      = var.location
  repository_id = "game-service"
  description   = "B-SAM Game Service Repository"
  format        = "DOCKER"
}


resource "google_artifact_registry_repository" "auth_repository" {
  location      = var.location
  repository_id = "auth-service"
  description   = "B-SAM Auth Service Repository"
  format        = "DOCKER"
}
