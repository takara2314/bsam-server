resource "google_bigquery_dataset" "race_log" {
  dataset_id  = "race_log"
  description = "B-SAM Race Log Dataset"
  location    = var.location
}

resource "google_bigquery_table" "geolocations" {
  dataset_id  = google_bigquery_dataset
  table_id    = "geolocations"
  description = "Geolocations of Race Log"
}
