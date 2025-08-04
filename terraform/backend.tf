terraform {
  backend "gcs" {
    bucket  = "wedding-468015-terraform-state-bucket"
    prefix  = "cloud-run/backend"
  }
}