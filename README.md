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

## How to get required parameters

In order to use Chrome Webstore API you should create an account that is autorized, the process is a bit tricky, let see step by step what we need to do.

1. Create a new application on Chorme Webstore
    - Go to (webstore developer dashboard](https://chrome.google.com/webstore/developer/dashboard) (need you're registered as developer in order to publish application to all users)
    - Create a new application (button *Add new item*) and fill form with required information, the first time shoul should upload the first version of your application.
    - Save an copy the application ID (you can see it on URL, eg: on URL `https://chrome.google.com/webstore/developer/edit/hcolchlidglfiofiefppcdnkpbgphkgb?authuser=0` the `hcolchlidglfiofiefppcdnkpbgphkgb` part is the Application ID
2. Create a new authorization, you should follow steps reported in [Using the Chrome Web Store Publish API
](https://developer.chrome.com/webstore/using_webstore_api) guide, o summarize:
    - Visit [https://console.developers.google.com/](https://console.developers.google.com/)
    - Create a new project
    - Enable *Chrome Webstore API* for this project
    - Create auhentication credentials
    - Save your *Client ID* and *Client Secret*
    - Generate your *Refresh Token* and save it

## Usage

### Options available

 - flag `--env-file`: `.env` file to load (useful for debugging)
 - env variable `$PLUGIN_APPLICATION` or flag `--application`: the application ID 
 - env variable `$PLUGIN_CLIENT_ID` or flag `--client-id`: Client ID
 - env variable `$PLUGIN_CLIENT_SECRET` or flag `--client-secret`: Client secret
 - env variable `$PLUGIN_REFRESH_TOKEN` or flag `--refresh-token`: Refresh token
 - env variable `$PLUGIN_SOURCE` or flag `--source`: Application source folder 
 - env variable `$PLUGIN_UPLOAD` or flag `--upload`: indicate if we should upload application to webstore (`true` by default)
 - env variable `$PLUGIN_PUBLISH` or flag `--publish`: indicate if we should publish application in webstore (`true` by default)
 - env variable `$PLUGIN_PUBLISH_TARGET` or flag `--publish-target`: Publish target, should be `default` or `trustedTesters` (`default` by default)

### Configure drone

Configure your drone instance to automatically upload / deploy your application. The configuration need some env variables that tipically are set in secrets section.

Your configuration should looks like:

```yaml
pipeline:
  # Previous steps

  upload-extension:
    image: mavimo/drone-chromewebstore
    secrets: [plugin_application, plugin_client_id, plugin_client_secret, plugin_refresh_token]
    source: ./src
    upload: true
    publish: false
    when:
      event: [push]
      branch: [master]

  deploy-extension:
    image: mavimo/drone-chromewebstore
    secrets: [plugin_application, plugin_client_id, plugin_client_secret, plugin_refresh_token]
    upload: false
    publish: true
    publish_target: trustedTesters
    when:
      event: deployment
      environment: staging

  deploy-extension:
    image: mavimo/drone-chromewebstore
    secrets: [plugin_application, plugin_client_id, plugin_client_secret, plugin_refresh_token]
    upload: false
    publish: true
    publish_target: default
    when:
      event: deployment
      environment: production

  # Next steps
```

Where the `upload-extension` step publish the application using the information configured on drone secrets:

 - `plugin_application`: contains the application ID
 - `plugin_client_id`: contains the client ID
 - `plugin_client_secret`: contains the client secret
 - `plugin_refresh_token`: contains the refresh token generated before

The `source` parameter indicate the root folder of your application (should contains the `manifest.json` file)

The `upload` parameter indicate that we are going to zip and upload a new application version. NB: `manifest.json` should contains a version number bigger than already published version.

The `publish` parameter indicate that we are going to publish uploaded application. By default it publish to `default` group, but you should publish also to `trustedTesters`, for example when deploy on staging env.

## Tips

Since is not possible publish the same version on webstore we should increase it each time, we should use the drone build ID. 
Edit `manifest.json` file

```json
{
    // ...
    "version": "1.0.3.BUILD_NUMBER",
    // ...
}
```
and add to `.drone.yml` as step before upload (tune manifest path and `when` conditions):
```yaml
pipeline:
  prepare:
    image: alpine
    commands:
      - sed -i "s/BUILD_NUMBER/$DRONE_BUILD_NUMBER/g" src/manifest.json
    when:
      event: push
```

## Special Thanks
 - [Bo-Yi Wu](https://github.com/appleboy) for create the drone plugin scaffolding I used for this project (see: [`drone-lambda`](https://github.com/appleboy/drone-lambda))
 - [Brad Rydzewski](https://github.com/bradrydzewski) for create (and open source) the [`drone`](https://github.com/drone/drone) project
