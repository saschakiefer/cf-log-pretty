# cf-log-pretty

`cf-log-pretty` is a command-line tool designed to format and colorize log output from SAP BTP Cloud Foundry. It parses the standard Cloud Foundry log format, including structured JSON logs, making them easier to read in a terminal.

## Features

- **Human-readable formatting**: Converts dense CF log lines into a clean, readable format.
- **Colorized output**: Highlights log levels (INFO, WARN, ERROR, etc.) for better visibility.
- **Filtering**: Filter logs by minimum log level.
- **Exclusion**: Exclude specific loggers from the output.
- **Truncation**: Truncate raw log messages to terminal width.

## Requirements

- **Go**: 1.25 or higher (as specified in `go.mod`)
- **Cloud Foundry CLI**: To pipe logs into this tool.

## Installation

To install `cf-log-pretty`, you can clone the repository and build it manually:

```bash
git clone https://github.com/saschakiefer/cf-log-pretty.git
cd cf-log-pretty
go install .
```

Alternatively, if you have Go installed, you can use:

```bash
go install github.com/saschakiefer/cf-log-pretty@latest
```

## Usage

The tool reads from standard input (`stdin`), so you can pipe the output of `cf logs` into it:

```bash
cf logs <app-name> | cf-log-pretty
```

### Options

```text
Flags:
  -e, --exclude-logger strings      exclude logs from given loggers. Supports exact match (e.g. "com.foo.Service") or package wildcard (e.g. "com.foo.core.*" for packages and sub-packages)
  -h, --help                        help for cf-log-pretty
  -l, --level string                minimum log level to include (TRACE, DEBUG, INFO, WARN, ERROR). (default "DEBUG")
  -r, --remove-logger-prefix string  remove given prefix from logger names (e.g. "com.foo.prod.")
  -t, --truncate-raw                truncate raw log messages to terminal width (if message is not in JSON format, e.g. platform logs)
```

### Example

Filter logs to show only `WARN` and `ERROR` levels:

```bash
cf logs my-app | cf-log-pretty --level WARN
```

Exclude specific loggers:

```bash
cf logs my-app | cf-log-pretty --exclude-logger "com.sap.cloud.sdk,org.springframework"
```

Remove logger prefix to shorten logger names:

```bash
cf logs my-app | cf-log-pretty --remove-logger-prefix "com.mycompany.prod."
```

Truncate raw log messages (e.g. for platform logs):

```bash
cf logs my-app | cf-log-pretty --truncate-raw
```

## Project Structure

- `main.go`: Entry point of the application.
- `cmd/`: CLI command definitions using Cobra.
- `internal/parser/`: Logic for parsing Cloud Foundry log lines.
- `internal/formatter/`: Logic for colorizing and formatting the output.
- `internal/filter/`: Logic for filtering logs based on level and logger.

## Development

### Building

```bash
go build -o cf-log-pretty main.go
```

### Running Tests

```bash
go test ./...
```

## Environment Variables

Currently, no specific environment variables are used. All configuration is done via CLI flags.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
Copyright (c) 2026 Sascha Kiefer
