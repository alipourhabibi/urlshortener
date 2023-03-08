# URL Shortener
Here is an outline for designing a backend service for a URL shortener using Go and the Hexagonal Architecture, with Domain-Driven Design (DDD) principles 

## How to run
```sh
$ go test ./...
$ go build -o shortener .
$ ./shortener migrate -- --config CONFIGFILE.yaml
$ ./shortener run -- --config CONFIGFILE.yaml
```
docker:
```sh
$ docker-compose up
$ docker exec -it base sh
$ ./shortener migrate
```
and you are ready to go

## Featurs:
- Postgres for saving urls
- Redis for caching and JWT
- Docker and Docker compose
- Hexagon Architecture
- DDD Principles
- Tests
