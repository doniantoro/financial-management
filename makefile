GOLANG_VERSION := go.v1.19
SERVICE_VERSION := 1.0
GIN_PORT := 4001
run:
	# ./gin -i -p ${GIN_PORT} -d cmd/
	go run cmd/main.go

run-docker:
	sudo docker-compose down && sudo docker-compose up -d --remove-orphans

build:
	@echo "Building the binary..."
	sudo docker build -t multi-finance:${SERVICE_VERSION} .