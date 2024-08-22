module "api_service" {
  source = "../modules/services/api"

  project  = var.project
  location = var.location
  image    = var.api_service_image
}

module "game_service" {
  source = "../modules/services/game"

  project  = var.project
  location = var.location
  image    = var.game_service_image
}

module "auth_service" {
  source = "../modules/services/auth"

  project  = var.project
  location = var.location
  image    = var.auth_service_image
}

module "firebase" {
  source = "../modules/firebase"

  project  = var.project
  location = var.location
}

module "bigquery" {
  source = "../modules/bigquery"

  location = var.location
}
