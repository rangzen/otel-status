# Open Telemetry Status

`otel-status` is a Go library that allows you to create tools
that check various endpoints using [Open Telemetry](https://opentelemetry.io/).

These tools allow you to monitor the status of your endpoints in real-time,
with comprehensive traces and metrics data
to help you understand the performance of your system.

See [cmd/otel-status](cmd/otel-status) for an example of a tool that uses this library.

See [docs](docs) for example of results that you can extract.

<!-- TOC -->
* [Open Telemetry Status](#open-telemetry-status)
  * [Information](#information)
  * [CLI tool](#cli-tool)
    * [Installation](#installation)
    * [Usage](#usage)
  * [Tools](#tools)
    * [Revive](#revive)
  * [Related projects](#related-projects)
  * [Contributing](#contributing)
  * [Versioning](#versioning)
<!-- TOC -->

## Information

Our approach is to focus on the core functionality and
keep the feature set streamlined for the typical use case.
Features such as alerting, graphical user interface, etc.
are typically handled in subsequent stages of the Open Telemetry flow.
See https://opentelemetry.io/docs/concepts/what-is-opentelemetry/#what-opentelemetry-is-not.

We try to limit dependencies to the bare minimum:
* [Open Telemetry](https://opentelemetry.io/) for tracing and metrics (obviously...)
* [Cron scheduler](https://github.com/go-co-op/gocron) to run the checks
* [YAML](https://github.com/go-yaml/yaml) for configuration

## CLI tool

### Installation

```shell
go install github.com/rangzen/otel-status/cmd/otel-status@latest
```

### Usage

```shell
otel-status -c config.yaml
```

See [tests/otel-status-compose/otel-status.yaml](tests/otel-status-compose/otel-status.yaml) for an example of configuration file.

## Tools

### Revive

Linter for Go source code.

This project use the [default](https://github.com/mgechev/revive#default-configuration) configuration.

See https://github.com/mgechev/revive#text-editors for more information on installation in IDEs.

## Related projects

* [Open Telemetry](https://opentelemetry.io/)
* [HTTP Check Receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/httpcheckreceiver)

## Contributing

Please:
* Check you code with revive before submitting a PR.
* Use https://www.conventionalcommits.org for commit messages.

## Versioning

We use [SemVer](http://semver.org/) for versioning.
For the versions available, see the [tags on this repository](https://github.com/rangzen/otel-status/tags). 