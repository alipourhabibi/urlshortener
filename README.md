# URL Shortener
Here is an outline for designing a backend service for a URL shortener using Go and the Hexagonal Architecture, with Domain-Driven Design (DDD) principles 

## How to run
```sh
$ go build -o shortener .
$ ./shortener run -- --config CONFIGFILE.yaml
```
docker:
```sh
$ docker-compose up
```

## Featurs:
- Postgres for saving urls
- Redis for caching and JWT
- Docker and Docker compose
- Hexagon Architecture
- DDD Principles
