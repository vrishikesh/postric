BROKER_BINARY=broker.app
AUTH_BINARY=auth.app

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: down build_broker build_auth
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down --remove-orphans
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./server
	@echo "Done!"

## build_auth: builds the authentication binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./server
	@echo "Done!"
