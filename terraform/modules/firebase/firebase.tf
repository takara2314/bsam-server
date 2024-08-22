resource "google_firestore_database" "bsam" {
  provider = google-beta

  project     = var.project
  name        = "bsam"
  location_id = var.location
  type        = "FIRESTORE_NATIVE"
}
