all: build-go build-docker

build-go:
	hack/build-go.sh

build-docker:
	hack/build-docker.sh

clean:
	hack/clean-k8s,sh

.PHONY: all test clean
