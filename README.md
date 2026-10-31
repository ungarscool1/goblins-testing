# goblins-testing
Repo used to test ungarscool1/goblins

## Hello world API

A minimal Go HTTP API built with the standard library only.

### Prerequisites

- Go 1.22 or later

### Run

```sh
go run .
```

The server listens on port `8080` by default. Set the `PORT` environment
variable to use another port:

```sh
PORT=9090 go run .
```

Stop it with `Ctrl+C` (SIGINT) or SIGTERM; in-flight requests are given up to
10 seconds to finish.

### Endpoints

| Method | Path       | Response                         |
| ------ | ---------- | -------------------------------- |
| GET    | `/`        | `{"message":"Hello, World!"}`    |
| GET    | `/healthz` | `{"status":"ok"}`                |

```sh
curl localhost:8080/
# {"message":"Hello, World!"}

curl localhost:8080/healthz
# {"status":"ok"}
```

Unknown paths return `404`, and unsupported methods return `405`.

### Test

```sh
go vet ./...
go test ./...
```
