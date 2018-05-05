all: test build-go build-docker deploy

build-go:
	hack/build-go.sh

build-docker:
	hack/build-docker.sh

deploy:
	hack/deploy-k8s.sh

test:
	go test ./...

push:
	docker push sebassch/lb-agent
	docker push sebassch/lb-controller
	docker push sebassch/k8s-lb-controller

clean:
	hack/clean-k8s,sh

.PHONY: all test clean
