resource "google_project" "bsam" {
  provider = google-beta

  project_id = var.project
  name       = var.project

  labels = {
    "firebase" = "enabled"
  }
}

resource "google_firebase_project" "bsam" {
  provider = google-beta
  project  = google_project.bsam.project_id
}
