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
resource "digitalocean_droplet" "test_droplet" {
    name   = "test-droplet"
    region = "fra1"
    size   = "s-1vcpu-1gb"
    image  = "ubuntu-24-04-x64"

    user_data = file("./cloud-init.yml")

    tags = ["aiops"]
}
