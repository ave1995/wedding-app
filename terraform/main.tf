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
            value = "https://bednarovi.today"
          }

          env {
             name = "SECRETKEY"
             value = "12345"
          }

          env {
            name  = "REDEPLOY_TRIGGER"
            value = timestamp() # updates every apply
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

# --- Google Cloud Storage Bucket for React Static Files ---
resource "google_storage_bucket" "react_app_static_bucket" {
  name          = "${var.project_id}-react-app-static" # Unique bucket name
  location      = var.region 
  uniform_bucket_level_access = true # Recommended for consistent permissions
  website {
    main_page_suffix = "index.html"
    not_found_page   = "index.html" # Essential for Single Page Application (SPA) routing
  }
  # WARNING: force_destroy = true will delete all objects in the bucket when the bucket is destroyed.
  # Use with caution in production environments!
  force_destroy = true
}

# Make the GCS bucket publicly readable (required for static website hosting)
resource "google_storage_bucket_iam_member" "react_app_static_bucket_public_access" {
  bucket = google_storage_bucket.react_app_static_bucket.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

# --- NEW: Build and Deploy React App to GCS ---
# This null_resource will execute local commands to build and deploy your React app.
resource "null_resource" "build_and_deploy_react_app" {
  # Triggers ensure this resource re-runs if the bucket or service URL changes.
  triggers = {
    bucket_name = google_storage_bucket.react_app_static_bucket.name
    # Directly use the URL from your existing Cloud Run service resource
    service_url = google_cloud_run_service.backend.status[0].url
    # Add a timestamp to force re-run on every 'terraform apply' if desired
    build_timestamp = timestamp()
  }

  # Provisioner to build the React application
  provisioner "local-exec" {
    # Ensure npm install is run to get dependencies, then build the app.
    # VITE_API_URL is injected as an environment variable for Vite.
    # IMPORTANT: Adjust this path to the root directory of your React app
    command = "npm install && VITE_API_URL=${self.triggers.service_url} VITE_GCS_BASE_URL=https://storage.googleapis.com/${self.triggers.bucket_name}/ npm run build"
    working_dir = "${path.module}/../web"
    environment = {
      # Ensure npm is in your system's PATH, or provide the full path to npm
      PATH = "/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin" # call which npm
    }
  }

  # Provisioner to deploy the built React app to the GCS bucket
  provisioner "local-exec" {
    # gsutil -m rsync -r will synchronize the 'dist' folder with the GCS bucket.
    # IMPORTANT: Adjust this path to the root directory of your React app's build output
    command = "gsutil -m rsync -r ${path.module}/../web/dist gs://${self.triggers.bucket_name}"
    # This command can be run from the main.tf directory
    working_dir = "${path.module}"
    environment = {
      # Ensure gsutil is in your system's PATH, or provide the full path to gsutil
      PATH = "/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin" # call which gsutil
    }
  }

  # Ensure the bucket and its public access are fully set up before attempting to upload files.
  # Also ensure the backend service is deployed and its URL is available.
  depends_on = [
    google_storage_bucket_iam_member.react_app_static_bucket_public_access,
    google_cloud_run_service.backend # Ensure backend service is ready before getting its URL
  ]
}

output "backend_service_url" {
  value       = google_cloud_run_service.backend.status[0].url
  description = "The URL of the backend service that the React app will call."
}

output "static_website_url" {
  value       = "https://storage.googleapis.com/${google_storage_bucket.react_app_static_bucket.name}/index.html"
  description = "The public URL of the deployed static React application."
}