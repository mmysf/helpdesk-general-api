# Service

## Overview

This service is used for handle all endpoint and data about **Service**. [Golang](https://golang.org/) is the main weapon of this service. This service used to implement [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). We basically implement [this go clean architecture](https://github.com/bxcodec/go-clean-arch) but with a lot of customization. Please read their article first for explanation of this architecture.

## Tech Stack
- Golang v1.22
- MongoDB 6.0+
- [Logrus](https://github.com/sirupsen/logrus) for logger

## How to run

1. Clone it
1. Copy paste `.env.example` and rename it into `.env`
1. Adjust the database in your `.env`
1. Run `make run` or `go run main.go`. It will download all dependencies and running your application

## Build the binary

Just run `make build` it build the binary named `build-app`

## Dockerize

- Run `docker build -t build-app . && docker image prune -f` to build the docker image and remove the unused image after building.
- Run `docker run -p 5050:5050 --env-file .env -d build-app` to run the container as daemon

**note : the command above will read the environment variable from `.env` file. You can also pass the env var as command args. Follow [this link](https://docs.docker.com/engine/reference/commandline/run/#set-environment-variables--e---env---env-file).**