# **************************************************************************** #
#                                                                              #
#                                                         :::      ::::::::    #
#    Makefile                                           :+:      :+:    :+:    #
#                                                     +:+ +:+         +:+      #
#    By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+         #
#                                                 +#+#+#+#+#+   +#+            #
#    Created: 2026/08/21 15:23:30 by dande-je          #+#    #+#              #
#    Updated: 2026/08/28 17:48:19 by dande-je         ###   ########.fr        #
#                                                                              #
# **************************************************************************** #

# ============================================================================ #
#                                   SHELL                                      #
# ============================================================================ #

SHELL       := /usr/bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
MAKEFLAGS   += --warn-undefined-variables --no-builtin-rules

# ============================================================================ #
#                                 SETTINGS                                     #
# ============================================================================ #

.SILENT:
.NOTPARALLEL:

# ============================================================================ #
#                           ANSI COLOR DEFINITIONS                             #
# ============================================================================ #

COLOR_RESET  := \033[0m
COLOR_GREEN  := \033[32m
COLOR_YELLOW := \033[33m
COLOR_RED    := \033[31m
COLOR_BLUE   := \033[34m
 
# ============================================================================ #
#                                PROJECT VARS                                  #
# ============================================================================ #

PROJECT_NAME         = scorpion
BIN_DIR              = bin
COMPOSE_PATH         = ./
COMPOSE_FILE         = $(COMPOSE_PATH)docker-compose.yml

DOCKER               := docker
DOCKER_PRUNE_FILTERS := --filter label=project=$(PROJECT_NAME)

VERSION              := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT               := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE           := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

export VERSION
export COMMIT
export BUILD_DATE
 
# ============================================================================ #
#                              DOCKER COMMANDS                                 #
# ============================================================================ #

COMPOSE := $(if $(shell command -v docker 2>/dev/null), \
             $(if $(shell docker compose version 2>/dev/null), \
               docker compose --project-name $(PROJECT_NAME) -f $(COMPOSE_FILE), \
               docker-compose --project-name $(PROJECT_NAME) -f $(COMPOSE_FILE)), \
             $(error "Neither 'docker compose' nor 'docker-compose' was found on PATH"))
 
# ============================================================================ #
#                                 functions                                    #
# ============================================================================ #

define build_prod
	$(DOCKER) buildx build \
		--target export \
		--build-arg app_name=$(PROJECT_NAME) \
		--build-arg version=$(VERSION) \
		--build-arg commit=$(COMMIT) \
		--build-arg build_date=$(BUILD_DATE) \
		--output type=local,dest=$(BIN_DIR) \
		.
endef

# ============================================================================ #
#                                   TARGETS                                    #
# ============================================================================ #

.DEFAULT_GOAL := help

##@ General

.PHONY: help
help: ## Show list of targets
	awk 'BEGIN {FS = ":.*##"; printf "\n$(COLOR_CYAN)Usage:$(COLOR_RESET)\n  make $(COLOR_GREEN)<target>$(COLOR_RESET)\n"} \
		/^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(COLOR_GREEN)%-15s$(COLOR_RESET) %s\n", $$1, $$2 } \
		/^##@/ { printf "\n$(COLOR_CYAN)%s$(COLOR_RESET)\n", substr($$0, 5) }' $(MAKEFILE_LIST)

.PHONY: all
all: prod ## create the program

##@ Lifecycle

.PHONY: up
up: ## Build (if needed) and start all containers
	$(COMPOSE) up dev
 
build: ## Build images (parallel, no cache, pull latest base)
	$(COMPOSE) build --parallel --no-cache --pull
 
down: ## Stop containers, keep volumes
	printf "$(COLOR_YELLOW)Stopping containers...$(COLOR_RESET)\n"
	$(COMPOSE) down --remove-orphans --timeout 30
	$(DOCKER) container prune -f $(DOCKER_PRUNE_FILTERS)
	printf "$(COLOR_YELLOW)Containers down; volumes intact.$(COLOR_RESET)\n"

.PHONY: shell
shell: ## Open an interactive shell inside the dev container
	$(COMPOSE) run --rm dev sh

.PHONY: stop
stop: ## Stop containers without removing them
	$(COMPOSE) stop

.PHONY: prod
prod: ## Build the production binary and export it to $(BIN_DIR)/
	printf "$(COLOR_BLUE)Building $(PROJECT_NAME) ($(VERSION))...$(COLOR_RESET)\n"
	$(call build_prod)
	printf "$(COLOR_GREEN)✓ Binary exported to $(BIN_DIR)/$(PROJECT_NAME)$(COLOR_RESET)\n"

.PHONY: run-prod
run-prod: ## Build and run the containerized production image
	$(COMPOSE) --profile prod run --rm --entrypoint sh app

.PHONY: restart
restart: down up ## Restart the stack

##@ CI-ready

.PHONY: lint
lint: ## Run static analysis (golangci-lint) in an isolated container
	$(COMPOSE) run --rm lint
 
.PHONY: test
test: ## Run the test suite with race detection and coverage
	$(COMPOSE) run --rm test

##@ Cleanup

.PHONY: clean
clean: ## Remove containers, volumes and dangling images for this project
	printf "$(COLOR_YELLOW)Cleaning Docker resources...$(COLOR_RESET)\n"
	$(COMPOSE) down --rmi local --volumes --remove-orphans --timeout 30
	printf "$(COLOR_GREEN)✓ Clean complete (volumes preserved)$(COLOR_RESET)\n"
 
.PHONY: fclean
fclean: clean ## Full clean including project volumes
	$(DOCKER) network prune -f
	$(DOCKER) system prune -f $(DOCKER_PRUNE_FILTERS)
	rm -rf $(BIN_DIR)
	printf "$(COLOR_RED)✓ Full clean complete$(COLOR_RESET)\n"
 
.PHONY: re
re: fclean up ## Full clean then start fresh

##@ Observability

.PHONY: logs
logs: ## Tail logs from the dev container
	$(COMPOSE) logs -f dev

.PHONY: ps
ps: ## Show running containers
	@$(COMPOSE) ps
