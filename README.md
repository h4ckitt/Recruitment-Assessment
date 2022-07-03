# What Is This?
This is my submission for the Jumia recruitment assessment to build a single page application that serves content from
a provided database.

# How To Run

## Pre-Requisites
- go 1.16+
- docker
- docker-compose

## Things To Note
- This assessment uses the newer version of docker compose to execute docker related tasks
- The backend listens on `localhost:9942/phone-numbers`
- The frontend listens on `localhost:9943`

## Run With Makefile (Recommended)
- Run Tests
```shell
$ make test
```
This runs all the tests in packages with coverage.

- Run The Project Without Building
```shell
$ make run
```
This is the default target if no argument is provided to make.

- Run the project using docker (Recommended)
```shell
$ make docker
```
This approach uses docker-compose to bring up the frontend and backend in different containers respectively.

- Clean up created docker images pertaining to this assessment
```shell
$ make clean
```

## Run With docker-compose
- Running the project
```shell
$ docker compose up
```
This will be build the necessary images on first run and bring up the reequired services

- Clean Up
```shell
$ docker compose down && docker image rm -f jumia_assessment
```

## Running Manually With Go
- Startup
```shell
cd backend && go build . && ./assessment &
cd ../frontend && go build && ./frontend
```

- Clean Up
```shell
ctrl + c
killall assessment
```
This interrupts execution of the frontend service and kills the assessment process
