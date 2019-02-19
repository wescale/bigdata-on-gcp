output "bastion-ip" {
  value = "${google_compute_instance.bastion-europe-1b.network_interface.0.access_config.0.nat_ip}"
}
