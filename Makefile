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
	swag init -g ./api/main.go -o ./api/docs