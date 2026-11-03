# goblins-testing
Repo used to test ungarscool1/goblins

## GET /time

Returns the server's current UTC time, for checking clock drift between hosts:

```json
{"utc":"2026-11-03T12:34:56.000000100Z","unix_nano":1793709296000000100}
```

Both fields describe the same instant. `utc` is RFC 3339 with exactly 9
fractional digits; `unix_nano` is nanoseconds since the Unix epoch as a JSON
integer (in JavaScript, which stores numbers as float64, use `utc` if you need
exact nanoseconds). The response is sent with `Cache-Control: no-store`.

To estimate drift, record the local time just before sending (`t0`) and just
after receiving (`t1`), then compare against the midpoint:
`drift ≈ unix_nano - (t0 + t1) / 2`. The uncertainty is about half the round
trip, `(t1 - t0) / 2`.
