terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.42.0"
    }
  }

  backend "gcs" {
    bucket = "${var.project}_tf-state-bucket"
  }
}

provider "google" {
  # Configuration options
}
