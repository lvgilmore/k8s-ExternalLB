#! /bin/bash

kubectl apply -f docker/k8s-controller/k8s-deployment.yaml

kubectl apply -f docker/haproxycluster/k8s-deployment.yaml