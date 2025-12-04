package terraform.security

# -------------------------------
# ✅ ALLOWED REGIONS
# -------------------------------
allowed_regions := {
  "us-central1",
  "us-east1"
}

deny[msg] {
  r := input.resource_changes[_]
  r.type == "google_compute_instance"

  region := split(r.change.after.zone, "/")[length(split(r.change.after.zone, "/")) - 1]
  not startswith(region, "us-central1")
  not startswith(region, "us-east1")

  msg := sprintf("❌ Region not allowed: %s", [region])
}

# -------------------------------
# ✅ ALLOWED MACHINE TYPES
# -------------------------------
allowed_machine_types := {
  "n1-standard-1",
  "e2-small"
}

deny[msg] {
  r := input.resource_changes[_]
  r.type == "google_compute_instance"

  machine := r.change.after.machine_type
  not allowed_machine_types[machine]

  msg := sprintf("❌ Machine type not allowed: %s", [machine])
}

# -------------------------------
# ❌ PUBLIC IP NOT ALLOWED
# -------------------------------
deny[msg] {
  r := input.resource_changes[_]
  r.type == "google_compute_instance"

  r.change.after.network_interface[_].access_config

  msg := "❌ Public IP is not allowed on compute instances"
}














# package terraform.security

# deny[msg] {
#   r := input.resource_changes[_]
#   r.type == "google_compute_instance"
#   r.change.after.network_interface[_].access_config
#   msg := "Public IP is not allowed on compute instances"
# }
