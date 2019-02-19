resource "google_container_cluster" "test-cluster" {
  provider = "google-beta"
  name     = "test-cluster"
  region   = "${var.region}"

  private_cluster_config {
    enable_private_endpoint = false
    enable_private_nodes    = true
    master_ipv4_cidr_block  = "192.168.16.0/28"
  }

  master_authorized_networks_config {
    cidr_blocks {
      cidr_block   = "81.56.12.49/32"
      display_name = "chez_moi"
    }

    cidr_blocks {
      cidr_block   = "81.250.133.68/32"
      display_name = "wescale.fr"
    }

    cidr_blocks {
      cidr_block   = "195.137.181.15/32"
      display_name = "client1"
    }

    cidr_blocks {
      cidr_block   = "195.81.225.200/32"
      display_name = "client2"
    }

    cidr_blocks {
      cidr_block   = "195.81.224.200/32"
      display_name = "client3"
    }

    cidr_blocks {
      cidr_block   = "${var.myip}/32"
      display_name = "dyn"
    }

    cidr_blocks {
      cidr_block = "81.31.9.164/32"
      display_name = "other"
    }
    
  }

  min_master_version = "1.11.5-gke.5"
  node_version       = "1.11.5-gke.5"

  network    = "projects/slavayssiere-sandbox/global/networks/demo-net"
  subnetwork = "projects/slavayssiere-sandbox/regions/europe-west1/subnetworks/demo-subnet"

  addons_config {
    kubernetes_dashboard {
      disabled = true
    }
  }

  ip_allocation_policy {
    cluster_secondary_range_name  = "c0-pods"
    services_secondary_range_name = "c0-services"
  }

  lifecycle {
    ignore_changes = ["node_pool"]
  }

  node_pool {
    name = "default-pool"
  }
}

resource "google_container_node_pool" "np-default" {
  provider   = "google-beta"
  name       = "np-default"
  region     = "${var.region}"
  cluster    = "${google_container_cluster.test-cluster.name}"
  node_count = 1

  node_config {
    machine_type = "n1-standard-4"

    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/ndev.clouddns.readwrite",
      "https://www.googleapis.com/auth/cloud-language"
    ]

    labels {
      Name = "test-cluster"
    }

    tags = ["kubernetes", "test-cluster"]
  }
}

output "cluster-endpoint" {
  value = "${google_container_cluster.test-cluster.endpoint}"
}
