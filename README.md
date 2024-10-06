# Discord Bot

[![Go](https://img.shields.io/github/go-mod/go-version/tibia-oce/discord-bot)](https://golang.org/doc/go1.16)
![GitHub repo size](https://img.shields.io/github/repo-size/tibia-oce/discord-bot)
[![GitHub pull request](https://img.shields.io/github/issues-pr/tibia-oce/discord-bot)](https://github.com/tibia-oce/discord-bot/pulls)
[![GitHub issues](https://img.shields.io/github/issues/tibia-oce/discord-bot)](https://github.com/tibia-oce/discord-bot/issues)


## Project

Describe your project **HERE**

### Getting **Started**

Describe your project **HERE**

**Enviroment Variables**

|       NAME          |            HOW TO USE                |
| :------------------ | :----------------------------------  |
|`MYSQL_DBNAME`       | `database default database name`     |
|`MYSQL_HOST`         | `database host`                      |
|`MYSQL_PORT`         | `database port`                      |
|`MYSQL_PASS`         | `database password`                  |
|`MYSQL_USER`         | `database username`                  |
|`ENV_LOG_LEVEL`      | `logrus log level for verbose` [ref](https://pkg.go.dev/github.com/sirupsen/logrus#Level)   |
|`APP_IP`             | `app ip address`                     |
|`APP_HTTP_PORT`      | `app http port`                      |
|`APP_GRPC_PORT`      | `app grpc port`                      |
|`RATE_LIMITER_BURST` | `rate limiter same request burst`    |
|`RATE_LIMITER_RATE`  | `rate limit request per sec per user`|

**Tests**  
`go test ./tests -v`

**Build**  
`go build -o TARGET_NAME ./src/`

## Docker
`docker pull ghcr.io/tibia-oce/discord-bot:latest`<br><br>
