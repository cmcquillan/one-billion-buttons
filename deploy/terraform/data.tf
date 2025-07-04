resource "digitalocean_database_cluster" "primary_db" {
  name       = "obb-db-primary"
  engine     = "pg"
  version    = "17"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_db" "primary_db" {
  cluster_id = digitalocean_database_cluster.primary_db.id
  name       = "obb"
}

resource "digitalocean_database_user" "primary_db_user" {
  cluster_id = digitalocean_database_cluster.primary_db.id
  name       = "obb_user"
}

resource "digitalocean_database_firewall" "primary_db_allow" {
  cluster_id = digitalocean_database_cluster.primary_db.id

  rule {
    type  = "app"
    value = digitalocean_app.obb_webapp.id
  }
}
