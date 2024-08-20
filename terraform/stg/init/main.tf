module "init" {
  source = "../../modules/init"

  project           = var.project
  location          = var.location
  github_repository = var.github_repository
}
