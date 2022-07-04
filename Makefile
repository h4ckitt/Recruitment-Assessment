SHELL := /bin/bash

.PHONY: run
run:
	@./runhelper

.PHONY: test
test:
	@cd backend && go test -v ./service
	@cd backend && go test -v ./interface/mux/controller

.PHONY: start
start: docker-compose.yml
	@docker compose up -d

.PHONY: stop
stop: docker-compose.yml
	@docker compose stop

.PHONY: clean
clean:
	@docker compose down
	@docker image rm -f jumia_assessment:frontend
	@docker image rm -f jumia_assessment:backend