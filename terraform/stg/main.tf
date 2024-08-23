module "api_service" {
  source = "../modules/services/api"

  project     = var.project
  environment = var.environment
  location    = var.location
  image       = var.api_service_image
  domain_name = var.domain_name
}

module "game_service" {
  source = "../modules/services/game"

  project     = var.project
  environment = var.environment
  location    = var.location
  image       = var.game_service_image
  domain_name = var.domain_name
}

module "auth_service" {
  source = "../modules/services/auth"

  project     = var.project
  environment = var.environment
  location    = var.location
  image       = var.auth_service_image
  domain_name = var.domain_name
}

module "bigquery" {
  source = "../modules/bigquery"

  location = var.location
}
