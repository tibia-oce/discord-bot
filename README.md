# Discord Bot

[![Version](https://img.shields.io/github/v/release/tibia-oce/discord-bot)](https://github.com/tibia-oce/discord-bot/releases/latest)
[![Go](https://img.shields.io/github/go-mod/go-version/tibia-oce/discord-bot)](https://golang.org/doc/go1.16)
![GitHub repo size](https://img.shields.io/github/repo-size/tibia-oce/discord-bot)
[![GitHub pull request](https://img.shields.io/github/issues-pr/tibia-oce/discord-bot)](https://github.com/tibia-oce/discord-bot/pulls)
[![GitHub issues](https://img.shields.io/github/issues/tibia-oce/discord-bot)](https://github.com/tibia-oce/discord-bot/issues)


## Project

Describe your project **HERE**

## Builds
| Platform       | Build        |
| :------------- | :----------: |
| MacOS          | [![MacOS Build](https://github.com/tibia-oce/discord-bot/actions/workflows/ci-build-macos.yml/badge.svg?branch=main)](https://github.com/tibia-oce/discord-bot/actions/workflows/ci-build-macos.yml)   |
| Ubuntu         | [![Ubuntu Build](https://github.com/tibia-oce/discord-bot/actions/workflows/ci-build-ubuntu.yml/badge.svg?branch=main)](https://github.com/tibia-oce/discord-bot/actions/workflows/ci-build-ubuntu.yml) |
| Windows        | [![Windows Build](https://github.com/tibia-oce/discord-bot/actions/workflows/ci-build-windows.yml/badge.svg?branch=main)](https://github.com/tibia-oce/discord-bot/actions/workflows/ci-build-windows.yml) |

[![Workflow](https://github.com/tibia-oce/discord-bot/actions/workflows/ci-multiplat-release.yml/badge.svg)](https://github.com/tibia-oce/discord-bot/actions/workflows/ci-multiplat-release.yml)

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
`docker pull tibia-oce/discord-bot:latest`<br><br>
[![Automation](https://img.shields.io/docker/cloud/automated/tibia-oce/discord-bot)](https://hub.docker.com/r/tibia-oce/discord-bot)
[![Image Size](https://img.shields.io/docker/image-size/tibia-oce/discord-bot)](https://hub.docker.com/r/tibia-oce/discord-bot/tags?page=1&ordering=last_updated)
![Pulls](https://img.shields.io/docker/pulls/tibia-oce/discord-bot)
[![Build](https://img.shields.io/docker/cloud/build/tibia-oce/discord-bot)](https://hub.docker.com/r/tibia-oce/discord-bot/builds)
