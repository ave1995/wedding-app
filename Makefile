include .env

NETWORK_NAME=wedding-app_default

up: 
	docker-compose up -d

down:
	docker-compose down

mongo-express:
	docker run -d \
		--name mongo-express \
		--network $(NETWORK_NAME) \
		-p 8081:8081 \
		-e ME_CONFIG_MONGODB_SERVER=mongo \
		-e ME_CONFIG_MONGODB_ADMINUSERNAME=$(MONGO_INITDB_ROOT_USERNAME) \
		-e ME_CONFIG_MONGODB_ADMINPASSWORD=$(MONGO_INITDB_ROOT_PASSWORD) \
		-e ME_CONFIG_MONGODB_AUTH_DATABASE=admin \
		mongo-express

mongo-express-stop:
	docker rm -f mongo-express

swag:
	cd backend && swag init -g main.go -o api/restapi/docs

go:
	@echo "Running Go app from backend/ with .env"
	@bash -c 'set -a; source .env; set +a; cd backend && go run main.go'

GO_BUILD_NAME=my-go-app

go-build:
	cd backend && docker build -t ${GO_BUILD_NAME} .

go-build-run:
	cd backend && docker run --rm ${GO_BUILD_NAME}

## Terraform basic
TF_DIR 	:= ./terraform
TF_VARS :=	terraform.tfvars

## Values from tfvars
GCP_PROJECT_ID := $(shell grep project_id $(TF_DIR)/$(TF_VARS) | cut -d'"' -f2)
REGION     := $(shell grep region $(TF_DIR)/$(TF_VARS) | cut -d'"' -f2)
REPO       := $(shell grep repo_name $(TF_DIR)/$(TF_VARS) | cut -d'"' -f2)
IMAGE      := $(shell grep image_name $(TF_DIR)/$(TF_VARS) | cut -d'"' -f2)
DATABASE      := $(shell grep database_name $(TF_DIR)/$(TF_VARS) | cut -d'"' -f2)

TF_BUCKET	:= $(GCP_PROJECT_ID)-terraform-state-bucket

IMG_URL := $(REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/$(REPO)/$(IMAGE):latest

create-service-account: ## Create service account in gCloud for Terraform deploying.
	@if ! gcloud iam service-accounts list \
		--project=$(GCP_PROJECT_ID) \
		--filter="email=terraform@$(GCP_PROJECT_ID).iam.gserviceaccount.com" \
		--format="value(email)" | grep -q terraform@; then \
	echo "🔧 Creating service account..."; \
	gcloud iam service-accounts create terraform \
		--project=$(GCP_PROJECT_ID) \
		--description="Terraform deployer" \
		--display-name="Terraform deployer"; \
	else \
		echo "✅ Service account already exists."; \
	fi

	@echo "🔐 Assigning roles to service account..."
	@for ROLE in roles/editor roles/storage.admin roles/run.admin roles/iam.serviceAccountUser roles/datastore.owner; do \
		gcloud projects add-iam-policy-binding $(GCP_PROJECT_ID) \
			--member="serviceAccount:terraform@$(GCP_PROJECT_ID).iam.gserviceaccount.com" \
			--role="$$ROLE" \
			--condition=None \
			--quiet; \
	done

create-credentials: ## Create service account credentials file
	@if [ ! -f credentials.json ]; then \
		echo "🔑 Creating service account key..."; \
		gcloud iam service-accounts keys create credentials.json \
			--iam-account=terraform@$(GCP_PROJECT_ID).iam.gserviceaccount.com \
			--project=$(GCP_PROJECT_ID); \
	else \
		echo "✅ credentials.json already exists."; \
	fi

create-backend-bucket: ## Create Cloud Storage bucket for saving Terraform state.
	@if ! gsutil ls -p $(GCP_PROJECT_ID) | grep "gs://$(TF_BUCKET)/" > /dev/null; then \
		echo "📦 Creating Terraform backend bucket..."; \
		gsutil mb -p $(GCP_PROJECT_ID) -l $(REGION) -c standard gs://$(TF_BUCKET); \
		gsutil versioning set on gs://$(TF_BUCKET); \
	else \
		echo "✅ Bucket gs://$(TF_BUCKET) already exists."; \
	fi

grant-storage-access:
	@echo "📦 Granting storage access to Terraform bucket..."
	@gsutil iam ch \
		serviceAccount:terraform@$(GCP_PROJECT_ID).iam.gserviceaccount.com:roles/storage.admin \
		gs://$(TF_BUCKET)

create-registry-repository:
	gcloud artifacts repositories create $(REPO) \
		--repository-format=docker \
		--location=$(REGION) \
		--project=$(GCP_PROJECT_ID)

build:
	docker build -t $(IMG_URL) ./backend

## Push image to Artifact Registry
push:
	gcloud auth configure-docker $(REGION)-docker.pkg.dev
	docker push $(IMG_URL)

## Init Terraform in TF_DIR
terraform-init:
	cd $(TF_DIR) && terraform init

## Apply Terraform with tfvars
terraform-apply:
	cd $(TF_DIR) && terraform apply -var-file="$(TF_VARS)" -auto-approve

terraform-destroy:
	gcloud run services delete $(IMAGE) --region=$(REGION) --project=$(GCP_PROJECT_ID);

## Run everything in order
deploy: create-service-account create-credentials create-backend-bucket terraform-init build push terraform-apply

CREDENTIALS_PATH := /Users/ave/wedding-app/credentials.json

set-google-credentials:
	export GOOGLE_APPLICATION_CREDENTIALS=$(CREDENTIALS_PATH)    


.PHONY: \
	create-service-account \
	create-credentials \
	create-backend-bucket \
	terraform-init \
	terraform-apply \
	build \
	push \
	deploy