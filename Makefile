DOCKER ?= docker

.PHONY: docker-build
docker-build:
	$(DOCKER) build -t nisshiee/kibela-group-set:local .

.PHONY: docker-run
docker-run:
	$(DOCKER) run --rm --env-file .env nisshiee/kibela-group-set:local
