resource "google_compute_firewall" "training_fw_rules" {
  name    = "training-fw-rules"
  network = "demo-net"

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  target_tags = ["bastion"]
}

resource "google_compute_firewall" "bastion_to_cluster_fw_rules" {
  name    = "bastion-to-cluster-training-fw-rules"
  network = "demo-net"

  allow {
    protocol = "tcp"
    ports    = ["80", "8080", "443"]
  }

  source_tags = ["bastion"]
  target_tags = ["test-cluster"]
}

data "google_dns_managed_zone" "public-gcp-wescale" {
  name = "slavayssiere-soa"
}

resource "google_dns_record_set" "bastion" {
  name = "bastion.${data.google_dns_managed_zone.public-gcp-wescale.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${data.google_dns_managed_zone.public-gcp-wescale.name}"

  rrdatas = ["${google_compute_instance.bastion-europe-1b.network_interface.0.access_config.0.nat_ip}"]
}

resource "google_compute_instance" "bastion-europe-1b" {
  name                      = "bastion-europe-1b"
  machine_type              = "n1-standard-1"
  zone                      = "${var.region}-b"
  allow_stopping_for_update = true

  tags = ["bastion", "public"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    # network = "demo-net"
    subnetwork = "demo-subnet"

    access_config {
      // Ephemeral IP
    }
  }

  metadata {
    Name     = "bastion"
    ssh-keys = "admin:${file("~/.ssh/id_rsa.pub")}"
  }

  metadata_startup_script = "${file("${path.cwd}/install-vm.sh")}"

  service_account {
    scopes = ["cloud-platform", "compute-rw", "storage-rw"]
  }
}
