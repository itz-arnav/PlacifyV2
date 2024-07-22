# Detect OS
ifeq ($(OS),Windows_NT)
    detected_OS := Windows
else
    detected_OS := $(shell uname -s)
endif

# OS-specific adjustments
ifeq ($(detected_OS),Linux)
    RM = rm -f
    MKDIR = mkdir -p
    RMDIR = rmdir
endif
ifeq ($(detected_OS),Darwin)
    RM = rm -f
    MKDIR = mkdir -p
    RMDIR = rmdir
endif
ifeq ($(detected_OS),Windows)
    RM = del /Q /F
    MKDIR = mkdir
    RMDIR = rd /s /q
    BINARY_NAME := $(BINARY_NAME).exe
endif

# Binary name
BINARY_NAME = placify

# Build directory
BUILD_DIR = ./backend/build
SRC_DIR = ./backend/server
TEST_DIR = ./backend/src
MAIN_GO = $(SRC_DIR)/main.go

# Go commands
GOBUILD = go build
GORUN = go run
GOCLEAN = go clean
GOTEST = go test
GOMOD = go mod

# Default command
all: clean build run

# Build binary
build:
	@echo Building binary...
	@$(GOBUILD) -o "$(BUILD_DIR)/$(BINARY_NAME)" "$(MAIN_GO)"

# Run project
run:
	@echo Running application...
	@$(GORUN) "$(MAIN_GO)"

# Clean up
clean:
	@echo Cleaning up...
	@$(GOCLEAN)
	@if exist "$(BUILD_DIR)\$(BINARY_NAME)" $(RM) "$(BUILD_DIR)\$(BINARY_NAME)"

# Test application
test:
	@echo Running tests...
	@$(GOTEST) "$(TEST_DIR)/..."

# Tidy up dependencies
tidy:
	@echo Tidying up dependencies...
	@$(GOMOD) tidy

# Help
help:
	@echo Available commands:
	@echo all    - Compiles and runs the application.
	@echo build  - Compiles the application and produces a binary
	@echo run    - Runs the Go application
	@echo clean  - Cleans up the binary
	@echo test   - Runs all the tests in the project
	@echo tidy   - Tidies up the dependencies
	@echo help   - Displays available commands
