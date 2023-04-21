# `switchproxy` Server

A customizable HTTP proxy server written in Go, capable of handling requests to different domains by routing them to different upstream proxy servers or handling them locally.

## Features

- Routes requests to different domains based on a YAML configuration file
- Supports routing to upstream proxy servers with authentication
- Handles HTTP CONNECT requests locally using the `goproxy` module
- Can handle requests locally (e.g., for testing or internal services)

## Configuration
The proxy server uses a YAML configuration file to define routing rules. Each rule consists of an array of domains (with regex support) and a proxy server address. If the proxy address is set to local, the request will be handled locally.

Here's an example config.yaml file:

```yaml
rules:
  - domains:
      - ".*domain1.com"
    proxy: "http://username:password@proxy1:8080"
  - domains:
      - ".*domain2.com"
    proxy: "local"
  - domains:
      - ".*"
    proxy: "http://username:password@proxy3:8080"
```

### Usage
Start the proxy server:
```bash
switchproxy -config path/to/your/config.yaml
```

Test the proxy server with curl:
```bash
curl --proxy http://127.0.0.1:8080 http://example.domain1.com
```

Replace http://example.domain1.com with the URL you want to request through the proxy server. Make sure the domain you use matches one of the rules in your configuration file.