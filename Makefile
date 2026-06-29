.RECIPEPREFIX := >
APP_NAME=k8s-cluster-inspector
IMAGE_NAME=k8s-cluster-inspector:latest

.PHONY: run build metrics clean docker-build docker-run

run:
>go run ./cmd/inspector

metrics:
>go run ./cmd/inspector --metrics

build:
>mkdir -p bin
>go build -o bin/$(APP_NAME) ./cmd/inspector

clean:
>rm -rf bin report.json

docker-build:
>docker build -t $(IMAGE_NAME) .

docker-run:
>docker run --rm --network host -v $(HOME)/.kube:/root/.kube:ro $(IMAGE_NAME)
