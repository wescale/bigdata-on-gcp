output "instance_group" {
  value = "${element(google_container_cluster.test-cluster.instance_group_urls, 0)}"
}

output "nodepool-groups" {
  value = "${google_container_node_pool.np-default.instance_group_urls}"
}
