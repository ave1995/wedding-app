provider "google" {
    credentials = file("../credentials.json")
    project = var.project_id
    region  = var.region
}

resource "google_artifact_registry_repository" "repo" {
    repository_id = var.repo_name
    format = "DOCKER"
    location = var.region
}

resource "google_cloud_run_service" "backend" {
    name = var.image_name
    location = var.region

    template {
      spec {
        containers {
          image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.repo_name}/${var.image_name}:latest"
        }
      }
    }

    traffic {
        percent         = 100
        latest_revision = true
    }
}