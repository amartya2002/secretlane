## Secretlane API (current state)

This repo currently focuses on:
- Auth: signup + login with JWT (cookie-based).
- Workspace CRUD, scoped to the authenticated user.
- Swappable DB backend: SQLite (default) or Postgres (via pgx).

Older endpoints for agents, nodes, secrets, and SSH configs are no longer backed by migrations and should be treated as experimental/disabled for now.

## Configuration

Configuration comes from three layers (highest priority first):
1. Environment variables (`.env`)
2. `config.yaml`
3. Built-in defaults

Key settings in `config.yaml`:

```yaml
app:
  port: 8080
  enable_frontend: true
  seed_default_user: true   # creates admin@local / ChangeMe123!

database:
  driver: sqlite            # or postgres

postgres:
  host: localhost
  port: 5432
  user: postgres
  password: ""
  dbname: secretlane
  sslmode: disable
```

Key env vars (see `.env` for full list):
- `PORT` – overrides `app.port`.
- `ENABLE_FRONTEND` – overrides `app.enable_frontend`.
- `SEED_DEFAULT_USER` – overrides `app.seed_default_user`.
- `DB_DRIVER` – overrides `database.driver` (`sqlite` / `postgres`).
- `PGHOST`, `PGPORT`, `PGUSER`, `PGPASSWORD`, `PGDATABASE`, `PGSSLMODE` – Postgres connection.
- `JWT_SECRET` – required, used for signing JWT tokens.

## Running the API

1. Install Go dependencies and tidy:
   ```bash
   go mod tidy
   ```
2. Set up `.env` with at least:
   ```bash
   JWT_SECRET=changeme-super-secret
   # DB_DRIVER=sqlite        # default
   # or DB_DRIVER=postgres   # when Postgres is configured
   ```
3. Start the server:
   ```bash
   go run .
   ```

On startup:
- Config is loaded from `config.yaml` + env.
- DB is initialised in SQLite or Postgres mode.
- Migrations create `users` and `workspaces`.
- If `seed_default_user` is enabled, a default user is added:
  - `username: admin@local`
  - `password: ChangeMe123!`

## API Versioning

All stable endpoints are currently served under:

- `v1` base path: `/api/v1`

The examples below assume the server is running on `http://localhost:8080`.

## Endpoints (v1)

### Signup

Creates a new user and logs them in (sets HttpOnly JWT cookie).

```bash
curl -i -X POST http://localhost:8080/api/v1/signup \
  -H "Content-Type: application/json" \
  -d '{"username": "alice@example.com", "password": "MyPassword123!"}'
```

### Login

Returns JSON and sets HttpOnly JWT cookie.

```bash
curl -i -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin@local", "password": "ChangeMe123!"}'
```

### Logout

Clears the auth cookie.

```bash
curl -i -X POST http://localhost:8080/api/v1/logout
```

### Workspaces (authenticated)

All workspace routes require the JWT cookie from signup/login.

Create workspace:

```bash
curl -i -X POST http://localhost:8080/api/v1/workspaces \
  -H "Content-Type: application/json" \
  --cookie "token=YOUR_JWT_HERE" \
  -d '{"name": "ws1", "description": "first workspace"}'
```

List workspaces:

```bash
curl -i http://localhost:8080/api/v1/workspaces \
  --cookie "token=YOUR_JWT_HERE"
```

Update workspace:

```bash
curl -i -X PUT http://localhost:8080/api/v1/workspaces/1 \
  -H "Content-Type: application/json" \
  --cookie "token=YOUR_JWT_HERE" \
  -d '{"name": "ws1-renamed", "description": "updated description"}'
```

Delete workspace:

```bash
curl -i -X DELETE http://localhost:8080/api/v1/workspaces/1 \
  --cookie "token=YOUR_JWT_HERE"
```
