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

source "googlecompute" "ubuntu" {
  project_id          = var.project_id
  zone                = var.zone
  source_image_family = "ubuntu-2204-lts"
  image_family        = "golden-ubuntu"
  image_name          = "golden-ubuntu-{{timestamp}}"
  ssh_username        = "packer"
}

build {
  name    = "gcp-golden-image"
  sources = ["source.googlecompute.ubuntu"]

  provisioner "shell" {
    inline = [
      "sudo apt update",
      "sudo apt install -y python3 python3-pip"
    ]
  }

  provisioner "ansible" {
    playbook_file = "../ansible/harden.yml"
  }
}
