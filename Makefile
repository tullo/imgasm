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

staticcheck-install:
	@GO111MODULE=on go install honnef.co/go/tools/cmd/staticcheck@v0.1.2
	@$$(go env GOPATH)/bin/staticcheck -debug.version
