// Raw topics

resource "google_pubsub_topic" "twitter-raw" {
  name = "twitter-raw"
}

resource "google_pubsub_subscription" "twitter-raw-sub" {
  name  = "twitter-raw-sub"
  topic = "${google_pubsub_topic.twitter-raw.name}"

  ack_deadline_seconds = 20
}

resource "google_pubsub_topic" "mastodon-raw" {
  name = "mastodon-raw"
}

resource "google_pubsub_subscription" "mastodon-raw-sub" {
  name  = "mastodon-raw-sub"
  topic = "${google_pubsub_topic.mastodon-raw.name}"

  ack_deadline_seconds = 20
}

// Normalized topics

resource "google_pubsub_topic" "messages-normalized" {
  name = "messages-normalized"
}

resource "google_pubsub_subscription" "messages-normalized-sub-bigtable" {
  name  = "messages-normalized-sub-bigtable"
  topic = "${google_pubsub_topic.messages-normalized.name}"

  ack_deadline_seconds = 20
}

resource "google_pubsub_subscription" "messages-normalized-sub-datastore" {
  name  = "messages-normalized-sub-datastore"
  topic = "${google_pubsub_topic.messages-normalized.name}"

  ack_deadline_seconds = 20
}

resource "google_pubsub_subscription" "messages-normalized-sub-dataproc" {
  name  = "messages-normalized-sub-dataproc"
  topic = "${google_pubsub_topic.messages-normalized.name}"

  ack_deadline_seconds = 20
}

// aggregator
resource "google_pubsub_topic" "aggregator-queue" {
  name = "aggregator-queue"
}

resource "google_pubsub_subscription" "aggregator-queue-sub" {
  name  = "aggregator-queue-sub"
  topic = "${google_pubsub_topic.aggregator-queue.name}"

  ack_deadline_seconds = 20
}

resource "google_pubsub_subscription" "aggregator-queue-sub-dataset" {
  name  = "aggregator-queue-sub-dataset"
  topic = "${google_pubsub_topic.aggregator-queue.name}"

  ack_deadline_seconds = 20
}

resource "google_pubsub_topic" "messages-public" {
  name = "messages-public"
}


resource "google_pubsub_subscription" "messages-public-sub" {
  name  = "messages-public-sub"
  topic = "${google_pubsub_topic.messages-public.name}"

  ack_deadline_seconds = 20
}
