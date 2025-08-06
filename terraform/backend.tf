terraform {
  backend "gcs" {
    bucket  = "wedding-468009-terraform-state-bucket"
    prefix  = "cloud-run/backend"
  }
}