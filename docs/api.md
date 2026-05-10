# Chirpy API Reference

Base URL: `http://localhost:8080`

---

## Authentication

Protected endpoints require a JWT access token passed in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

The Polka webhook endpoint requires an API key:

```
Authorization: ApiKey <polka_key>
```

Access tokens expire after **1 hour**. Use the refresh token flow to obtain a new one without re-authenticating.

---

## Endpoints

### Health

#### `GET /api/healthz`

Returns the server health status.

**Response** `200 OK`

```
OK
```

---

### Users

#### `POST /api/users`

Creates a new user account.

**Request body**

```json
{
  "email": "jane@example.com",
  "password": "supersecret"
}
```

**Response** `201 Created`

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "created_at": "2026-05-11T10:00:00Z",
  "updated_at": "2026-05-11T10:00:00Z",
  "email": "jane@example.com",
  "is_chirpy_red": false,
  "token": "",
  "refresh_token": ""
}
```

**Errors**

| Status | Reason |
|--------|--------|
| `400` | Invalid request body or failed to hash password |
| `500` | Database error |

---

#### `PUT /api/users` `AUTH REQUIRED`

Updates the authenticated user's email and password.

**Headers**

```
Authorization: Bearer <access_token>
```

**Request body**

```json
{
  "email": "newemail@example.com",
  "password": "newpassword"
}
```

**Response** `200 OK`

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "created_at": "2026-05-11T10:00:00Z",
  "updated_at": "2026-05-11T10:05:00Z",
  "email": "newemail@example.com",
  "is_chirpy_red": false,
  "token": "",
  "refresh_token": ""
}
```

**Errors**

| Status | Reason |
|--------|--------|
| `401` | Missing or invalid access token |
| `500` | Database error |

---

### Authentication

#### `POST /api/login`

Authenticates a user and returns a short-lived JWT access token and a long-lived refresh token.

**Request body**

```json
{
  "email": "jane@example.com",
  "password": "supersecret"
}
```

**Response** `200 OK`

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "created_at": "2026-05-11T10:00:00Z",
  "updated_at": "2026-05-11T10:00:00Z",
  "email": "jane@example.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a3f9c2e1b4d7..."
}
```

Token lifetimes:
- `token` (JWT): 1 hour
- `refresh_token`: 60 days

**Errors**

| Status | Reason |
|--------|--------|
| `401` | Invalid email or password |

---

#### `POST /api/refresh`

Issues a new JWT access token using a valid, non-expired refresh token.

**Headers**

```
Authorization: Bearer <refresh_token>
```

**Response** `200 OK`

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors**

| Status | Reason |
|--------|--------|
| `400` | Missing or malformed Authorization header |
| `401` | Refresh token expired or revoked |

---

#### `POST /api/revoke`

Revokes a refresh token, preventing it from being used to issue new access tokens.

**Headers**

```
Authorization: Bearer <refresh_token>
```

**Response** `204 No Content`

**Errors**

| Status | Reason |
|--------|--------|
| `400` | Missing or malformed Authorization header, or database error |

---

### Chirps

#### `GET /api/chirps`

Returns all chirps, ordered by creation time. Optionally filtered by author.

**Query parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `author_id` | UUID (optional) | Filter chirps to a specific user |

**Example request**

```
GET /api/chirps?author_id=a1b2c3d4-e5f6-7890-abcd-ef1234567890
```

**Response** `200 OK`

```json
[
  {
    "id": "c9d8e7f6-a5b4-3210-fedc-ba9876543210",
    "created_at": "2026-05-11T10:01:00Z",
    "updated_at": "2026-05-11T10:01:00Z",
    "body": "Hello, Chirpy world!",
    "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  },
  {
    "id": "d1e2f3a4-b5c6-7890-1234-567890abcdef",
    "created_at": "2026-05-11T10:02:00Z",
    "updated_at": "2026-05-11T10:02:00Z",
    "body": "Another chirp from the same user.",
    "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }
]
```

Returns an empty array `[]` when no chirps are found.

**Errors**

| Status | Reason |
|--------|--------|
| `400` | `author_id` is not a valid UUID |
| `500` | Database error |

---

#### `GET /api/chirps/{chirpID}`

Returns a single chirp by its ID.

**Path parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `chirpID` | UUID | ID of the chirp to retrieve |

**Response** `200 OK`

```json
{
  "id": "c9d8e7f6-a5b4-3210-fedc-ba9876543210",
  "created_at": "2026-05-11T10:01:00Z",
  "updated_at": "2026-05-11T10:01:00Z",
  "body": "Hello, Chirpy world!",
  "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**Errors**

| Status | Reason |
|--------|--------|
| `404` | Chirp not found or invalid UUID |

---

#### `POST /api/chirps` `AUTH REQUIRED`

Creates a new chirp on behalf of the authenticated user. Profanity is automatically filtered.

**Headers**

```
Authorization: Bearer <access_token>
```

**Request body**

```json
{
  "body": "Hello, Chirpy world!"
}
```

Chirps are limited to **140 characters**. Certain profane words are replaced with `****`.

**Response** `201 Created`

```json
{
  "id": "c9d8e7f6-a5b4-3210-fedc-ba9876543210",
  "created_at": "2026-05-11T10:01:00Z",
  "updated_at": "2026-05-11T10:01:00Z",
  "body": "Hello, Chirpy world!",
  "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**Errors**

| Status | Reason |
|--------|--------|
| `400` | Chirp body exceeds 140 characters |
| `401` | Missing or invalid access token |
| `500` | Failed to decode request body |

---

#### `DELETE /api/chirps/{chirpID}` `AUTH REQUIRED`

Deletes a chirp. Only the author of the chirp may delete it.

**Headers**

```
Authorization: Bearer <access_token>
```

**Path parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `chirpID` | UUID | ID of the chirp to delete |

**Response** `204 No Content`

**Errors**

| Status | Reason |
|--------|--------|
| `401` | Missing or invalid access token |
| `403` | Authenticated user is not the chirp author |
| `404` | Chirp not found or invalid UUID |

---

### Webhooks

#### `POST /api/polka/webhooks`

Receives webhook events from Polka. Currently handles the `user.upgraded` event to grant Chirpy Red status.

**Headers**

```
Authorization: ApiKey <polka_key>
```

The API key is validated only for recognized events.

**Request body**

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }
}
```

**Response** `204 No Content`

Returned for both successfully processed events and unrecognized events.

**Errors**

| Status | Reason |
|--------|--------|
| `400` | Invalid request body or malformed user ID |
| `401` | Missing or invalid API key (for recognized events) |
| `404` | User not found |

---

### Admin

#### `GET /admin/metrics`

Returns an HTML page displaying the number of requests served by the fileserver since the last reset.

**Response** `200 OK` — `text/html`

---

#### `POST /admin/reset`

Resets the fileserver hit counter and deletes all users from the database. Only available when the server is running in `dev` platform mode.

**Response** `200 OK` — `text/plain`

```
Reset server count back to 0
```

**Errors**

| Status | Reason |
|--------|--------|
| `401` | Server is not running in `dev` mode |
