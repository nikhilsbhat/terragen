GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')
APP_NAME?=terragen
APP_DIR?=$$(git rev-parse --show-toplevel)
DEV?=${DEVBOX_TRUE}
SRC_PACKAGES=$(shell go list -mod=vendor ./... | grep -v "vendor" | grep -v "mocks")

.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z0-9._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

local.fmt: ## Lints all the go code in the application.
	gofmt -w $(GOFMT_FILES)

local.check: local.fmt ## Loads all the dependencies to vendor directory
	go mod vendor
	go mod tidy

local.build: local.check ## Generates the artifact with the help of 'go build'
	go build -o $(APP_NAME) -ldflags="-s -w"

local.push: local.build ## Pushes built artifact to the specified location

local.run: local.build ## Generates the artifact and start the service in the current directory
	./${APP_NAME}

local.clean: ## Cleans directory from all built binary and other artifacts
	rm -rf $(APP_NAME) compose.db

dockerise: local.check ## Containerise the appliction
	docker build . --tag ${DOCKER_USER}/${PROJECT_NAME}:${VERSION}

docker.lint: ## Linting Dockerfile for
	if [ -z "${DEV}" ]; then hadolint Dockerfile ; else docker run --rm -v $(APP_DIR):/app -w /app hadolint/hadolint:latest-alpine hadolint Dockerfile ; fi

docker.login: ## Establishes the connection to the docker registry
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWD} ${DOCKER_REPO}

docker.publish.image: docker_login ## Publisies the image to the registered docker registry.
	docker push ${DOCKER_USER}/${PROJECT_NAME}:${VERSION}

coverage.lint: ## Lint's application for errors, it is a linters aggregator (https://github.com/golangci/golangci-lint).
	if [ -z "${DEV}" ]; then golangci-lint run --color always ; else docker run --rm -v $(APP_DIR):/app -w /app golangci/golangci-lint:v1.31-alpine golangci-lint run --color always ; fi

coverage.report: ## Publishes the go-report of the appliction (uses go-reportcard)
	if [ -z "${DEV}" ]; then goreportcard-cli -v ; else docker run --rm -v $(APP_DIR):/app -w /app basnik/goreportcard-cli:latest goreportcard-cli -v ; fi

dev.prerequisite.up: ## Sets up the development environment with all necessary components.
	$(APP_DIR)/scripts/prerequisite.sh

dev.prerequisite.purge: ## Teardown the development environment by removing all components.

install.hooks: ## install pre-push hooks for the repository.
	${APP_DIR}/scripts/hook.sh ${APP_DIR}

print.source: ## prints the source packages from which the mocks would be generated.
	echo ${SRC_PACKAGES}

generate.mock: ## generates mocks for the selected source packages.
	@go generate ${SRC_PACKAGES}