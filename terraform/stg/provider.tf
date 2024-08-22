terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.42.0"
    }
  }

  backend "gcs" {
    bucket = "bsam-stg_tf-state-bucket"
    prefix = "ordinary"
  }
}

provider "google" {
  project = var.project
  region  = var.location
}
