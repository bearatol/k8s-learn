build:
	CGO_ENABLED=0 go build -o ./app-controller/bin/main ./app-controller/main.go && \
	CGO_ENABLED=0 go build -o ./app-model/bin/main ./app-model/main.go
run: build
	@docker compose -f ./docker-compose.yml up -d --build

# from docker-compose to helm
kompose-helm:
	kompose convert -c
# from docker-compose to manifest kubernetes .yml
kompose:
	kompose convert -f docker-compose.yml

helm-build:
	helm template --debug learn-app ./kopose-k8s-learn > helm-k8s.yml

run-k8s: build
	kubectl apply -f k8s.yaml

run-controller: build
	MODEL_PORT="6002" CONTROLLER_PORT="6003" ./app-controller/bin/main

run-model: build
	REDIS_PASS="123" REDIS_PORT="6001" MODEL_PORT="6002" ./app-model/bin/main