# sjq - Stream JSON Separator

[![Go Report Card](https://goreportcard.com/badge/github.com/naoyafurudono/sjq)](https://goreportcard.com/report/github.com/naoyafurudono/sjq)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

**sjq** (Stream JSON Separator) is a lightweight CLI tool written in Go that separates JSON and non-JSON lines from mixed log streams. Perfect for processing application logs that contain both structured JSON logs and plain text output.

✨ **Zero dependencies** - Built with Go's standard library only!

## 🎯 Why sjq?

Modern applications often produce mixed output - structured JSON logs alongside plain text messages. **sjq** helps you:

- 📊 **Extract structured logs** for analysis tools (Elasticsearch, CloudWatch, etc.)
- 📝 **Isolate plain text** for debugging and human reading
- 🚀 **Process streams in real-time** with minimal memory footprint
- 🔧 **Integrate easily** into existing log processing pipelines
- 🎯 **Zero dependencies** - single binary, works everywhere Go runs

## 📦 Installation

### Using Go

```sh
go install github.com/naoyafurudono/sjq@latest
```

### From Source

```sh
git clone https://github.com/naoyafurudono/sjq.git
cd sjq
go build -o sjq
```

## 🚀 Quick Start

### Basic Usage

Extract only JSON lines from mixed logs:
```sh
cat app.log | sjq
```

Extract only non-JSON lines:
```sh
cat app.log | sjq -n
```

### Advanced Usage

Separate JSON and non-JSON into different files:
```sh
cat app.log | sjq --json structured.log --non-json plain.log
```

## 📖 Examples

Given a mixed log file (`app.log`):
```
{ "timestamp": "2023-04-01T12:00:00Z", "level": "INFO", "message": "Server started" }
Starting application...
{ "timestamp": "2023-04-01T12:00:01Z", "level": "ERROR", "message": "Connection failed" }
Retrying connection...
{ "timestamp": "2023-04-01T12:00:02Z", "level": "INFO", "message": "Connected" }
```

### Extract JSON logs for monitoring
```sh
cat app.log | sjq > structured.json
# Output: Only the JSON lines
```

### Extract plain text for debugging
```sh
cat app.log | sjq -n > debug.log
# Output: Only the non-JSON lines
```

### Real-time log processing
```sh
tail -f app.log | sjq --json metrics.json --non-json debug.log
```

## 🔧 Options

| Option | Description |
|--------|-------------|
| `-n` | Output only non-JSON lines to stdout |
| `--json FILE` | Write JSON lines to specified file |
| `--non-json FILE` | Write non-JSON lines to specified file |

## 🏗️ Use Cases

### 1. **Log Analysis Pipeline**
```sh
# Send structured logs to Elasticsearch
tail -f app.log | sjq | jq '.' | elasticsearch-bulk-insert

# Keep human-readable logs separately
tail -f app.log | sjq -n > debug.log
```

### 2. **Debugging Production Issues**
```sh
# Extract error messages from mixed logs
cat production.log | sjq | jq 'select(.level == "ERROR")'

# Get non-JSON error output
cat production.log | sjq -n | grep ERROR
```

### 3. **CI/CD Pipeline Integration**
```sh
# Separate test results from debug output
./run-tests.sh | sjq --json test-results.json --non-json test-debug.log
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Inspired by the need to handle mixed log formats in modern cloud-native applications.