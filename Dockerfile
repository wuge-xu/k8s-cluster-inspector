FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o k8s-cluster-inspector ./cmd/inspector

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/k8s-cluster-inspector /usr/local/bin/k8s-cluster-inspector

ENTRYPOINT ["k8s-cluster-inspector"]
