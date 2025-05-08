# Proxy-Go

Proxy-Go is a simple reverse proxy server written in Go, created as a learning project.
It provides basic proxy functionality with rate limiting capabilities.

## Features

- **Reverse Proxy**: Forwards requests to configurable backend servers
- **IP-based Rate Limiting**: Restricts requests based on client IP address
- **Request Logging**: Records details of each request
  including method, path, IP, and status
- **Configurable**: Supports customization through configuration options

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/proxy-go.git
   cd proxy-go
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build the project:

   ```bash
   go build -o proxy-go main.go
   ```

## Usage

### Basic Configuration

You can configure the proxy through a configuration file or command-line flags:

```yaml
# Example config.yaml
target: "http://localhost:3000"
proxy_port: 8080
rate_limit: 60
burst_limit: 10
cache_size: 1000
```

### Running the Proxy

```bash
./proxy-go --config config.yaml
```

The proxy will start on the specified port and forward requests to the target
server while applying rate limits.

## Development Status

This project is primarily a learning exercise for understanding Go's HTTP
capabilities, middleware patterns, and rate limiting implementation.
Contributions and suggestions for improvements are welcome!

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE)
file for details.
