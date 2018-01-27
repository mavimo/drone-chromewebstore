# drone-chromewebstore

[![GoDoc](https://godoc.org/github.com/mavimo/drone-chromewebstore?status.svg)](https://godoc.org/github.com/mavimo/drone-chromewebstore)
[![Build Status](https://ci.mavimo.tech/api/badges/mavimo/drone-chromewebstore/status.svg)](https://ci.mavimo.tech/mavimo/drone-chromewebstore)
[![codecov](https://codecov.io/gh/mavimo/drone-chromewebstore/branch/master/graph/badge.svg)](https://codecov.io/gh/mavimo/drone-chromewebstore)
[![Go Report Card](https://goreportcard.com/badge/github.com/mavimo/drone-chromewebstore)](https://goreportcard.com/report/github.com/mavimo/drone-chromewebstore)
[![Docker Pulls](https://img.shields.io/docker/pulls/mavimo/drone-chromewebstore.svg)](https://hub.docker.com/r/mavimo/drone-chromewebstore/)
[![](https://images.microbadger.com/badges/image/mavimo/drone-chromewebstore.svg)](https://microbadger.com/images/mavimo/drone-chromewebstore "Get your own image badge on microbadger.com")
[![Release](https://github-release-version.herokuapp.com/github/mavimo/drone-chromewebstore/release.svg?style=flat)](https://github.com/mavimo/drone-chromewebstore/releases/latest)
[![Build status](https://ci.appveyor.com/api/projects/status/cuioqombam9yufdy/branch/master?svg=true)](https://ci.appveyor.com/project/mavimo/drone-chromewebstore/branch/master)


Deploying application on Chrome Webstore with drone CI to an existing function

## Build or Download a binary

The pre-compiled binaries can be downloaded from [release page](https://github.com/mavimo/drone-chromewebstore/releases). Support the following OS type.

* Windows amd64/386
* Linux amd64/386
* Darwin amd64/386

With `Go` installed

```
$ go get -u -v github.com/mavimo/drone-chromewebstore
``` 

or build the binary with the following command:

```
$ make build
```

## Docker

Build the docker image with the following commands:

```
$ make docker
```

Please note incorrectly building the image for the correct x64 linux and with
CGO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/drone-chromewebstore' not found or does not exist..
```

## Usage

TBD
