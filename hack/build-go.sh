#! /bin/bash

go build -o docker/k8s-controller/k8s cmd/controller/controller.go

go build -o docker/agent/agent cmd/haproxy/agent/agent.go

go build -o docker/haproxycluster/controller cmd/haproxy/controller/main.go