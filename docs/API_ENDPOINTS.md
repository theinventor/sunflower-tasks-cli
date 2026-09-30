# Sunflower Tasks mobile API

The executable is the canonical, versioned endpoint catalog. Run `sunflower-tasks api endpoints --json` to obtain the exact machine-readable list. The catalog is generated from the Rails mobile routes and includes method, path, summary, and expected body fields where known.

All endpoints under `/api/mobile/v1` except session creation, registration, confirmation, and password reset require a bearer token. Use `X-Sunflower-Account-ID` when an account other than the user's selected account is needed.

The generic command is intentionally stable:

```sh
sunflower-tasks api request METHOD PATH [--query key=value] [--data JSON|@file|@-]
```

This keeps agents usable while the server adds endpoints between CLI releases.
