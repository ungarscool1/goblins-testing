# Conformance

Checks every change has to pass.

## Endpoints

- `GET /` returns 200, `text/plain; charset=utf-8`, body `hello world`.
- `GET /health` returns 200, `application/json`, body `{"status":"ok"}`.
- `GET /version` returns 200, `application/json`, body `{"version":"1"}`.

Other methods on these paths (including HEAD) return 405 with `Allow: GET`.
