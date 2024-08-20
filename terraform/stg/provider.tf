terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.42.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "5.42.0"
    }
  }

  backend "gcs" {
    bucket = "tf-state-bucket"
  }
}

provider "google" {
  # Configuration options
}

provider "google-beta" {
  # Configuration options
}