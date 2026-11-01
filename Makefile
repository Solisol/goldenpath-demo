VERSION ?= 0.1.0
IMAGE   ?= goldenpath-api:$(VERSION)
CLUSTER ?= goldenpath

.PHONY: test cli run image cluster deploy clean

test:
	go test ./...

cli: ## build the goldenpath CLI into ./bin
	go build -ldflags "-X main.version=$(VERSION)" -o bin/goldenpath ./cmd/goldenpath

run: ## run the service locally on :8080 without containers
	go run ./cmd/goldenpath-api

image:
	docker build -t $(IMAGE) .

cluster: ## local k3s cluster in Docker, service reachable on localhost:8080
	k3d cluster create $(CLUSTER) -p "8080:80@loadbalancer" --wait

deploy: image
	k3d image import $(IMAGE) -c $(CLUSTER)
	kubectl apply -f deploy/goldenpath.yaml
	kubectl -n idp rollout restart deployment/goldenpath-api
	kubectl -n idp rollout status deployment/goldenpath-api

clean:
	k3d cluster delete $(CLUSTER)
	rm -rf bin
