DOCKER_REMOTE=docker.rebase.com.au

build/connor:
	mkdir -p build
	go build -o build/connor github.com/daemonl/connor/main

.PHONY: docker
docker-build: build/connor
	docker build -t ${DOCKER_REMOTE}/connor .

.PHONY: docker-push
docker-push: docker-build
	docker push ${DOCKER_REMOTE}/connor


