# REST API

![Build](https://github.com/anurag-rajawat/rest-api/actions/workflows/build.yaml/badge.svg)
![CI tests](https://github.com/anurag-rajawat/rest-api/actions/workflows/test.yaml/badge.svg)
[![codecov](https://codecov.io/gh/anurag-rajawat/rest-api/branch/main/graph/badge.svg)](https://codecov.io/gh/anurag-rajawat/rest-api)
![License](https://img.shields.io/github/license/anurag-rajawat/rest-api?color=brightgreen)
![Go Version](https://img.shields.io/github/go-mod/go-version/anurag-rajawat/rest-api?color=brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/anurag-rajawat/rest-api)](https://goreportcard.com/report/github.com/anurag-rajawat/rest-api)

This REST API was built with Golang and the Gin web framework. It is backed by a PostgreSQL database and uses unit tests
to ensure its quality.

I built this REST API as a learning project to learn Golang and backend development. I am currently learning Golang and
wanted to build something to solidify my understanding of the language.

## Prerequisites

1. [Docker](https://www.docker.com/products/docker-desktop/)
2. [Go](https://go.dev/doc/install)

## Local Development

Start API server (locally)

```shell
$ docker run --rm -d -p 5432:5432 --name=api-db -e POSTGRES_PASSWORD=test postgres:15 
```

```shell
$ make run
```

Alternatively, using docker compose

```shell
$ make docker-run
```

Run tests
```shell
$ make test
```

## API Resources

- GET `/v1`
- GET `/v1/users`
- GET `/v1/users/{id}`
- PUT `v1/users/{id}`
- DELETE `v1/users/{id}`
- POST `/v1/signup`
- POST `/v1/signin`

