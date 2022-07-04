# What Is This?
This is my submission for the Jumia recruitment assessment to build a single page application that serves content from
a provided database.

# How To Run

## Pre-Requisites
- go 1.16+
- docker
- docker compose

## Things To Note
- This assessment uses the newer version of docker compose to execute docker related tasks i.e. command is `docker compose` not `docker-compose`
- The backend uses `localhost:9942/phone-numbers` to serve data
- The frontend listens for requests on `localhost:9943`

## Run With Makefile (Recommended)
### - Run Tests
```shell
$ make test
```
This runs all the tests in packages with coverage.

### - Run The Project Without Building
```shell
$ make run
```
This is the default target if no argument is provided to make.

### - Run the project using docker
```shell
$ make start
```
This uses docker-compose to bring up the frontend and backend in different containers respectively in detached mode.

```shell
$ make stop
```
This stops any running docker services that belong to this project.

### - Clean up created docker images pertaining to this assessment
```shell
$ make clean
```

## Run With docker-compose
### - Running the project
```shell
$ docker compose up -d
```
This will be build the necessary images on first run and bring up the required services in detached mode.

### - Stopping the project
```shell
$ docker compose stop
```
This will stop any running docker services that belong to this project

### - Clean Up
```shell
$ docker compose down
$ docker image rm -f jumia_assessment:backend
$ docker image rm -f jumia_assessment:frontend
```

## Running Manually With Go
### - Using Run Helper
```shell
$ ./runhelper
```

### - Manual Startup
```shell
$ cd backend && go build . && ./assessment &
$ cd ../frontend && go build && ./frontend
```

### - Manual Stop And Clean Up
```shell
$ ctrl + c
$ killall assessment
```
This interrupts execution of the frontend service and kills the assessment process
