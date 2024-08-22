resource "google_service_account" "github_actions" {
  project      = var.project
  account_id   = "github-actions"
  display_name = "Github Actions"
}

resource "google_project_iam_member" "github_actions" {
  project = var.project
  role    = "roles/owner"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

resource "google_iam_workload_identity_pool" "github_actions" {
  project                   = var.project
  workload_identity_pool_id = "github-actions-oidc"
}

resource "google_iam_workload_identity_pool_provider" "github_actions" {
  project                            = var.project
  workload_identity_pool_provider_id = "github-actions-oidc-provider"
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_actions.workload_identity_pool_id
  attribute_condition                = "\"${var.github_repository}\" == assertion.repository"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.repository" = "assertion.repository"
  }
}

resource "google_service_account_iam_member" "github_actions_iam_workload_identity_user" {
  service_account_id = "projects/${var.project}/serviceAccounts/github-actions@${var.project}.iam.gserviceaccount.com"
  role               = "roles/iam.workloadIdentityUser"
  member             = "principal://iam.googleapis.com/${google_iam_workload_identity_pool.github_actions.name}/subject/repo:${var.github_repository}:ref:refs/heads/main"
}
