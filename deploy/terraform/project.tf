resource "digitalocean_project" "one_billion_buttons" {
  name        = "One Billion Buttons"
  description = "The One Billion Buttons Web Application"
  purpose     = "Web Application"
  environment = "Production"
}

resource "digitalocean_project_resources" "databases" {
  project = digitalocean_project.one_billion_buttons.id
  resources = [
    digitalocean_database_cluster.primary_db.urn,
  ]
}

resource "digitalocean_project_resources" "platform_apps" {
  project = digitalocean_project.one_billion_buttons.id
  resources = [
    digitalocean_app.obb_webapp.urn,
  ]
}