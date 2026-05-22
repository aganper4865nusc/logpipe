# logpipe

Lightweight structured log aggregator that tails multiple sources and forwards to configurable sinks.

---

## Installation

```bash
go install github.com/yourorg/logpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/logpipe.git && cd logpipe && go build ./...
```

---

## Usage

Define your sources and sinks in a YAML config file:

```yaml
sources:
  - type: file
    path: /var/log/app/*.log
  - type: stdin

sinks:
  - type: stdout
  - type: http
    endpoint: https://logs.example.com/ingest
    headers:
      Authorization: Bearer $LOG_TOKEN
```

Then run:

```bash
logpipe --config logpipe.yaml
```

logpipe will tail all configured sources and forward structured JSON log entries to every configured sink in real time.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `logpipe.yaml` | Path to config file |
| `--log-level` | `info` | Internal log verbosity |
| `--dry-run` | `false` | Parse config and exit without running |

---

## License

[MIT](LICENSE) © yourorg