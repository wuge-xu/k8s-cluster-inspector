.RECIPEPREFIX := >
APP_NAME=k8s-cluster-inspector

.PHONY: run build metrics clean

run:
>go run ./cmd/inspector

metrics:
>go run ./cmd/inspector --metrics

build:
>mkdir -p bin
>go build -o bin/$(APP_NAME) ./cmd/inspector

clean:
>rm -rf bin report.json
