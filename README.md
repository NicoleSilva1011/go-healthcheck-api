# Go Health Check API

A simple API that exposes **health**, **readiness**, and **metrics** endpoints, simulating the behavior of a cloud‑native application running in environments such as **Docker** and **Kubernetes**.

This project demonstrates backend fundamentals, containerization, and basic observability concepts.

---

## Features

* Liveness endpoint (`/health`)
* Readiness endpoint with simulated dependencies (`/ready`)
* Metrics endpoint (`/metrics`)
* Request counting with concurrency safety
* Docker multi‑stage build
* Docker container health check

---

## Endpoints

### `GET /health`

Indicates whether the application is alive.

**Response:**

```json
{
  "status": "ok",
  "uptime": "2m30s"
}
```

Used as a **liveness probe** in cloud environments.

---

### `GET /ready`

Indicates whether the application is ready to receive traffic.
Simulates external dependencies such as a database and cache.

**Response (ready):**

```json
{
  "status": "ready",
  "dependencies": {
    "database": "ok",
    "cache": "ok"
  }
}
```

**Response (not ready):**

* HTTP Status: `503 Service Unavailable`

```json
{
  "status": "not ready",
  "dependencies": {
    "database": "down",
    "cache": "ok"
  }
}
```

Used as a **readiness probe** by load balancers or Kubernetes.

---

### `GET /metrics`

Exposes basic runtime metrics.

```json
{
  "requests": 42,
  "uptime_seconds": 150
}
```

> The metrics endpoint is intentionally excluded from request counting to avoid inflating metrics with monitoring traffic.

---

## Running Locally (without Docker)

### Prerequisites

* Go 1.22+

### Run

```bash
go run main.go
```

### Simulate dependency failure

```bash
export DB_DOWN=true
go run main.go
```

### Test endpoints

```bash
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:8080/metrics
```

---

## Running with Docker

### Build the image

```bash
docker build -t go-healthcheck-api .
```

### Run the container

```bash
docker run -p 8080:8080 go-healthcheck-api
```

### Simulate dependency failure

```bash
docker run -p 8080:8080 -e DB_DOWN=true go-healthcheck-api
```

---

## Docker Health Check

The container includes a Docker `HEALTHCHECK` that verifies application liveness via the `/health` endpoint.

```dockerfile
HEALTHCHECK CMD wget -qO- http://localhost:8080/health || exit 1
```

Docker will report the container as `healthy` or `unhealthy` based on the result.

---

## Cloud‑Native Concepts Demonstrated

* Liveness vs readiness probes
* Environment‑based configuration
* Thread‑safe metrics using atomic operations
* Container health checks
* Separation of application and infrastructure concerns
