package gen

import (
	"os"
	"path/filepath"

	"html/template"
)

func (i *Input) createMakefile(app string) error {
	file, err := os.Create(filepath.Join(i.Path, "Makefile"))
	if err != nil {
		return err
	}
	defer file.Close()

	type makefile struct {
		App string
	}

	make := makefile{
		App: app,
	}
	tmpl := template.Must(template.New("makefile").Parse(makefileTemplate))
	if err := tmpl.Execute(file, make); err != nil {
		return err
	}
	return nil
}

const makefileTemplate = `
GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')
APP_NAME?={{ .App }}
APP_DIR?=$$(git rev-parse --show-toplevel)
SRC_PACKAGES=$(shell go list -mod=vendor ./... | grep -v "vendor" | grep -v "mocks")
VERSION?=0_0_1

.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z0-9._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; \
{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

local.fmt: ## Lints all the go code in the application.
	gofmt -w $(GOFMT_FILES)

local.check: local.fmt ## Loads all the dependencies to vendor directory
	go mod vendor
	go mod tidy

local.build: local.check ## Generates the artifact with the help of 'go build'
	go build -o $(APP_NAME)_$(VERSION) -ldflags="-s -w"

local.push: local.build ## Pushes built artifact to the specified location

local.run: local.build ## Generates the artifact and start the service in the current directory
	./${APP_NAME}

dockerise: local.check ## Containerise the appliction
	docker build . --tag ${DOCKER_USER}/${PROJECT_NAME}:${VERSION}

docker.lint: ## Linting Dockerfile for
	docker run --rm -v $(APP_DIR):/app -w /app hadolint/hadolint:latest-alpine hadolint Dockerfile

docker.login: ## Establishes the connection to the docker registry
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWD} ${DOCKER_REPO}

docker.publish.image: docker_login ## Publisies the image to the registered docker registry.
	docker push ${DOCKER_USER}/${PROJECT_NAME}:${VERSION}

coverage.lint: ## Lint's application for errors, it is a linters aggregator (https://github.com/golangci/golangci-lint).
	docker run --rm -v $(APP_DIR):/app -w /app golangci/golangci-lint:v1.31-alpine golangci-lint run --color always

coverage.report: ## Publishes the go-report of the appliction (uses go-reportcard)
	docker run --rm -v $(APP_DIR):/app -w /app basnik/goreportcard-cli:latest goreportcard-cli -v

dev.prerequisite.up: ## Sets up the development environment with all necessary components.
	$(APP_DIR)/scripts/prerequisite.sh

generate.mock: ## generates mocks for the selected source packages.
	@go generate ${SRC_PACKAGES}
`
