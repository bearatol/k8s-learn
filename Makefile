build:
	@go build -o ./app-controller/bin/main ./app-controller/main.go && \
	go build -o ./app-model/bin/main ./app-model/main.go
run:
	@docker compose -f ./docker-compose.yml up -d --build

run-controller: build
	MODEL_PORT="6002" CONTROLLER_PORT="6003" ./app-controller/bin/main

run-model: build
	REDIS_PASS="123" REDIS_PORT="6001" MODEL_PORT="6002" ./app-model/bin/main