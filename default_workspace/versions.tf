terraform {
    backend "local" {
        path = "terraform.tfstate"
    }
    required_providers {
        digitalocean = {
            source = "digitalocean/digitalocean"
            version = "2.5.0"
        }
    }
}

provider "digitalocean" {
    token = var.do_token
}
