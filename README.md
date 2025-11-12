## Auth Bucket Service

Lightweight authentication REST API built with Go and Gin. It provides user creation, login, password reset, and basic user listing backed by PostgreSQL. Configuration is file-driven and supports optional TLS and SMTP for email flows.

### Features
- User auth endpoints: create, login, reset password, list users
- JWT-based request authorization with configurable bypass list
- PostgreSQL connection using pgx pool
- CORS enabled (development-friendly defaults)
- Optional TLS support
- Optional SMTP configuration for email
- sqlc-based query code generation

### Tech stack
- Go, Gin, Logrus, Viper
- pgx/v5 (pgxpool) for PostgreSQL
- sqlc for type-safe query codegen
- Swagger (optional, target present; generation command provided)


## Getting started

### Prerequisites
- Go (recommended 1.21+)
- Make
- sqlc (for regenerating query code) 
- PostgreSQL instance you can reach from your machine


### Configuration
Create or update a JSON config file. A sample is provided at `config/connection/dev-config.json`. Do not commit secrets.

Example (redacted values):

```json
{
  "dbhost": "localhost",
  "dbPort": 5432,
  "dbname": "postgres",
  "dbuid": "postgres",
  "dbpassword": "your-db-password",
  "timeout": 100,
  "connRetryCount": 1,
  "connRetryInterval": 5000,
  "jwtKey": "your-jwt-signing-key",
  "bypassAuth": [
    "/api/auth/create",
    "/api/auth/login",
    "/api/auth/resetpwd",
    "/apidoc/index.html",
    "/apidoc/swagger.yaml"
  ],
  "isTLS": false,
  "tlsKeyPath": "",
  "tlsCertPath": "",
  "senderEmail": "user@gmail.com",
  "password": "app-password-or-token",
  "smtp_host": "smtp.gmail.com",
  "smtp_port": 587,
  "url": {
    "uiurl": "http://localhost:3000"
  }
}
```


## Database Schema

The service uses PostgreSQL and requires the following table in the `common` schema:

```sql
CREATE TABLE common.users (
    user_id serial4 NOT NULL,
    user_name text NOT NULL,
    email text NOT NULL,
    phone text NOT NULL,
    pass text NOT NULL,
    pss_valid bool DEFAULT true NOT NULL,
    otp text NULL,
    otp_valid bool DEFAULT false NOT NULL,
    otp_exp timestamp NULL,
    "role" text NOT NULL
);
```

**Table Structure:**
- `user_id`: Auto-incrementing primary key (serial)
- `user_name`: User's username (text, required)
- `email`: User's email address (text, required)
- `phone`: User's phone number (text, required)
- `pass`: User's password hash (text, required)
- `pss_valid`: Password validity flag (boolean, default: true)
- `otp`: One-time password for verification (text, nullable)
- `otp_valid`: OTP validity flag (boolean, default: false)
- `otp_exp`: OTP expiration timestamp (timestamp, nullable)
- `role`: User role/permissions (text, required)

**Note:** Ensure the `common` schema exists in your PostgreSQL database before creating the table:
```sql
CREATE SCHEMA IF NOT EXISTS common;
```


## Build and run
The `Makefile` provides convenient targets.

- Build for Windows:

```bash
make winBuild
```

- Build and run in development on Windows (default port 7070):

```bash
make dev
```

This compiles `build/exec/service.exe` and runs it with:

```bash
./build/exec/service.exe -c ./config/connection/dev-config.json --port 7070
```

- Cross-compile (Windows and Linux binaries) without running:

```bash
make server
```

- Regenerate sqlc query code:

```bash
make sqlc
```

- Generate Swagger docs (optional; code references present but serving is commented):

```bash
make swag
```


### Run directly (without Make)

```bash
go mod download
go build -o build/exec/service.exe main.go
./build/exec/service.exe -c ./config/connection/dev-config.json --port 7070
```

Flags:
- `-c` path to config JSON (default `./config.json`)
- `--port` server port (default `7070`)
- `-v` verbose logs


## API
Base health:
- `GET /` → 200, basic availability

Auth endpoints:
- `POST /api/auth/create`
- `POST /api/auth/login`
- `POST /api/auth/resetpwd`
- `GET /api/auth/users`

Notes:
- Requests are intercepted by an auth middleware. Paths in `bypassAuth` are accessible without a token.
- Static API docs (if generated/copied) are served from `/apidoc`.


## Development
- CORS is open by default for development. Harden before production.
- sqlc config lives at `config/sqlc/db_query/db.yaml`. Update queries under `config/sqlc/db_query` and run `make sqlc`.
- Generated query code is in `internal/dbmodel/db_query/`.


## Debugging (VS Code)
Use a launch config similar to:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Auth Bucket Service",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/main.go",
      "args": ["-c", "./config/connection/dev-config.json", "--port", "7070"],
      "cwd": "${workspaceFolder}"
    }
  ]
}
```


## Notes
- Keep secrets out of git. Use environment-specific config files.
- For TLS, set `isTLS: true` and provide `tlsKeyPath` and `tlsCertPath`.
- SMTP fields are optional, required only if email flows are enabled.