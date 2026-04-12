# dps-go

Go library for dumb partition storage, backed by SQLite.

## Layout

Each top-level package (e.g. `blob/`) is a self-contained dumb partition
variant with its own interfaces, types, and implementation.

## Tooling

- `make` for all dev workflows
- `golangci-lint` v2 for linting and formatting
- Go 1.26+

**All dev workflows go through `make`:**

| Instead of               | Use          |
|--------------------------|--------------|
| `go build ./...`         | `make build` |
| `go test ./...`          | `make test`  |
| `go mod tidy`            | `make tidy`  |
| `go vet ./...`           | `make tidy`  |
| `golangci-lint run ./...`| `make tidy`  |

`make` runs all targets: build, tidy, test.

## Design

### Remove what doesn't earn its place

Every exported type, function, and option is a maintenance commitment and a
cognitive cost. If two options collapse into one without losing expressiveness,
collapse them. If a mode is the only sensible choice, don't make the user choose
it.

### Contain external dependencies at the seam

SQLite details and cgo boundaries stay behind clean interfaces. The rest of the
code works with plain Go types.

### Resist abstraction until the primitive proves insufficient

Don't create interfaces until there are two implementations. Complexity is added
in response to demonstrated need, not anticipated need.

### Fail at construction, not at use

When an argument is invalid or a configuration is impossible, return an error
from the constructor. Don't let the problem become a stale field that surfaces
three calls later with no indication of where it came from. If a constructor
returns without error, all its fields are known valid.

### Fail explicitly with named errors

Every error a caller might need to handle gets a named sentinel or typed error.
Callers check with `errors.Is` and `errors.As` — not string matching, not type
assertions on unexported types. A missing table should name the table. A
constraint violation should name the constraint.

No silent fallbacks, no swallowed errors.

### Names at the caller's level

Public APIs use names that reflect what the caller is doing, not how the
implementation works. The caller's mental model is the namespace — SQLite
internals stay behind the seam.

## Documentation

The audience is an experienced developer who saw "Dumb Partition Storage" and
wants something obvious. They want to see a working example, recognize the
pattern, and start building.

### Task-first, not type-first

Lead each section with what the reader is trying to do. Show a complete,
copy-pasteable example before explaining the parts. If someone can read the
example and start working, the doc succeeded.

When a pattern enables something (reuse, composition, flexibility), show the
payoff. A reusable function earns its example when called in both contexts.

### Stay at the caller's level

Write about what the caller sees and does. Details only earn space when they
change what the caller writes.

### One concept per section

Each doc section covers one task: basic usage, custom codecs, transactions,
writing generic helpers. A reader scanning headers should find exactly the
section they need.

### Tone

Direct and concise. Assume competence and skip preamble.
