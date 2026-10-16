# Backend

Go REST API using the standard library. Run with `go run ./cmd/api`.
The default port is 8080; override it with `PORT`.

## Architecture

The Go app uses the requested MVC-style packages:

| Package | Responsibility |
| --- | --- |
| `api/models` | Calculator data, commands, history entries, defaults, and typed errors. |
| `api/repositories` | History storage: save, list, find, clear, and retention. |
| `api/services` | Arithmetic, input/state rules, validation, and calculator workflows. |
| `api/routes` | Endpoint paths, HTTP methods, and controller registration. |
| `api/controllers` | Request decoding, service calls, and HTTP status mapping. |
| `api/views` | JSON request/response representations and response rendering. |

`cmd/api/main.go` wires the repository factory into the service and the service
into the controllers, then passes the controller to `api/routes/routes.go`.
React is the visual view layer in `../frontend`.

Request flow: **route → controller → service → repository/models → JSON view**.

The in-memory history repository is scoped to each request and seeded from
the submitted history. It returns copies of stored entries and retains the
last 20 results. This preserves the existing stateless REST API and keeps
browser sessions isolated. Services use the repository for saving results,
restoring history entries, and clearing history.

Services and repositories do not depend on HTTP or JSON. No calculations
run in React.

## Calculate

`POST /api/calculate` with `Content-Type: application/json`.

Operands and results are unformatted decimal strings with a dot separator.
Scientific notation is accepted. Display formatting belongs to the frontend.

```json
{"operation":"add","left":"12.5","right":"3"}
```

Response:

```json
{"result":"15.5"}
```

| Operation | Required operands | Optional operand |
| --- | --- | --- |
| `add` | `left`, `right` | — |
| `subtract` | `left`, `right` | — |
| `multiply` | `left`, `right` | — |
| `divide` | `left`, `right` | — |
| `power` | `left`, `right` | — |
| `sqrt` | `left` | — |
| `percent` | `left` | `base` |

Standalone percentages divide `left` by 100. With `base`, the API computes
that percentage of the base. For example:

```json
{"operation":"percent","left":"10","base":"200"}
```

Returns `{"result":"20"}`. The calculator action API applies that percentage
when processing `200 + 10%`, producing `220`.

Numbers use finite float64 values; results round to 15 significant digits.

## Errors

Errors return JSON: `{"error":"Cannot divide by zero."}`.

- `400`: malformed requests, unsupported operations, missing/invalid operands.
- `413`: request body exceeds 4 KiB.
- `415`: content type is not JSON.
- `422`: division by zero, negative square root, non-real or overflowing result.

Unknown fields and additional JSON documents are rejected.

## Health and checks

Go tests live separately in `tests/controllers`, `tests/services`, and
`tests/repositories`, plus routing checks in `tests/routes`. They exercise application packages through their public APIs.
Run only the test folders with `go test ./tests/...`; `go test ./...` also includes them.

`GET /api/health` returns `{"status":"ok"}`.

```sh
go test ./...
go vet ./...
go build -o bin/server ./cmd/api
```

## Calculator UI actions

`POST /api/calculator/action` processes every UI action in Go, including
number entry/editing, operation selection, chaining, contextual percentages,
Power toggling, result history, and history reuse.

The client sends its last returned state and the next button action:

```json
{
  "state": {
    "display": "0",
    "stored": null,
    "operator": null,
    "awaitingOperand": false,
    "justCalculated": false,
    "expression": "",
    "operandLabel": null,
    "error": null,
    "history": [],
    "nextId": 1
  },
  "action": {"type": "digit", "value": "2"}
}
```

The response is `{"state": {...}}` with the complete updated state. Pass that
state back unchanged with the next action. This keeps the service stateless;
the browser caches state for its current session.

Supported actions:

- `digit`: `value` is one digit.
- `operator`: `value` is `+`, `-`, `*`, `/`, or `^`.
- `restore`: `value` is an entry from the returned history.
- `decimal`, `equals`, `clear`, `backspace`, `sign`, `sqrt`, `percent`,
  and `clear-history` take no value.

Accepted actions return HTTP 200. Domain errors, such as division by zero,
are represented in `state.error`. Invalid actions/states return 400;
non-JSON requests return 415; bodies over 32 KiB return 413.

The standalone `POST /api/calculate` endpoint documented above remains
available for direct arithmetic API consumers.
