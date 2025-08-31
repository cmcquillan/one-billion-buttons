resource "random_password" "obb_function_secret" {
  length = 24
}

resource "digitalocean_app" "obb_webapp" {
  spec {
    name   = "obb-webapp"
    region = "nyc"

    env {
      key   = "PG_CONNECTION_STRING"
      value = "host=${digitalocean_database_cluster.primary_db.host} port=${digitalocean_database_cluster.primary_db.port} dbname=${digitalocean_database_db.primary_db.name} user=${digitalocean_database_cluster.primary_db.user} password=${digitalocean_database_cluster.primary_db.password}"
      scope = "RUN_TIME"
      type  = "SECRET"
    }

    env {
      key   = "OOB_FUNCTION_SECRET"
      value = random_password.obb_function_secret.result
      scope = "RUN_AND_BUILD_TIME"
      type  = "SECRET"
    }

    env {
      key   = "RUN_MINIMAP_IN_MAIN"
      value = "true"
      scope = "RUN_AND_BUILD_TIME"
      type  = "GENERAL"
    }

    service {
      name               = "api"
      instance_count     = 1
      instance_size_slug = "basic-xxs"

      github {
        repo   = "cmcquillan/one-billion-buttons"
        branch = "main"
      }

      dockerfile_path = "app/Dockerfile"
      http_port       = "8080"

      health_check {
        http_path             = "/healthcheck/live"
        port                  = 8080
        initial_delay_seconds = 10
        period_seconds        = 10
        timeout_seconds       = 10
      }
    }

    job {
      name               = "makedb"
      kind               = "PRE_DEPLOY"
      instance_count     = 1
      instance_size_slug = "apps-s-1vcpu-0.5gb"
      dockerfile_path    = "makedb/Dockerfile"

      github {
        repo   = "cmcquillan/one-billion-buttons"
        branch = "main"
      }
    }

    # function {
    #   name       = "cronjobs"
    #   source_dir = "functions"

    #   github {
    #     repo   = "cmcquillan/one-billion-buttons"
    #     branch = "main"
    #   }
    # }

    ingress {
      rule {
        component {
          name = "api"
        }
        match {
          path {
            prefix = "/"
          }
        }
      }
    }
  }

  depends_on = [
    digitalocean_database_cluster.primary_db,
    digitalocean_database_db.primary_db,
    digitalocean_database_user.primary_db_user,
    random_password.obb_function_secret,
  ]

  lifecycle {
    ignore_changes = [
      #spec[0].env[0]
    ]
    prevent_destroy = true
  }
}
