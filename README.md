go-api-server-starter
-------

## Inspired by following projects

* Laravel
* https://github.com/irahardianto/service-pattern-go

## Feature

* Config
* Log
* JWT
* Rate limiter
* Database CRUD(Todo)

## Development

* go install

## Dependencies

* Config: github.com/spf13/viper
* Database: github.com/jmoiron/sqlx
* Database: github.com/lib/pq
* Log: github.com/sirupsen/logrus
* Web Framework: github.com/gin-gonic/gin
* JWT: github.com/dgrijalva/jwt-go

## Middleware
* JWT: github.com/appleboy/gin-jwt
* Rate limit: github.com/ulule/limiter

## Deploy

* Nginx load banancing
* Supervisor