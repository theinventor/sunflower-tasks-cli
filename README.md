# `sunflower-tasks` — AI-first CLI for Sunflower Tasks

`sunflower-tasks` is the terminal interface to the Sunflower Tasks mobile API. It is designed for AI agents and automation: every command has useful `--help`, successful responses are JSON on stdout, diagnostics are on stderr, credentials are never echoed, and the complete API contract is available as machine-readable context.

## Install

### Pre-built release

Download the matching archive from [GitHub Releases](https://github.com/theinventor/sunflower-tasks-cli/releases), put `sunflower-tasks` on `$PATH`, then run `sunflower-tasks update` to check for future releases.

### Go

```sh
go install github.com/theinventor/sunflower-tasks-cli@latest
```

### From source

```sh
git clone https://github.com/theinventor/sunflower-tasks-cli
cd sunflower-tasks-cli
go build -o ~/.local/bin/sunflower-tasks .
```

## Agent setup

Start with the introspection document instead of scraping help:

```sh
sunflower-tasks agent-context > /tmp/sunflower-tasks-context.json
sunflower-tasks api endpoints --json
```

Login stores a profile in a mode-0600 config file. Passwords are accepted from an environment variable and never persisted:

```sh
export SUNFLOWER_TASKS_PASSWORD='...'
sunflower-tasks auth login --profile work --email you@example.com
sunflower-tasks whoami
```

For CI, use `SUNFLOWER_TASKS_API_KEY` and optionally `SUNFLOWER_TASKS_API_URL` and `SUNFLOWER_TASKS_ACCOUNT_ID`. A saved profile can be selected per invocation with `--profile`.

## API calls

The endpoint catalog documents all 106 operations currently exposed by `/api/mobile/v1`, including inventory counts, maintenance, cleanings, templates, account management, and notifications. `api request` is a forward-compatible escape hatch for new server endpoints:

```sh
sunflower-tasks api request GET /api/mobile/v1/properties
sunflower-tasks api request GET '/api/mobile/v1/schedule?from=2026-10-01&to=2026-10-31'
sunflower-tasks api request PATCH /api/mobile/v1/profile --data '{"preferred_language":"es"}'
sunflower-tasks api request POST /api/mobile/v1/inventories/42/counts --data @count.json
```

The API uses `Authorization: Bearer <token>` and accepts `X-Sunflower-Account-ID` for account selection. Non-2xx responses include the server's JSON error and return a non-zero status.

## Failure behavior

- stdout contains only successful command data, suitable for piping to `jq`.
- stderr contains HTTP status, progress, and actionable errors.
- Exit codes are stable and listed by `agent-context`: usage 2, auth 3, not found 4, API 5, network 6, config 7.
- Tokens and passwords never appear in normal output. Config files are mode 0600.

## Development

```sh
go test ./...
go run . --help
go run . api endpoints --json | jq 'length'
```

Releases are built by GoReleaser for Linux, macOS, and Windows on amd64 and arm64. Release archives include checksums and the README. `update` selects the matching archive for the current platform and refuses to mutate development builds.
