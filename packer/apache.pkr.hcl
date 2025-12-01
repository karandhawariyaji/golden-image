packer {
  required_plugins {
    ansible = {
      source  = "github.com/hashicorp/ansible"
      version = "~> 1"
    }
    googlecompute = {
      source  = "github.com/hashicorp/googlecompute"
      version = "~> 1"
    }
  }
}

variable "project_id" {
}
variable "region" {
  default = "us-central1"
}
variable "zone" {
  default = "us-central1-a"
}

source "googlecompute" "apache" {
  project_id          = var.project_id
  zone                = var.zone
  source_image_family = "ubuntu-2204-lts"
  image_family        = "golden-apache"
  image_name          = "golden-apache-{{timestamp}}"
  ssh_username        = "packer"
}

build {
  name    = "gcp-apache-image"
  sources = ["source.googlecompute.apache"]

  provisioner "shell" {
    inline = [
      "sudo apt update",
      "sudo apt install -y python3"
    ]
  }

  provisioner "ansible" {
    playbook_file = "../ansible/apache.yml"
  }
}
