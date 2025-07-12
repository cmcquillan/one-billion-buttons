terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.59.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.6.0"
    }
  }

  required_version = ">= 1.12"
}

provider "digitalocean" {
  token = var.DigitalOceanToken
}
