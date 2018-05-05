#! /bin/bash

docker build -f docker/k8s-controller/Dockerfile-k8s-controller -t sebassch/k8s-lb-controller docker/k8s-controller

docker build -f docker/agent/Dockerfile-agent -t sebassch/lb-agent docker/agent

docker build -f docker/haproxycluster/Dockerfile-haproxycluster -t sebassch/lb-controller docker/haproxycluster