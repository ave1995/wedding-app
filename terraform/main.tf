provider "google" {
    credentials = file("../credentials.json")
    project = var.project_id
    region  = var.region
}

resource "google_cloud_run_service" "backend" {
    name = var.image_name
    location = var.region

    template {
      spec {
        containers {
          image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.repo_name}/${var.image_name}:latest"

          env {
            name  = "BACKEND_DBNAME"
            value = "your-db-name"
          }

          env {
            name  = "BACKEND_DBUSERNAME"
            value = "your-db-user"
          }

          env {
            name  = "BACKEND_DBPASSWORD"
            value = "your-db-password"
          }
        }
      }
    }

    traffic {
        percent         = 100
        latest_revision = true
    }
}