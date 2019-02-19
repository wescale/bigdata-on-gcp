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

variable "myip" {
  default = "192.168.0.1"
}

terraform {
  backend "gcs" {
    bucket = "tf-slavayssiere-wescale"
    prefix = "terraform/layer-kubernetes"
  }
}

data "terraform_remote_state" "layer-base" {
  backend = "gcs"

  config {
    bucket = "tf-slavayssiere-wescale"
    prefix = "terraform/layer-base"
  }
}
