# goblins-testing
Repo used to test ungarscool1/goblins

## Endpoints

The server listens on `:8080`.

### `GET /health`

Liveness probe. It doesn't check any dependencies. It only confirms that the
process is up and serving HTTP.

- Returns `200 OK` with `Content-Type: application/json` and the body
  `{"status":"ok"}`.
- Any other method returns `405 Method Not Allowed` with an `Allow: GET` header.

```sh
curl http://localhost:8080/health
# {"status":"ok"}
```
