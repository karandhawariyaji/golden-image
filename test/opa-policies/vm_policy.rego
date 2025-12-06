package vm.policies

# Policy 1: Only allow specific regions
allowed_regions := {"us-central1", "us-west1", "europe-west1"}

deny[msg] {
    input.resource_changes[_].type == "google_compute_instance"
    zone := input.resource_changes[_].change.after.zone
    region := regex.split(zone, "-")[0]  # Extract region from zone
    not allowed_regions[region]
    msg := sprintf("Region '%s' not allowed. Use: %v", [region, allowed_regions])
}

# Policy 2: Only allow specific machine types
allowed_machine_types := {"n1-standard-2", "e2-standard-2", "e2-medium"}

deny[msg] {
    input.resource_changes[_].type == "google_compute_instance"
    machine_type := input.resource_changes[_].change.after.machine_type
    not allowed_machine_types[machine_type]
    msg := sprintf("Machine type '%s' not allowed. Use: %v", [machine_type, allowed_machine_types])
}

# Policy 3: No public IP allowed
deny[msg] {
    input.resource_changes[_].type == "google_compute_instance"
    input.resource_changes[_].change.after.network_interface[_].access_config != null
    msg := "Public IP not allowed for VMs"
}











# package terraform.security

# deny[msg] {
#   r := input.resource_changes[_]
#   r.type == "google_compute_instance"
#   r.change.after.network_interface[_].access_config
#   msg := "Public IP is not allowed on compute instances"
# }
