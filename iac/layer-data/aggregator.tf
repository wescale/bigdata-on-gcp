resource "google_redis_instance" "aggregator" {
  name           = "aggregator"
  memory_size_gb = 1
  tier           = "BASIC"

  location_id = "europe-west1-b"

  authorized_network = "demo-net"

  reserved_ip_range = "10.254.0.0/29"

  redis_version = "REDIS_3_2"
  display_name  = "aggregator-test"

  labels {
    app    = "aggregator"
    usage = "redis"
  }
}

data "google_dns_managed_zone" "private" {
  name     = "private-dns-zone"
}

resource "google_dns_record_set" "redis" {
  name = "redis.${data.google_dns_managed_zone.private.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${data.google_dns_managed_zone.private.name}"

  rrdatas = ["${google_redis_instance.aggregator.host}"]
}


