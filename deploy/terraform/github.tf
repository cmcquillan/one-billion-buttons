import {
  to = github_repository.obb_repo
  id = "one-billion-buttons"

}

resource "github_repository" "obb_repo" {
  name        = "one-billion-buttons"
  description = "Approximately one billion buttons"
  visibility  = "public"
}

resource "github_repository_environment" "obb_prod" {
  environment = "Production"
  repository  = github_repository.obb_repo.name
}

resource "github_actions_environment_variable" "obb_prod_do_app" {
  repository    = github_repository.obb_repo.name
  environment   = github_repository_environment.obb_prod.environment
  variable_name = "DO_APP_ID"
  value         = digitalocean_app.obb_webapp.id
}

resource "github_actions_secret" "obb_do_token" {
  repository      = github_repository.obb_repo.name
  secret_name     = "DO_TOKEN"
  plaintext_value = var.DigitalOceanToken
}
