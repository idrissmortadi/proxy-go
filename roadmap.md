# Roadmap for Go Reverse Proxy with Rate Limiting

## **1. Requirements Definition**

### Core Features

- Reverse proxy to one or more backend servers
- IP-based rate limiting (e.g. 10 requests per minute)
- Logs requests (method, path, IP, status)
- Configurable backend via CLI flag or config file

---

## **2. Project Structure**

```
/proxy-go
├── main.go           # Entry point
├── config.json       # Optional config for backend target, rate limits
└── go.mod
```

---

## **3. Basic Reverse Proxy**

### 3.1 Set up `main.go`

- Parse config or flags (`flag` package)
  - `--target` for backend (e.g. `http://localhost:3000`)
  - `--port` for proxy (default: 8080)

### 3.2 Implement reverse proxy

- Use `httputil.NewSingleHostReverseProxy`
- Replace the `Director` if needed (e.g. to modify headers)

### 3.3 Serve requests

- Use `http.ListenAndServe`

---

## **4. Add Logging Middleware**

### Log

- Client IP
- Method and URL
- Status code
- Timestamp

Wrap the proxy handler with a middleware that logs using `log.Printf`.

---

## **5. Implement IP Rate Limiting**

### 5.1 Choose method

- Use `golang.org/x/time/rate` (token bucket)

### 5.2 Maintain a map

```go
map[string]*rate.Limiter // IP → limiter
```

### 5.3 Middleware logic

- Extract IP from `http.Request`
- Lookup or create limiter
- If limiter rejects, respond `429 Too Many Requests`

Add eviction logic or use an LRU cache (`golang/groupcache/lru`) to prevent memory bloat.

---

## **6. Optional: Config File Support**

- JSON or YAML (use `encoding/json` or `gopkg.in/yaml.v2`)
- Load config at startup:

  ```json
  {
    "target": "http://localhost:3000",
    "port": 8080,
    "rate_limit": {
      "requests_per_minute": 60
    }
  }
  ```

---

## **7. Final Testing and Polish**

- Test with real backend (e.g. a local HTTP server)
- Test rate limit (curl loop)
- Test with multiple IPs (use `curl --header` spoofing)
- Print startup info (port, target, limits)
- Handle errors gracefully

---

## **8. Possible Extensions**

- Per-route rate limits
- Burst control (e.g. allow 5 reqs instantly, then limit)
- Dashboard with metrics (expose `/metrics`)
- TLS support
- Auth token limiting (instead of IP)
