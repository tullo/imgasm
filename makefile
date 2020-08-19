export COMPOSE_DOCKER_CLI_BUILD = 1
export DOCKER_BUILDKIT = 1

up:
	@docker-compose up --build

down:
	@docker-compose down --remove-orphans