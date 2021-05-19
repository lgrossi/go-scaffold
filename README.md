# Application

[![Version](https://img.shields.io/github/v/release/lgrossi/go-scaffold)](https://github.com/lgrossi/go-scaffold/releases/latest)
[![Go](https://img.shields.io/github/go-mod/go-version/lgrossi/go-scaffold)](https://golang.org/doc/go1.16)
![GitHub repo size](https://img.shields.io/github/repo-size/lgrossi/go-scaffold)

[![Discord Channel](https://img.shields.io/discord/528117503952551936.svg?style=flat-square&logo=discord)](https://discord.gg/3NxYnyV)
[![GitHub pull request](https://img.shields.io/github/issues-pr/lgrossi/go-scaffold)](https://github.com/lgrossi/go-scaffold/pulls)
[![GitHub issues](https://img.shields.io/github/issues/lgrossi/go-scaffold)](https://github.com/lgrossi/go-scaffold/issues)


## Project

Describr your project **HERE**

## Builds
| Platform       | Build        |
| :------------- | :----------: |
| MacOS          | [![MacOS Build](https://github.com/lgrossi/go-scaffold/actions/workflows/ci-build-macos.yml/badge.svg?branch=main)](https://github.com/lgrossi/go-scaffold/actions/workflows/ci-build-macos.yml)   |
| Ubuntu         | [![Ubuntu Build](https://github.com/lgrossi/go-scaffold/actions/workflows/ci-build-ubuntu.yml/badge.svg?branch=main)](https://github.com/lgrossi/go-scaffold/actions/workflows/ci-build-ubuntu.yml) |
| Windows        | [![Windows Build](https://github.com/lgrossi/go-scaffold/actions/workflows/ci-build-windows.yml/badge.svg?branch=main)](https://github.com/lgrossi/go-scaffold/actions/workflows/ci-build-windows.yml) |

[![Workflow](https://github.com/lgrossi/go-scaffold/actions/workflows/ci-multiplat-release.yml/badge.svg)](https://github.com/lgrossi/go-scaffold/actions/workflows/ci-multiplat-release.yml)

### Getting **Started**

To run it, simply download the latest release and define your environment variables.
You can set environment type as `dev` if you want to use a `.env` file (store it in the same folder of the login server).

You can also download our docker image and apply the environment variables to your container.

**Enviroment Variables**

|       NAME          |            HOW TO USE                |
| :------------------ | :----------------------------------  |
|`MYSQL_DBNAME`       | `database default database name`     |
|`MYSQL_HOST`         | `database host`                      |
|`MYSQL_PORT`         | `database port`                      |
|`MYSQL_PASS`         | `database password`                  |
|`MYSQL_USER`         | `database username`                  |
|`ENV_LOG_LEVEL`      | `logrus log level for verbose` [ref](https://pkg.go.dev/github.com/sirupsen/logrus#Level)   |
|`LOGIN_IP`           | `login ip address`                   |
|`LOGIN_HTTP_PORT`    | `login http port`                    |
|`LOGIN_GRPC_PORT`    | `login grpc port`                    |
|`RATE_LIMITER_BURST` | `rate limiter same request burst`    |
|`RATE_LIMITER_RATE`  | `rate limit request per sec per user`|
|`SERVER_IP`          | `game server IP address`             |
|`SERVER_LOCATION`    | `game server location`               |
|`SERVER_NAME`        | `game server name`                   |
|`SERVER_PORT`        | `game server game port`              |

**Tests**  
`go test ./tests -v`

**Build**  
`RUN go build -o TARGET_NAME ./src/`

## Docker
`docker pull lgrossi/go-scaffold:latest`<br><br>
[![Automation](https://img.shields.io/docker/cloud/automated/lgrossi/go-scaffold)](https://hub.docker.com/r/lgrossi/go-scaffold)
[![Image Size](https://img.shields.io/docker/image-size/lgrossi/go-scaffold)](https://hub.docker.com/r/lgrossi/go-scaffold/tags?page=1&ordering=last_updated)
![Pulls](https://img.shields.io/docker/pulls/lgrossi/go-scaffold)
[![Build](https://img.shields.io/docker/cloud/build/lgrossi/go-scaffold)](https://hub.docker.com/r/lgrossi/go-scaffold/builds)
