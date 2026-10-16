# Disclaimer 
Well, project didn't really needed a Auth system neither a database, i just keept the history in the memory. 
The arquiteture is a Layeared Arquiteture, usually is mny decision to build scalable systemsm
Decision was based in scalability, obviusly this project didnt need it.
Also a few changes into the GO workflow, i like to use a difent test folder to write the testes, also i preffer to use a dedicated router file.

Use of AI
I included the prompts I used to create this project. I was also testing the new GPT SOL 6.1 model.
This project was a good candidate for a single prompt since it was pretty simple, but I preferred to build it step by step so you guys could get a better understanding of how I think and design systems.
I always like to say that I only use AI to write code I could have written myself, and I always review the code it generates. Now, that's not really a strict rule, since some things need less review than others, but it is what it is.
At the end of the day, I always want to understand the code I'm pushing to production.

This readme were all writen by AI, a good case schenario to generate some API exemples and all.
Time, i took probable around 2 hours to do it.
# Calculator

A calculator with a Go REST API and a React + TypeScript interface. It supports
addition, subtraction, multiplication, division, exponentiation, square root,
and percentage, plus keyboard input and a reusable calculation history.

Go processes every calculator action. React sends actions, formats the returned
numbers as `1.234,56`, and displays the server response.

## Setup

### Docker Desktop

Requirements: Docker Desktop running with Linux containers and Docker Compose.
Local Go and Node.js installations are not required for this option.

Run these commands from the project root:

```sh
docker compose up --build -d --wait
docker compose ps
```

| Service                           | Default address                  |
| --------------------------------- | -------------------------------- |
| Calculator                        | http://localhost:3000            |
| API through the frontend proxy    | http://localhost:3000/api        |
| Backend directly                  | http://localhost:8080            |
| Health through the frontend proxy | http://localhost:3000/api/health |

Compose binds published ports to `127.0.0.1`. The frontend forwards `/api/*`
to the backend container, so the browser uses one origin.

#### Port configuration

Compose reads an optional `.env` in the project root. Defaults are documented
in [.env.example](.env.example):

```dotenv
BACKEND_PORT=8080
FRONTEND_PORT=3000
```

If no `.env` exists, you can copy `.env.example` to `.env` and edit the ports.
Preserve an existing `.env`; it may already contain local overrides.

For example, set `BACKEND_PORT=8081` if a development backend is using 8080.
The direct Docker API then becomes `http://localhost:8081`; the frontend still
uses `http://localhost:3000/api`. The container's internal backend port remains 8080. `FRONTEND_PORT` changes the browser-facing port.

The root `.env` is ignored by Git and sits outside both Docker build contexts.

#### Logs, rebuilds, and shutdown

```sh
docker compose logs -f
docker compose up --build -d --wait
docker compose down
```

Rebuild after changing application code. Docker serves the production frontend;
it does not mount source code or run Vite's development server.

### Local development

Requirements: Go 1.27 or newer, Node.js 24 LTS, and npm.

From the project root, start the backend in one terminal:

```sh
go -C backend run ./cmd/api
```

Start the frontend in another terminal:

```sh
npm --prefix frontend ci
npm --prefix frontend run dev
```

Open http://localhost:5173. Vite forwards `/api/*` to
`http://localhost:8080`. Keep both processes running.

The backend reads `PORT` and defaults to 8080. If you change that port for local
development, update the proxy target in
[frontend/vite.config.ts](frontend/vite.config.ts) as well.
The root Compose `.env` does not configure these local processes.

When working with the project in WSL, run both development commands inside WSL.

## Using the calculator

Enter numbers with buttons or the keyboard. The display uses dots for thousands
and a comma for decimals; either comma or dot on the keyboard enters a decimal.

| Input                   | Action                        |
| ----------------------- | ----------------------------- |
| `0`–`9`                 | Enter digits                  |
| `+`, `-`, `*`, `/`, `^` | Choose an operation           |
| `Enter` or `=`          | Calculate                     |
| `R`                     | Square root                   |
| `%`                     | Percentage                    |
| `Backspace`             | Remove the last character     |
| `Esc` or `Delete`       | Reset the current calculation |

Power toggles off when selected again. Select a history entry to reuse its
result. All clear preserves history; the history panel's Clear button removes it.

## REST API examples

The examples below use Bash/WSL and the default Docker frontend proxy at
`http://localhost:3000`. For local development, use
`http://localhost:8080` instead. Adjust host ports for your configuration.

Requests to POST endpoints require `Content-Type: application/json`.
Operands and results are decimal **strings**, using a dot as the decimal
separator and no thousands separators. Scientific notation is accepted.
For example, send `"1234.56"`; `"1.234,56"` is a display value and is rejected.

### Health

```sh
curl -sS --fail-with-body http://localhost:3000/api/health
```

Response:

```json
{
  "status": "ok"
}
```

### Calculate directly

`POST /api/calculate` evaluates one operation without managing UI state or
adding a history entry.

```sh
curl -sS --fail-with-body http://localhost:3000/api/calculate \
  -H 'Content-Type: application/json' \
  --data-raw '{"operation":"add","left":"12.5","right":"3"}'
```

Response:

```json
{
  "result": "15.5"
}
```

All seven operations:

| Operation      | Request body                                      | Result               |
| -------------- | ------------------------------------------------- | -------------------- |
| Addition       | `{"operation":"add","left":"8","right":"3"}`      | `"11"`               |
| Subtraction    | `{"operation":"subtract","left":"8","right":"3"}` | `"5"`                |
| Multiplication | `{"operation":"multiply","left":"8","right":"3"}` | `"24"`               |
| Division       | `{"operation":"divide","left":"8","right":"3"}`   | `"2.66666666666667"` |
| Exponentiation | `{"operation":"power","left":"2","right":"10"}`   | `"1024"`             |
| Square root    | `{"operation":"sqrt","left":"81"}`                | `"9"`                |
| Percentage     | `{"operation":"percent","left":"25"}`             | `"0.25"`             |

Binary operations require `left` and `right`. Square root and percentage take
only `left`. Percentage optionally accepts `base`:

```sh
curl -sS --fail-with-body http://localhost:3000/api/calculate \
  -H 'Content-Type: application/json' \
  --data-raw '{"operation":"percent","left":"10","base":"200"}'
```

Response: `{"result":"20"}`.

PowerShell example:

```powershell
Invoke-RestMethod -Uri 'http://localhost:3000/api/calculate' `
  -Method Post `
  -ContentType 'application/json' `
  -Body '{"operation":"add","left":"12.5","right":"3"}'
```

### Process a calculator action

`POST /api/calculator/action` is the endpoint used by the UI. Send the previous
server state and one action. This initial request enters the digit 2:

```sh
curl -sS --fail-with-body http://localhost:3000/api/calculator/action \
  -H 'Content-Type: application/json' \
  --data-raw '{
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
    "action": {
      "type": "digit",
      "value": "2"
    }
  }'
```

Response:

```json
{
  "state": {
    "display": "2",
    "stored": null,
    "operator": null,
    "awaitingOperand": false,
    "justCalculated": false,
    "expression": "",
    "operandLabel": null,
    "error": null,
    "history": [],
    "nextId": 1
  }
}
```

Pass the returned `state` unchanged with each following action. For example,
send `{"type":"operator","value":"+"}`, then `{"type":"digit","value":"3"}`,
then `{"type":"equals"}`, each alongside the latest returned state. The display
will become `"5"`, and the response history will include the calculation.

| Action type                                                                           | Value                              |
| ------------------------------------------------------------------------------------- | ---------------------------------- |
| `digit`                                                                               | A single digit string              |
| `operator`                                                                            | `+`, `-`, `*`, `/`, or `^`         |
| `restore`                                                                             | An entry from the returned history |
| `decimal`, `equals`, `clear`, `backspace`, `sign`, `sqrt`, `percent`, `clear-history` | Omit `value`                       |

### Error responses

For example, division by zero through the direct calculation endpoint:

```sh
curl -sS -i http://localhost:3000/api/calculate \
  -H 'Content-Type: application/json' \
  --data-raw '{"operation":"divide","left":"9","right":"0"}'
```

Returns HTTP 422 with:

```json
{
  "error": "Cannot divide by zero."
}
```

| Status | Meaning                                                                   |
| ------ | ------------------------------------------------------------------------- |
| 400    | Malformed JSON, unknown fields, invalid state/action, or invalid operands |
| 404    | Unknown endpoint                                                          |
| 405    | HTTP method is not supported for that endpoint                            |
| 413    | Body exceeds 4 KiB for calculations or 32 KiB for actions                 |
| 415    | Content type is not JSON                                                  |
| 422    | Undefined arithmetic result from `/api/calculate`                         |
| 500    | Unexpected internal failure; implementation details are hidden            |

Accepted UI actions return HTTP 200 even when the calculation has a business
error. In that case, `state.error` contains the message and the UI displays it.
Clients must inspect that field, not just the HTTP status.
Unknown fields and additional JSON documents in one request are rejected.

## Structure and design decisions

```text
backend/
  cmd/api/main.go       Server startup and dependency wiring
  api/
    routes/             Endpoint paths and HTTP method registration
    controllers/        Request decoding, service calls, HTTP error mapping
    services/           Arithmetic, validation, and calculator state rules
    repositories/       Request-scoped calculation history
    models/             State, commands, history entries, and typed errors
    views/              JSON request/response representations
  tests/                Controller, route, service, and repository tests

frontend/
  src/
    controllers/        React state binding and keyboard events
    services/           Ordered requests, cancellation, loading/error handling
    repositories/       REST transport and response validation
    models/             Shared state and action types
    views/              Calculator interface and display formatting
    App.tsx             Connects the controller to the view
  tests/                Tests grouped by layer
```

### Go owns the calculations and input rules

All operations, input edits, operator toggles, percentages, and history actions
run in Go. This gives the REST API and UI one implementation of calculator
behavior. React never evaluates arithmetic or supplies a local fallback.
The tradeoff is that entering a digit also requires a backend request.

### Layers keep responsibilities separate

Routes declare endpoints; controllers translate HTTP requests into service calls;
services apply the rules using models and repositories; backend views render
JSON. Services and repositories do not depend on HTTP.

The frontend uses the same named layers for its own responsibilities. Its
repository sends REST requests, its service orders requests and discards obsolete
responses after reset/history reuse, and its controller connects those results
to the view. Keeping the queue ordered ensures each action uses the state
returned by the preceding action.

### State is stateless on the server

The client holds the latest state in React memory and sends it with each action.
A new in-memory history repository is seeded from that request, retains the last
20 results, and returns updated history. There is no database or global session
store. Separate clients do not share a history; refreshing the page resets it.

This keeps deployment simple and avoids sticky sessions. History is client-supplied
session state, not a persistent or authoritative audit log.

### Precision and calculator behavior

Calculator digit entry is capped at 15 numeric digits. The direct API accepts
longer operand strings (up to 64 characters), but still uses finite `float64`
values. Results round to **15 significant digits**, preserving inputs such
as `123456789012345 + 0`. This is not a fixed number of decimal places or
arbitrary-precision decimal arithmetic; binary floating-point limits still apply.

Calculations run in entry order: `2 + 3 × 4 = 20`. Percentages use the pending
operation's context: `200 + 10% = 220`, while `200 × 10% = 20`.
The direct API's optional `base` calculates a percentage of that base.

Formatting is separate from arithmetic: JSON uses decimal strings with a dot,
while the React view displays Brazilian separators without converting those
strings to JavaScript numbers.

### Docker builds include checks and small runtime images

Each service uses its own directory as a build context and has an allowlist-based
`.dockerignore`. Local dependencies, binaries, generated output, environment
secrets, and unrelated project files stay out of the build context.

Both Dockerfiles use multiple stages. The backend build runs Go tests and vet,
then copies a static binary into a `scratch` image running as a non-root user.
The frontend build uses `npm ci`, runs tests and the TypeScript/Vite build, then
copies only production assets into Nginx. Development tools and source are not
copied into the runtime images.

Nginx serves the UI and proxies API requests to Go. The frontend container's
health check calls the proxied Go health endpoint, so Compose's `--wait` checks
the connection between both services.

## Tests, builds, and formatting

Run from the project root:

```sh
go -C backend test ./...
go -C backend vet ./...
go -C backend build -o bin/server ./cmd/api

npm --prefix frontend test
npm --prefix frontend run build
npm --prefix frontend run format:check
```

Tests live under `backend/tests/` and `frontend/tests/`, separate from the
application source. Go tests cover arithmetic, state transitions, history,
routing, and request validation. Frontend tests cover keyboard mapping, display
formatting, REST transport, request ordering, and cancellation.

To format source:

```sh
gofmt -w backend/api backend/cmd backend/tests
npm --prefix frontend run format
```

Go uses tabs; frontend files use two spaces. The root
[.editorconfig](.editorconfig) and
[frontend/.prettierrc.json](frontend/.prettierrc.json) record those conventions.

For more detail, see the [backend README](backend/README.md) and
[frontend README](frontend/README.md).
