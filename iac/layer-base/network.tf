resource "google_compute_network" "demo-net" {
  name                    = "demo-net"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "demo-private-subnet" {
  name          = "demo-subnet"
  ip_cidr_range = "192.168.0.0/20"
  network       = "${google_compute_network.demo-net.self_link}"
  region        = "${var.region}"

  secondary_ip_range {
    range_name    = "c0-pods"
    ip_cidr_range = "10.0.0.0/16"
  }

  secondary_ip_range {
    range_name    = "c0-services"
    ip_cidr_range = "10.1.0.0/16"
  }

  private_ip_google_access = true
}
