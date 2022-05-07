.DEFAULT_GOAL := help

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.netlify:
	netlify init
	netlify link

.PHONY: login
login: ## Netlify login
login:
	netligy login

.PHONY: deploy
deploy: ## Deploy to draft
deploy: .netlify
	netlify deploy

.PHONY: publish
publish: ## Deploy to production
publish: .netlify
	netlify deploy --prod

.PHONY: open
open: ## Open netlify
open: .netlify
	netlify open

.PHONY: build
build: ## build
build:
	mkdir -p functions
	go get ./...
	go install ./...
	go build -o ./functions/main ./...

.PHONY: test
test: ## test
test:
	go vet ./...
	go test ./...