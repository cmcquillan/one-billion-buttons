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
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.7.2"
    }
  }

  required_version = ">= 1.12"
}

provider "digitalocean" {
  token = var.DigitalOceanToken
}


provider "cloudflare" {
  api_token = var.CloudflareToken
}
