# Chirpy

Chirpy is a lightweight social media backend inspired by X/Twitter. It exposes a RESTful API that supports **GET**, **POST**, **PUT**, and **DELETE** operations, allowing clients to create accounts, authenticate, post chirps, manage tokens, and interact with a PostgreSQL‑backed data store.

---

## Why I Built This Project

This project is part of the Boot.dev curriculum and demonstrates how to build a real HTTP backend in Go. It showcases:

- Running an HTTP server  
- Handling and validating client requests  
- Persisting data in PostgreSQL  
- Structuring a Go backend with clean separation of concerns  
- Using **goose** for database migrations  
- Using **sqlc** to generate type‑safe Go code from SQL queries  
- Exchanging all data between client and server in JSON  

All SQL files, migrations, and generated code are included in the repository.

---

## How to Build and Run Chirpy

### 1. Install Go
Install Go from the official website or through your IDE’s package manager.

### 2. Install PostgreSQL
Download PostgreSQL from the official site and accept the default installation options.

- Set a password for the default superuser `postgres`
- Keep the default port `5432`

### 3. Create a `.env` File
Create a `.env` file in the project root (excluded from Git). Add the following values:

```
DB_URL="postgres://<postgres_user>:<user_password>@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
SECRET="your_secret_string_goes_here-a_random_64_byte_string_will_work"
POLKA_KEY="your_key_goes_here"
```

**Notes:**

- `DB_URL` is your PostgreSQL connection string.  
- macOS uses a slightly different connection format — check PostgreSQL documentation if needed.  
- The `PLATFORM` variable controls access to destructive endpoints (e.g., resetting all users and chirps). Only `dev` environments may call these endpoints.

---

## API Endpoints

Below is a complete list of valid API endpoints supported by Chirpy.

### Health & Metrics

- **GET `/api/healthz`**  
  Returns server health status.

- **GET `/admin/metrics`**  
  Returns the number of requests handled since the last reset.

- **POST `/admin/reset`**  
  Deletes all users and chirps. Only allowed when `PLATFORM=dev`.

---

### User Management

- **POST `/api/users`**  
  Creates a new user.  
  **Body:**  
  ```json
  { "email": "...", "password": "..." }
  ```

- **POST `/api/login`**  
  Logs a user in.  
  **Body:**  
  ```json
  { "email": "...", "password": "..." }
  ```

- **PUT `/api/users`**  
  Updates a user’s email or password.  
  **Body:**  
  ```json
  { "email": "...", "password": "..." }
  ```

---

### Chirps (Posts)

- **POST `/api/chirps`**  
  Creates a new chirp (max 140 characters).  
  Profanity is filtered using `replaceBadWords`.  
  **Body:**  
  ```json
  { "body": "..." }
  ```

- **GET `/api/chirps`**  
  Returns all chirps.  
  Optional query parameters:  
  - `author_id=<uuid>` — filter by author  
  - `sort=asc|desc` — sort by creation time  

- **GET `/api/chirps/{chirpID}`**  
  Returns a single chirp by ID.

- **DELETE `/api/chirps/{chirpID}`**  
  Deletes a chirp owned by the authenticated user.

---

### Tokens

- **POST `/api/refresh`**  
  Issues a refresh token.  
  **Body:**  
  ```json
  { "userID": "..." }
  ```

- **POST `/api/revoke`**  
  Revokes a refresh token.  
  **Body:**  
  ```json
  { "userID": "..." }
  ```

---

### Polka Webhooks

- **POST `/api/polka/webhooks`**  
  Used by a third‑party (Polka) to confirm user upgrades to **Chirpy Red**.  
  Requires an API key in the following header:

  ```
  Authorization: ApiKey <key>
  ```


---

## Development Notes

- All database interactions are type‑safe thanks to **sqlc**.  
- Migrations are managed with **goose**.  
- The server is designed to be stateless and easy to deploy.  
- Sensitive values (DB credentials, secrets, API keys) must be stored in `.env` and never committed.

---