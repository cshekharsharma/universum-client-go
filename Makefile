# Makefile for Universum

# This Makefile defines several commands to manage and build the project universum.
# The commands include building, testing, cleaning, and running the project,
# as well as configuring the project and generating test coverage reports.

# Colors
RED         := \033[0;31m
GREEN       := \033[0;32m
YELLOW      := \033[1;33m
BLUE        := \033[0;34m
NC          := \033[0m

BINARYNAME := universum

# Go Commands
GOCMD       := go # Command to run Go
GOBUILD     := $(GOCMD) build  # Command to build Go binaries
GOCLEAN     := $(GOCMD) clean  # Command to clean Go binaries and caches
GOGET       := $(GOCMD) get    # Command to download Go modules
GOFMT       := gofmt           # Command to format Go source code
GOINSTALL   := $(GOCMD) install # Command to install Go packages
STATICCHECK := staticcheck     # Command to run static analysis on Go code

# Target: help
# Description: List all make targets with descriptions
.PHONY: help
help:
	@printf "\n${YELLOW}Makefile Targets:${NC}\n\n"
	@printf "  ${GREEN}configure${NC}      - Configure the project\n"
	@printf "  ${GREEN}test${NC}           - Run unit tests\n"

# Target: configure
# Description: Configure the project by tidying and verifying the modules, formatting the code, and running static analysis.
.PHONY: configure
configure:
	@printf "\n${YELLOW}CONFIGURING THE PACKAGE...${NC}\n\n"
	$(GOINSTALL) honnef.co/go/tools/cmd/staticcheck@latest
	$(GOCMD) mod tidy
	$(GOCMD) mod verify
	$(GOFMT) -s -w .
	$(STATICCHECK) ./...
	@printf "\n"

# Target: test
# Description: Run unit tests for the project.
.PHONY: test
test:
	@printf "\n${YELLOW}RUNNING UNIT TESTS...${NC}\n\n"
	$(GOCMD) test ./...
	@printf "\n"
