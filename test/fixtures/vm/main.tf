provider "google" {
  project = var.project_id
  zone    = var.zone
}

resource "google_compute_instance" "test_vm" {
  name         = "image-test-vm"
  machine_type = "e2-medium"

  boot_disk {
    initialize_params {
      image = var.image_name   # ✅ Your golden image
    }
  }

  network_interface {
    network = "default"
    access_config {}
  }

  tags = ["http-server"]
}
