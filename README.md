# go-mumble-ping

Microservice that translates the [Mumble Ping](https://wiki.mumble.info/wiki/Protocol) to JSON.

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kaibloecker/go-mumble-ping)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaibloecker/go-mumble-ping)](https://goreportcard.com/report/github.com/kaibloecker/go-mumble-ping)
![License](https://img.shields.io/github/license/kaibloecker/go-mumble-ping)
![GitHub issues](https://img.shields.io/github/issues-raw/kaibloecker/go-mumble-ping)
![GitHub last commit](https://img.shields.io/github/last-commit/kaibloecker/go-mumble-ping)

## Why?

The murmur server offers a way to see how many users are active without connecting/logging in. Since Mumble is used mostly by communities, it would be great if they could show a user counter on their website. Unfortunately this functionality is only exposed through a custom formatted UDP datagram and therefore not usable by ordinary javascript.

**go-mumble-ping** aims to bridge that gap by exposing a webserver that returns JSON on its `/` route (see [usage](#usage)). The mumble server is queried when the webserver is queried. Successful responses are cached for 15 seconds.


## Requirements

Have a local go install.

## Usage
Build an run:
```bash
$ go run .
```
In another shell, query the server:
```bash
$ curl -s http://localhost:8080/
{
    "server_version": "1.3.4",
    "last_update": 1674767226,
    "connected_users": 5,
    "max_users": 100,
    "bandwidth": 72000
}
```
In case of an error, the server returns an HTTP 5xx status code and a JSON object with a `message` key containing the error message.

You can run **go-mumble-ping** on its own, but it's probably best to integrate it into an existing nginx as an upstream server and binding it to a path like `/mumble_status.json`.


## Configuration

All configuration is done via environment variables.

### Environment Variables

| Variable      | Description                                                  |
| ------------- | ------------------------------------------------------------ |
| `MUMBLE_HOST` | The hostname of your mumble server, defaults to `localhost`. |
| `MUMBLE_PORT` | The port of your mumble server, defaults to `64738`.         |
| `PORT`        | The port go-mumble-ping will listen on, defaults to `8080`.  |


## Project Status

The basic functionality is stable. No further development to be expected.


## Contributing

PRs are open.


## License

This project is released under the [MIT License](https://github.com/kaibloecker/go-mumble-ping/blob/main/LICENSE).

Copyright © 2023 [Kai Blöcker](https://github.com/kaibloecker)
