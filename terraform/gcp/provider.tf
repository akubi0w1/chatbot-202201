# https://cloud.google.com/compute/docs/regions-zones?hl=ja
provider "google" {
  project = var.project_id
  region  = var.region

  credentials = file("./credentials/terraform-credential.json")
}

provider "google-beta" {
  project = var.project_id
  region  = var.region

  credentials = file("./credentials/terraform-credential.json")
}
