# Open Telemetry Status

```
--------------------------------
--- /!\ WORK IN PROGRESS /!\ ---
--------------------------------
```

`otel-status` is a utility that verifies the status of various endpoints and
returns results with [Open Telemetry](https://opentelemetry.io/) traces and metrics.

This tool enables you to monitor the status of your endpoints in real-time,
with comprehensive tracing and metrics data
to give you an understanding of the performance of your system.

<!-- TOC -->
* [Open Telemetry Status](#open-telemetry-status)
  * [Information](#information)
  * [Tools](#tools)
    * [Revive](#revive)
  * [Contributing](#contributing)
<!-- TOC -->

## Information

Our approach is to focus on the core functionality and
keep the feature set streamlined for the typical use case.
Features such as alerting, graphical user interface, etc.
are typically handled in subsequent stages of the Open Telemetry flow.
See https://opentelemetry.io/docs/concepts/what-is-opentelemetry/#what-opentelemetry-is-not.

We try to limit dependencies to the bare minimum:
* [Open Telemetry](https://opentelemetry.io/) for tracing and metrics (obviously...)
* Cron scheduler to run the checks (https://github.com/go-co-op/gocron)
* YAML for configuration

## Tools

### Revive

Linter for Go source code.

This project use the [default](https://github.com/mgechev/revive#default-configuration) configuration.

See https://github.com/mgechev/revive#text-editors for more information on installation in IDEs.

## Contributing

Please:
* Check you code with revive before submitting a PR.
* Use https://www.conventionalcommits.org for commit messages.
