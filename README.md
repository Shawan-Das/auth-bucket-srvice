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

**Satcom Data Table:**

```sql
CREATE TABLE common.satcom_data (
    id serial4 NOT NULL,
    company text NOT NULL,
    category text NOT NULL,
    "type" text NOT NULL,
    "date" text NOT NULL,
    "time" text NOT NULL,
    db_port text NOT NULL,
    ui_port text NOT NULL,
    url text NOT NULL,
    ip text NOT NULL,
    status bool NOT NULL
);
```

**Table Structure:**
- `id`: Auto-incrementing primary key (serial)
- `company`: Company name (text, required)
- `category`: Category classification (text, required)
- `type`: Type of equipment/service (text, required)
- `date`: Date value (text, required)
- `time`: Time value (text, required)
- `db_port`: Database port (text, required)
- `ui_port`: UI port (text, required)
- `url`: URL address (text, required)
- `ip`: IP address (text, required)
- `status`: Active/inactive status (boolean, required)


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

### Base URL
Default: `http://localhost:7070`

### Endpoints

**Public Endpoints (No Authentication Required):**
- `GET /` - Health check
- `POST /api/auth/create` - Create new user
- `POST /api/auth/login` - Login and get JWT token
- `POST /api/auth/resetpwd` - Reset password

**Protected Endpoints (Require JWT Token):**
- `GET /api/auth/users` - Get all users
- `POST /api/satcom` - Create satcom data
- `GET /api/satcom` - Get all satcom data
- `GET /api/satcom/:id` - Get satcom data by ID
- `PUT /api/satcom/:id` - Update satcom data
- `DELETE /api/satcom/:id` - Delete satcom data

**Notes:**
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


## Testing with Postman

### Prerequisites
1. Start the server (see [Build and run](#build-and-run) section)
2. Ensure the database is set up with the required tables
3. Have Postman installed

### Step 1: Create a User (No Token Required)

**Request:**
- **Method:** `POST`
- **URL:** `http://localhost:7070/api/auth/create`
- **Headers:**
  ```
  Content-Type: application/json
  ```
- **Body (raw JSON):**
  ```json
  {
    "email": "test@example.com",
    "password": "testpassword123",
    "phone": "1234567890",
    "userName": "testuser",
    "role": "USER"
  }
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "User created successfully",
  "isSuccess": true,
  "ts": "2024-01-15-10:30:45.123"
}
```

### Step 2: Login and Get JWT Token

**Request:**
- **Method:** `POST`
- **URL:** `http://localhost:7070/api/auth/login`
- **Headers:**
  ```
  Content-Type: application/json
  ```
- **Body (raw JSON):**
  ```json
  {
    "login": "test@example.com",
    "pwd": "testpassword123"
  }
  ```
  > **Note:** The `login` field can be username, email, or phone number.

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Login successful",
  "isSuccess": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "payload": {
    "user_id": 1,
    "user_name": "testuser",
    "email": "test@example.com",
    "phone": "1234567890",
    "role": "USER"
  },
  "ts": "2024-01-15-10:31:00.456"
}
```

**Important:** Copy the `token` value from the response. You'll need it for authenticated requests.

### Step 3: Set Up Postman Environment Variable (Optional but Recommended)

1. In Postman, click on **Environments** (left sidebar)
2. Create a new environment or use the default
3. Add a variable:
   - **Variable:** `auth_token`
   - **Initial Value:** (leave empty)
   - **Current Value:** (leave empty)
4. Save the environment

### Step 4: Use Token for Authenticated Requests

#### Option A: Using Authorization Header (Recommended)

For each authenticated request, add this header:
```
Authorization: Bearer <your-token-here>
```

Replace `<your-token-here>` with the token you received from the login response.

#### Option B: Using Postman Environment Variable

1. After login, in the **Tests** tab of the login request, add:
   ```javascript
   if (pm.response.code === 200) {
       var jsonData = pm.response.json();
       pm.environment.set("auth_token", jsonData.token);
   }
   ```
2. In authenticated requests, set the Authorization header as:
   ```
   Authorization: Bearer {{auth_token}}
   ```

### Step 5: Test Protected Endpoints

#### Get All Users

**Request:**
- **Method:** `GET`
- **URL:** `http://localhost:7070/api/auth/users`
- **Headers:**
  ```
  Authorization: Bearer <your-token-here>
  Content-Type: application/json
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Users retrieved successfully",
  "isSuccess": true,
  "payload": [
    {
      "code": 1,
      "name": "testuser",
      "email": "test@example.com"
    }
  ],
  "ts": "2024-01-15-10:32:00.789"
}
```

#### Create Satcom Data

**Request:**
- **Method:** `POST`
- **URL:** `http://localhost:7070/api/satcom`
- **Headers:**
  ```
  Authorization: Bearer <your-token-here>
  Content-Type: application/json
  ```
- **Body (raw JSON):**
  ```json
  {
    "company": "SatCom Inc",
    "category": "Network",
    "type": "Router",
    "date": "2024-01-15",
    "time": "10:30:00",
    "db_port": "5432",
    "ui_port": "8080",
    "url": "https://example.com",
    "ip": "192.168.1.100",
    "status": true
  }
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Satcom data created successfully",
  "isSuccess": true,
  "ts": "2024-01-15-10:33:00.123"
}
```

#### Get All Satcom Data

**Request:**
- **Method:** `GET`
- **URL:** `http://localhost:7070/api/satcom`
- **Headers:**
  ```
  Authorization: Bearer <your-token-here>
  Content-Type: application/json
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Satcom data retrieved successfully",
  "isSuccess": true,
  "payload": [
    {
      "id": 1,
      "company": "SatCom Inc",
      "category": "Network",
      "type": "Router",
      "date": "2024-01-15",
      "time": "10:30:00",
      "db_port": "5432",
      "ui_port": "8080",
      "url": "https://example.com",
      "ip": "192.168.1.100",
      "status": true
    }
  ],
  "ts": "2024-01-15-10:34:00.456"
}
```

#### Get Satcom Data by ID

**Request:**
- **Method:** `GET`
- **URL:** `http://localhost:7070/api/satcom/1`
- **Headers:**
  ```
  Authorization: Bearer <your-token-here>
  Content-Type: application/json
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Satcom data retrieved successfully",
  "isSuccess": true,
  "payload": {
    "id": 1,
    "company": "SatCom Inc",
    "category": "Network",
    "type": "Router",
    "date": "2024-01-15",
    "time": "10:30:00",
    "db_port": "5432",
    "ui_port": "8080",
    "url": "https://example.com",
    "ip": "192.168.1.100",
    "status": true
  },
  "ts": "2024-01-15-10:35:00.789"
}
```

#### Update Satcom Data

**Request:**
- **Method:** `PUT`
- **URL:** `http://localhost:7070/api/satcom/1`
- **Headers:**
  ```
  Authorization: Bearer <your-token-here>
  Content-Type: application/json
  ```
- **Body (raw JSON):**
  ```json
  {
    "company": "SatCom Inc Updated",
    "category": "Network",
    "type": "Router",
    "date": "2024-01-15",
    "time": "11:00:00",
    "db_port": "5432",
    "ui_port": "8080",
    "url": "https://example.com",
    "ip": "192.168.1.100",
    "status": false
  }
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Satcom data updated successfully",
  "isSuccess": true,
  "ts": "2024-01-15-10:36:00.123"
}
```

#### Delete Satcom Data

**Request:**
- **Method:** `DELETE`
- **URL:** `http://localhost:7070/api/satcom/1`
- **Headers:**
  ```
  Authorization: Bearer <your-token-here>
  Content-Type: application/json
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Satcom data deleted successfully",
  "isSuccess": true,
  "ts": "2024-01-15-10:37:00.456"
}
```

### Step 6: Reset Password (No Token Required)

**Request:**
- **Method:** `POST`
- **URL:** `http://localhost:7070/api/auth/resetpwd`
- **Headers:**
  ```
  Content-Type: application/json
  ```
- **Body (raw JSON):**
  ```json
  {
    "email": "test@example.com",
    "newPwd": "newpassword123"
  }
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Password reset successfully",
  "isSuccess": true,
  "ts": "2024-01-15-10:38:00.789"
}
```

### Troubleshooting

**401 Unauthorized / "Unauthorized" response:**
- Check that you've included the `Authorization: Bearer <token>` header
- Verify the token is still valid (tokens expire after 1 hour)
- Make sure there's no extra spaces in the token
- Try logging in again to get a fresh token

**400 Bad Request:**
- Verify the JSON body is valid
- Check that all required fields are present
- Ensure field names match exactly (case-sensitive)

**404 Not Found:**
- Verify the endpoint URL is correct
- Check that the resource ID exists (for GET/PUT/DELETE by ID)
- Ensure the server is running on the correct port

**500 Internal Server Error:**
- Check server logs for detailed error messages
- Verify database connection is working
- Ensure database tables exist and schema is correct


## Notes
- Keep secrets out of git. Use environment-specific config files.
- For TLS, set `isTLS: true` and provide `tlsKeyPath` and `tlsCertPath`.
- SMTP fields are optional, required only if email flows are enabled.