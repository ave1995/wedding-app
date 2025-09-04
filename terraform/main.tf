provider "google" {
    credentials = file("../credentials.json")
    project = var.project_id
    region  = var.region
}

# --- Firestore Database ---
resource "google_firestore_database" "database" {
  project          = var.project_id
  name             = var.database_name
  location_id      = var.region
  type             = "FIRESTORE_NATIVE"
  database_edition = "ENTERPRISE"
}

# --- Cloud Run Service Account and IAM ---
resource "google_service_account" "cloudrun_sa" {
  depends_on = [ google_firestore_database.database ]
  account_id   = "my-cloudrun-sa"
  display_name = "Cloud Run Service Account"
}

resource "google_project_iam_member" "cloudrun_firestore_access" {
  depends_on = [ google_service_account.cloudrun_sa ]
  project = var.project_id
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.cloudrun_sa.email}"
}

# --- Cloud Run Service (BACKEND) ---
resource "google_cloud_run_service" "backend" {
    name = var.image_name
    location = var.region

    depends_on = [ google_firestore_database.database, google_service_account.cloudrun_sa ]

    template {
      spec {

        service_account_name = google_service_account.cloudrun_sa.email
        containers {
          image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.repo_name}/${var.image_name}:latest"

          env {
            name  = "DBURL"
            value = "mongodb://${google_firestore_database.database.uid}.${var.region}.firestore.goog:443/${google_firestore_database.database.id}?loadBalanced=true&tls=true&retryWrites=false&authMechanism=MONGODB-OIDC&authMechanismProperties=ENVIRONMENT:gcp,TOKEN_RESOURCE:FIRESTORE"
          }
          env {
            name  = "DBNAME"
            value = var.database_name
          }

          env {
            name = "WEB_ORIGIN"
            value = "https://bednarovi.today,https://www.bednarovi.today"
          }

          env {
             name = "SECRETKEY"
             value = "12345"
          }

          env {
            name  = "REDEPLOY_TRIGGER"
            value = timestamp() # updates every apply
          }

          env {
            name = "USERICONS_BUCKET"
            value = "wedding-user-icons"
          }

          env {
            name = "DURATION"
            value = "24h"
          }
        }
      }
    }

    traffic {
        percent         = 100
        latest_revision = true
    }
}

resource "google_cloud_run_service_iam_member" "public_access" {
  location    = google_cloud_run_service.backend.location
  project     = var.project_id
  service     = google_cloud_run_service.backend.name
  role        = "roles/run.invoker"
  member      = "allUsers"  # ← veřejně dostupné pro každého
}

# --- Bucket for User Icons ---
resource "google_storage_bucket" "user_icons" {
  name                        = var.user_icon_name
  location                    = var.region
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "public_read" {
  bucket = google_storage_bucket.user_icons.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

output "backend_service_url" {
  value       = google_cloud_run_service.backend.status[0].url
  description = "The URL of the backend service that the React app will call."
}
