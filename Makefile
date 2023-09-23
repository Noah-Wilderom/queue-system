LOGGER_BINARY = "loggerService"
QUEUE_LISTENER_BINARY = "queueListener"
QUEUE_WORKER_BINARY = "queueWorker"

up_build: docker_build
	@echo "Starting docker-compose"
	docker-compose up -d
	@echo "Done starting docker-compose"

docker_build:
	@echo "Building docker-compose"
	docker-compose build
	@echo "Done building docker-compose"