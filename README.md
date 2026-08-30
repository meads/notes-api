# Notes App

This application achieves two primary objectives. One is to serve as a template for horizontally layered architecture written in Go and two is to provide an example of a backend that authenticates web clients using a token based architecture which is described in the subsequent sections of this document.

This document describes the stateless-but-revocable JWT authentication system used in this Notes application. The architecture prioritizes simplicity, clean database state, and seamless user experience across browser sessions.

The application allows a user to register using a username and password. Following the initial signup the user can create, read, update and delete notes. The functionality was intentionally made to be simple so as just to provide the proof of concept.

## 🛠️ Tech Stack

* **Frontend:** [Vite/React]
* **Styling:** [CSS]
* **HTTP Client:** [Fetch API]
* **HTTP Server:** [Gin]
* **Backend:** [Go] (version 1.26)
* **Tokens:** [JWT]
* **Database:** [sqlc/Postgresql]

## 📦 Getting Started

### Prerequisites
* Docker Desktop [typical docker install steps](https://www.docker.com/get-started/)
* A running instance of the [Notes API](https://github.com/meads/notes-api)


### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/meads/notes-api.git
    cd notes-api
    ```

2. Configure environment variables:
   Create a `.env` file in the root directory using the example found in CONTRIBUTING.md

3. Build and start Docker containers:
   ```bash
   docker compose build 
   docker compose up
   ```


### Running the App

Open `http://localhost:3000` in your browser to view the UI.



<hr>


# 🏗️ Architecture Overview

The system uses a dual-token architecture (Access Token + Refresh Token) persisted through http only cookies. 

To maintain control over active sessions without checking the database on every single API request, the system uses a Sessions table. The sessions table PRIMARY KEY is a uuid that mirrors the refresh token RegisteredClaims.ID. The user_id represents the relation to the user the session is for. The refresh_token is the actual JWT refresh token string. The is_revoked column is a boolean that allows revocation of a session/refresh_token. The expires_at is another field that mirrors the JWT refresh token but for the expires timestamp.

### 🗄️ Database Schema

```sql
CREATE TABLE sessions (
    id VARCHAR(36) PRIMARY KEY NOT NULL,  -- Unique UUID or Session ID
    user_id BIGINT NOT NULL,              -- Foreign key to Users table
    refresh_token VARCHAR(512) NOT NULL,       -- The issued 24-hour Refresh JWT
    is_revoked BOOLEAN DEFAULT FALSE NOT NULL, -- Flag to invalidate the session preventing token rotation
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### 🔄 Lifecycle Flows

#### Registration & Authentication (/register & /login)

  1. The client submits credentials.

  2. The backend verifies credentials.

  3. The backend generates a 15-minute Access Token cookie and a 24-hour Refresh Token cookie.

  4. The backend inserts a new record into the sessions table.

  5. The backend returns the following JSON payload to the client:

```json
{
  "sessionId": "uuid-string",
  "userId": 123,
}
```

#### Client-Side Storage

The client stores all returned properties directly in sessionStorage. Both access_token and refresh_token are 
unreachable by javascript but are persisted in transport via credentials: 'include' option in fetch requests.

```javascript
const data = await response.json()
sessionStorage.setItem("sessionId", data.sessionId)
sessionStorage.setItem("userId", data.userId)
```

#### Requests to protected routes

1. The client makes a request to create, save, delete or fetch some note(s) from the api. The fetch options include 
the access_token in each request via a credentials option being specified.
  ```javascript
  fetch(url,{
  credentials: 'include'
  ...
  })
  ```

2. The access token passed in the http only cookie is seen to be invalid by the backend
  JWT middleware and a 401 response is returned to the client before completing the api request.

3. The client queues the original request and any subsequent requests, then begins the transparent 
  token refresh. See below...

#### Transparent Token Refresh (The Interceptor Queue)

To keep the user logged in seamlessly, the client wraps the native browser fetch function by a function named (fetchClient). The fetchClient function intercepts all 401 responses and performs this refresh/retry logic.

```
Expired Access Token Triggered -> 401 Unauthorized
                                      |
                         Queue all pending requests
                                      |
                     POST /refresh (Send Refresh Token via credentials: 'include' option)
                                      |
              [Backend verifies Session exists in DB]
                                      |
                     Returns NEW 15-min Access Token
                                      |
              Replay queued requests with updated cookie sent to the client 
```

 * Failure Catch: If the /refresh call fails (e.g., token expired or deleted from DB), the request queue is rejected, and the user is redirected to the login screen.

#### Explicit Logout (/logout)

  1. The user clicks "Log Out".
  2. The client fires a POST /logout passing the sessionId.
  3. The backend explicitly deletes the corresponding row from the sessions table.
  4. The client is redirected to the login screen.

### 🧹 Automated Database Cleanup

Because a new session is created upon every new login, the server is never notified of "abandoned" sessions. Left unchecked, the sessions table would grow indefinitely.
To prevent this, create a scheduled routine that runs once every 24 hours to prune expired sessions from the database.


SQL Cleanup Query

Any session older than 24 hours is guaranteed to have an expired Refresh Token and is safe to delete.

```sql
DELETE FROM sessions 
WHERE created_at < NOW() - INTERVAL '24 hours';
```

### ⚠️ Security Warning: CSRF Protection in Dual-Token Cookie Authentication

When implementing a **dual-token architecture** (using an **Access Token** and a **Refresh Token** stored in cookies), your application is inherently vulnerable to **Cross-Site Request Forgery (CSRF)**. Because browsers automatically include cookies with every matching request, a malicious third-party site can force a user's browser to perform unauthorized actions on your platform. 

To secure this architecture against CSRF attacks, you **must** strictly enforce the following security precautions:

### 1. Enforce strict Cookie Attributes
Do not rely on default browser behaviors. Explicitly configure your token cookies with these security flags:
*   **`SameSite=Strict` or `SameSite=Lax`**: Set your Access Token cookie to `Strict` whenever possible to prevent it from being sent during cross-site navigations. Use `Lax` only if top-level cross-site navigations must remain authenticated.
*   **`HttpOnly`**: This prevents Cross-Site Scripting (XSS) from reading the tokens, though it does not block CSRF on its own.
*   **`Secure`**: Ensures cookies are only transmitted over encrypted (HTTPS) connections.

### 2. Implement Anti-CSRF Tokens (Synchronizer Token Pattern)
Because `SameSite` flags can be bypassed in older browsers or through specific sub-domain vulnerabilities, you must implement an explicit defense-in-depth mechanism:
*   **Generate a cryptographically secure, random CSRF token** tied to the user's current session.
*   **Deliver the CSRF token** to the client via a custom HTTP response header or a non-HttpOnly cookie (separate from the authentication cookies).
*   **Require the client application to read this token** and include it in a custom request header (e.g., `X-CSRF-Token`) for all state-changing requests (POST, PUT, DELETE, PATCH).
*   **Verify the header on the server**. Reject any state-changing request where the header is missing, malformed, or does not match the expected session value.

### 3. Require Custom Request Headers
Modern browsers enforce a **Same-Origin Policy (SOP)** for custom headers. 
*   Ensure your backend api rejects state-changing requests that lack standard custom headers like `X-Requested-With` or your specific `X-CSRF-Token`.
*   A malicious site cannot easily forge custom headers via standard HTML forms or cross-origin `fetch`/`xhr` requests unless your CORS policy explicitly allows it.

### 4. Maintain a Strict Cross-Origin Resource Sharing (CORS) Policy
A misconfigured CORS policy can completely neutralize your CSRF defenses.
*   **Never** use `Access-Control-Allow-Origin: *` alongside `Access-Control-Allow-Credentials: true`.
*   **Explicitly whitelist** trusted origins. Do not dynamically mirror the incoming `Origin` header in the response.
*   Restrict allowed methods and headers to the absolute minimum required by your client application.

### 5. Separate Refresh Token and Access Token Lifecycles
*   **The Access Token** should have a very short lifespan (e.g., 5 to 15 minutes) to minimize the window of opportunity for an active CSRF exploit.
*   **The Refresh Token** endpoint (used to generate a new Access Token) **must also be protected against CSRF**. If a malicious site can silently hit your `/refresh` endpoint, they can keep an expired session alive indefinitely. Use strict `SameSite=Strict` attributes and anti-CSRF tokens on the refresh route.

### 📄 License

This project is open-source and available under the **MIT License**.

