provider "google" {
  region  = "${var.region}"
  project = "slavayssiere-sandbox"
}

variable "region" {
  default = "europe-west1"
}

terraform {
  backend "gcs" {
    bucket = "tf-slavayssiere-wescale"
    prefix = "terraform/layer-base"
  }
}
