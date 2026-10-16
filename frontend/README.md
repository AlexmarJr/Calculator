# Frontend

React + TypeScript UI for the Go calculator. This directory contains no
arithmetic engine, operation evaluation, percentage rules, input-state rules,
or history-processing rules.

Every button/key action is sent to `POST /api/calculator/action` with the
last state returned by Go. The frontend renders the complete response.
It only defines controls and API types, formats numbers as `1.234,56`,
orders requests, and ignores responses cancelled by reset/history reuse.

The Go backend must run for all calculator input.

## Structure

```text
src/
  controllers/   React state binding and keyboard events
  services/      Request ordering, cancellation, and loading/error state
  repositories/  REST requests and response validation
  models/        Calculator state and action types
  views/         Calculator UI and display formatting
  App.tsx        Connects the controller to the view
tests/
  controllers/
  services/
  repositories/
  views/
```

The view sends actions through the controller and service to the REST repository.
Go returns the complete calculator state. None of these layers evaluates
operations or implements calculator input rules.

## Development

Requires Node.js 24 LTS. Start the Go backend on port 8080 first.

```sh
npm ci
npm run dev
```

Vite proxies `/api/*` to `http://localhost:8080`. Docker's Nginx forwards
those requests to the backend container.

Keyboard: digits, comma/dot for decimals, +, -, *, /, ^, %, Enter or = for
equals, R for square root, Backspace to delete, Esc or Delete to reset.
Use Tab to move between buttons and Space to activate them.

## Checks

```sh
npm test
npm run build
```

Frontend tests cover formatting, transport, response rendering, ordering, and
cancellation. All calculator behavior and arithmetic are tested in `../backend`.

## Formatting

Use two spaces for frontend code and tabs for Go. Editor settings are in
the root `.editorconfig`; frontend rules are in `.prettierrc.json`.

```sh
npm run format
npm run format:check
```
