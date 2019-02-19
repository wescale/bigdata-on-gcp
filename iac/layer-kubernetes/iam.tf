// resource "google_service_account" "account-pubsub" {
//   account_id   = "pubsub-rw"
//   display_name = "PubSub read and write"
// }
// resource "google_service_account_key" "pubsubkey" {
//   service_account_id = "${google_service_account.account-pubsub.name}"
// }
// resource "kubernetes_secret" "pubsub-credentials" {
//   metadata {
//     name = "pubsub-credentials"
//   }
//   data {
//     credentials.json = "${base64decode(google_service_account_key.pubsubkey.private_key)}"
//   }
// }
// data "google_iam_policy" "kubernetes-admin" {
//   binding {
//     role    = "roles/pubsub.editor"
//     members = ["serviceAccount:pubsub-rw@slavayssiere-sandbox.iam.gserviceaccount.com"]
//   }
// }

