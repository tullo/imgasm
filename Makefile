SHELL = /bin/bash -o pipefail
export COMPOSE_DOCKER_CLI_BUILD = 1
export DOCKER_BUILDKIT = 1

build:
	@docker-compose build

up:
	@docker-compose up --build

down:
	@docker-compose down --remove-orphans

#   -d flag ...download the source code needed to build ...
#   -t flag ...consider modules needed to build tests ...
#   -u flag ...use newer minor or patch releases when available 
deps-upgrade:
	@go get -u -d -v ./...
	@go mod tidy
	@go mod vendor

check:
	$(shell go env GOPATH)/bin/staticcheck -go 1.15 \
		-tests ./backblaze/... ./db/... ./file/... ./models/... ./ui/templates/...

.PHONY: clone
clone:
	@git clone git@github.com:dominikh/go-tools.git /tmp/go-tools \
		&& cd /tmp/go-tools \
		&& git checkout "2020.1.5" \

.PHONY: install
install:
	@cd /tmp/go-tools && go install -v ./cmd/staticcheck
	$(shell go env GOPATH)/bin/staticcheck -debug.version
