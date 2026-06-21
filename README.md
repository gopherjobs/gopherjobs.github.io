# Gopher Jobs

Golang Jobs — a job board for Go developers, published at [gopherjobs.github.io](https://gopherjobs.github.io/).

## About

Gopher Jobs is a static site built with Go that aggregates and lists job openings for Go (Golang) developers. The site is hosted via GitHub Pages directly from this repository.

## Tech Stack

- **Go** — core build/aggregation logic (`cmd/`, `internal/`)
- **HTML** — site templates/output
- **Shell** — supporting scripts (`scripts/`)
- **GitHub Actions** — CI/CD (`.github/workflows/`)

## Project Structure

```
.
├── cmd/              # Main application entry point
├── internal/         # Internal Go packages
├── scripts/          # Shell scripts (e.g. job fetching/processing)
├── .github/workflows # CI/CD pipelines
├── Makefile          # Build and run commands
├── go.mod / go.sum    # Go module dependencies
```

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) (version as specified in `go.mod`)
- `make`

### Build

```bash
make build
```

This compiles the binary to `./bin/out`.

### Run

```bash
make run
```

Builds and runs the binary.

### Run job script

```bash
make run-job
```

Runs `scripts/job.sh`.

## Contributing

Issues and pull requests are welcome. If you'd like to add a job listing, fix a bug, or improve the site, feel free to open a PR.

## License

This project is licensed under the [GPL-3.0 License](LICENSE).
