# MiniIRC - v0.0.1
# By Lguedes -> lguedes@student.42.rio
#
# Inspired By tsolding

all: build

build:
	@echo "Building..."
	@go build -o main cmd/main.go

run:
	@go run cmd/main.go

clean:

	@echo "Cleaning..."
	@rm -f main

.PHONY: all build run clean

