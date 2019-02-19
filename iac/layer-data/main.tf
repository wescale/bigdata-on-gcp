provider "google" {
  region  = "${var.region}"
  project = "slavayssiere-sandbox"
}

provider "google-beta" {
  region  = "${var.region}"
  project = "slavayssiere-sandbox"
}

variable "region" {
  default = "europe-west1"
}

terraform {
  backend "gcs" {
    bucket = "tf-slavayssiere-wescale"
    prefix = "terraform/layer-data"
  }
}

data "terraform_remote_state" "layer-base" {
  backend = "gcs"

  config {
    bucket = "tf-slavayssiere-wescale"
    prefix = "terraform/layer-base"
  }
}
