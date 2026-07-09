# Public Calls Portal — API Contract

> **Version:** 1.2.0
> **Base URL:** `http://localhost:8080`
> **Auth:** JWT Bearer Token (except where noted)
> **Content-Type:** `application/json`

---

## Endpoints Overview

| # | Method | Path | Auth Required | Description |
|---|--------|------|:---:|-------------|
| 1 | `GET` | `/health` | ❌ | Health check |
| 2 | `POST` | `/user-create` | ❌ | Register a new user |
| 3 | `POST` | `/login` | ❌ | Authenticate and get JWT |
| 4 | `POST` | `/user-info` | ✅ | Get current user profile |
| 5 | `POST` | `/user-modify` | ✅ | Update current user profile fields |
| 6 | `GET` | `/public-proposals` | ✅ | List public proposals from datos.gov.co |
| 7 | `GET` | `/saved_proposals` | ✅ | List user's saved proposals (with full data) |
| 8 | `POST` | `/saved-proposals` | ✅ | Save a proposal for current user |
| 9 | `POST` | `/proposal-save` | ✅ | Save (upsert) a full proposal locally |
| 10 | `POST` | `/logout` | ✅ | Invalidate current session (logout) |

---

## 1. Health Check

**Protocol:** `GET /health`

Check if the backend service is running.

### Input

None.

### Response — `200 OK`

```json
{
  "status": "ok"
}
```

---

## 2. Create User

**Protocol:** `POST /user-create`

Register a new user in the portal. No authentication required.

### Input Parameters

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `id` | string | ✅ | Document number (unique identifier) |
| `first_name` | string | ✅ | User's first name |
| `last_name` | string | ✅ | User's last name |
| `gender` | string | ❌ | Gender |
| `email` | string | ✅ | Email address (must be unique) |
| `phone_number` | string | ❌ | Phone number |
| `password` | string | ✅ | Password (will be bcrypt-hashed) |
| `status` | string | ❌ | Default: `"AC"` (Active) |

### Example Request

```json
{
  "id": "1234567890",
  "first_name": "Sebastián",
  "last_name": "Paniagua",
  "gender": "M",
  "email": "sebastian@example.com",
  "phone_number": "3001234567",
  "password": "MiClaveSegura2026",
  "status": "AC"
}
```

### Response — `201 Created`

```json
{
  "id": "1234567890",
  "first_name": "Sebastián",
  "last_name": "Paniagua",
  "gender": "M",
  "email": "sebastian@example.com",
  "phone_number": "3001234567",
  "status": "AC",
  "created_at": "2026-07-06T14:30:00.000000-05:00",
  "updated_at": "2026-07-06T14:30:00.000000-05:00",
  "deleted_at": null
}
```

### Error Responses

| Status | Condition | Body |
|--------|-----------|------|
| `409 Conflict` | Duplicate document number | `{"error":"document already exists"}` |
| `409 Conflict` | Duplicate email | `{"error":"email already exists"}` |
| `400 Bad Request` | Missing required field | `{"error":"first_name is required"}` |

---

## 3. Login

**Protocol:** `POST /login`

Authenticate with email and password. Returns a JWT token valid for **1 hour**.

### Input Parameters

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `email` | string | ✅ | Registered email |
| `password` | string | ✅ | Account password |

### Example Request

```json
{
  "email": "sebastian@example.com",
  "password": "MiClaveSegura2026"
}
```

### Response — `200 OK`

```json
{
  "token": "eyJhbG...c123...",
  "user": {
    "id": "1234567890",
    "first_name": "Sebastián",
    "last_name": "Paniagua",
    "gender": "M",
    "email": "sebastian@example.com",
    "phone_number": "3001234567",
    "status": "AC",
    "created_at": "2026-07-06T14:30:00.000000-05:00",
    "updated_at": "2026-07-06T14:30:00.000000-05:00",
    "deleted_at": null
  }
}
```

### Error Responses

| Status | Condition | Body |
|--------|-----------|------|
| `401 Unauthorized` | Invalid email or password | `{"error":"invalid email or password"}` |
| `400 Bad Request` | Invalid JSON body | `{"error":"invalid request body"}` |

---

## 4. Get User Info

**Protocol:** `POST /user-info`

**Auth:** `Authorization: Bearer <jwt-token>` — returns the authenticated user's profile information. Password is never included in the response.

### Input

| Header | Value |
|--------|-------|
| `Authorization` | `Bearer <jwt-token>` |

No body parameters required.

### Response - 200 OK

```json
{
  "id": "1234567890",
  "first_name": "Sebastian",
  "last_name": "Paniagua",
  "gender": "M",
  "email": "sebastian@example.com",
  "phone_number": "3001234567",
  "status": "AC",
  "created_at": "2026-07-06T14:30:00.000000-05:00",
  "updated_at": "2026-07-06T14:30:00.000000-05:00",
  "deleted_at": null
}
```

### Error Responses

| Status | Condition | Body |
|--------|-----------|------|
| 401 Unauthorized | Missing/invalid/expired token | `{"error":"Credentials invalid"}` |
| 404 Not Found | User not found | `{"error":"Usuario no registrado"}` |

---

## 5. Modify User

**Protocol:** `POST /user-modify`

**Auth:** `Authorization: Bearer <jwt-token>` — updates the authenticated user's profile. Only the fields provided are updated.

### Input Parameters (all optional)

| Field | Type | Description |
|-------|------|-------------|
| `first_name` | string | New first name |
| `last_name` | string | New last name |
| `gender` | string | New gender |
| `phone_number` | string | New phone number |
| `status` | string | New status (AC / IN) |

> **Note:** The fields `id`, `email`, and `password` cannot be modified through this endpoint. If provided, they will be ignored.

### Validation Rules

- At least one allowed field must be provided.
- If no allowed fields are present, the request is rejected.

### Example Request

```json
{
  "first_name": "Sebastian Alejandro",
  "phone_number": "3009876543"
}
```

### Response - 200 OK

```json
{
  "message": "Informacion actualizada correctamente"
}
```

### Error Responses

| Status | Condition | Body |
|--------|-----------|------|
| 401 Unauthorized | Missing/invalid/expired token | `{"error":"Credentials invalid"}` |
| 400 Bad Request | No allowed fields provided | `{"error":"Debe modificarse un campo por lo menos para actualizar"}` |
| 400 Bad Request | Invalid JSON body | `{"error":"invalid request body"}` |

---

## 6. List Public Proposals

**Protocol:** `GET /public-proposals`

**Auth:** `Authorization: Bearer <jwt-token>`

Fetches public calls for proposals **live from datos.gov.co** (SECOP II SODA API). Does not cache — each request calls the external API.

### Input

| Header | Value |
|--------|-------|
| `Authorization` | `Bearer <jwt-token>` |

### Query Parameters (all optional)

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `query` | string | Full-text search on procedure name | `?query=convocatoria` |
| `fase` | string | Filter by phase (exact match) | `?fase=Presentación de oferta` |
| `entidad` | string | Filter by entity name (partial) | `?entidad=DANE` |
| `limit` | integer | Max results (default: `100`) | `?limit=10` |
| `offset` | integer | Pagination offset (default: `0`) | `?offset=0` |

### Response — `200 OK`

Returns a JSON array of objects. Each object uses the **original column names from datos.gov.co**.

```json
[
  {
    "entidad": "DEPARTAMENTO ADMINISTRATIVO NACIONAL DE ESTADISTICA (DANE)",
    "nit_entidad": "899999027",
    "departamento_entidad": "Distrito Capital de Bogotá",
    "ciudad_entidad": "Bogotá",
    "ordenentidad": "Nacional",
    "id_del_proceso": "CO1.REQ.2577563",
    "referencia_del_proceso": "EDP-545-2022",
    "nombre_del_procedimiento": "Convocatoria DSNFT-0001-FEEC-2024",
    "descripci_n_del_procedimiento": "Aunar esfuerzos entre las partes para ejecutar el proyecto...",
    "fase": "Presentación de oferta",
    "fecha_de_publicacion_del": "2024-02-19T00:00:00.000",
    "fecha_de_ultima_publicaci": "2024-03-12T00:00:00.000",
    "modalidad_de_contratacion": "Licitación pública",
    "precio_base": "500000000.00",
    "duracion": "90",
    "unidad_de_duracion": "Días",
    "fecha_de_recepcion_de": "2024-04-15T17:00:00.000",
    "estado_del_procedimiento": "Publicado",
    "adjudicado": "No",
    "nombre_del_proveedor": null,
    "valor_total_adjudicacion": null,
    "urlproceso": "https://community.secop.gov.co/...",
    "codigo_principal_de_categoria": "72101500",
    "tipo_de_contrato": "Servicios",
    "estado_de_apertura_del_proceso": "Abierto",
    "estado_resumen": "Recepcion de ofertas",
    "proveedores_invitados": "0",
    "proveedores_que_manifestaron": "0",
    "respuestas_al_procedimiento": "0",
    "numero_de_lotes": "1"
  }
]
```

### Error Responses

| Status | Condition | Body |
|--------|-----------|------|
| `401 Unauthorized` | Missing/invalid/expired token | `{"error":"Credentials invalid"}` |
| `500 Internal Server Error` | External API failure | `{"error":"failed to fetch proposals: ..."}` |

---

## 7. List Saved Proposals (with full proposal data)

**Protocol:** `GET /saved_proposals`

**Auth:** `Authorization: Bearer <jwt-token>`

Returns the authenticated user's saved proposals **with full proposal data** (merged from the local proposals store). Returns an empty array `[]` if the user has no saved proposals.

### Input

| Header | Value |
|--------|-------|
| `Authorization` | `Bearer <jwt-token>` |

No query or body parameters.

### Response - 200 OK

Returns an array of saved proposals with the **full proposal data shape**. Each item includes the `user_id` and all proposal fields from `public_calls_proposals`, with the proposal ID serialised as `id_del_proceso`. Association metadata (`id`, `public_call_id`, `association_date`, `created_at`, `updated_at`, `deleted_at`) is not included — use the proposal data directly. Returns an empty array `[]` if the user has no saved proposals.

```json
[ 
  {
    "user_id": "1234567890",
    "entidad": "DEPARTAMENTO ADMINISTRATIVO NACIONAL DE ESTADISTICA (DANE)",
    "nit_entidad": "899999027",
    "departamento_entidad": "Distrito Capital de Bogotá",
    "ciudad_entidad": "Bogotá",
    "ordenentidad": "Nacional",
    "id_del_proceso": "CO1.REQ.2577563",
    "referencia_del_proceso": "EDP-545-2022",
    "nombre_del_procedimiento": "Convocatoria DSNFT-0001-FEEC-2024",
    "descripci_n_del_procedimiento": "Aunar esfuerzos entre las partes para ejecutar el proyecto...",
    "fase": "Presentación de oferta",
    "fecha_de_publicacion_del": "2024-02-19T00:00:00.000",
    "fecha_de_ultima_publicaci": "2024-03-12T00:00:00.000",
    "modalidad_de_contratacion": "Licitación pública",
    "precio_base": "500000000.00",
    "duracion": "90",
    "unidad_de_duracion": "Días",
    "fecha_de_recepcion_de": "2024-04-15T17:00:00.000",
    "estado_del_procedimiento": "Publicado",
    "adjudicado": "No",
    "nombre_del_proveedor": null,
    "valor_total_adjudicacion": null,
    "urlproceso": "https://community.secop.gov.co/...",
    "codigo_principal_de_categoria": "72101500",
    "tipo_de_contrato": "Servicios",
    "estado_de_apertura_del_proceso": "Abierto",
    "estado_resumen": "Recepcion de ofertas",
    "proveedores_invitados": "0",
    "proveedores_que_manifestaron": "0",
    "respuestas_al_procedimiento": "0",
    "numero_de_lotes": "1"
  }
]
```

### Empty Response (no saved proposals)

```json
[]
```

### Error Responses

| Status | Condition | Body |
|--------|-----------|------|
| 401 Unauthorized | Missing/invalid/expired token | `{"error":"Credentials invalid"}` |

---

## 8. Save a Proposal (associate)

**Protocol:** `POST /saved-proposals`

**Auth:** `Authorization: Bearer <jwt-token>`

Saves (associates) a public proposal for the authenticated user. This endpoint is **idempotent** — if the proposal is already saved, it still returns a success response without creating duplicate records.

### Input

| Header | Value |
|--------|-------|
| `Authorization` | `Bearer <jwt-token>` |
| `Content-Type` | `application/json` |

### Input Parameters

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `public_call_id` | string | Yes | ID of the proposal to save |

### Example Request

```json
{
  "public_call_id": "CO1.REQ.2577563"
}
```

### Response - 201 Created

```json
{
  "message": "Guardado satisfactoriamente"
}
```

### Error Responses

| Status | Condition | Body |
|--------|-----------|------|
| 401 Unauthorized | Missing/invalid/expired token | `{"error":"Credentials invalid"}` |
| 400 Bad Request | Invalid JSON body | `{"error":"invalid request body"}` |
| 400 Bad Request | Missing public_call_id | `{"error":"public_call_id is required"}` |

---

## 9. Save a Proposal (full data)

**Protocol:** `POST /proposal-save`

**Auth:** `Authorization: Bearer <jwt-token>`

Saves (upserts) a full public call proposal into the local repository. If a proposal with the same `id` already exists, it updates the existing record; otherwise it creates a new one.

Use this endpoint to preserve proposal details locally so they can be referenced later via the saved-proposals listing (which merges full proposal data).

### Input

| Header | Value |
|--------|-------|
| `Authorization` | `Bearer <jwt-token>` |
| `Content-Type` | `application/json` |

### Input Parameters

All fields matching the datos.gov.co response shape are accepted. Only `id` is required.

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `id` | string | ✅ | Unique proposal identifier |
| `nombre_del_procedimiento` | string | ❌ | Procedure name |
| `entidad` | string | ❌ | Entity name |
| `nit_entidad` | string | ❌ | Entity NIT |
| `departamento_entidad` | string | ❌ | Entity department |
| `ciudad_entidad` | string | ❌ | Entity city |
| `ordenentidad` | string | ❌ | Entity order |
| `referencia_del_proceso` | string | ❌ | Process reference |
| `descripci_n_del_procedimiento` | string | ❌ | Procedure description |
| `fase` | string | ❌ | Phase |
| `fecha_de_publicacion_del` | string | ❌ | Publication date |
| `fecha_de_ultima_publicaci` | string | ❌ | Last publication date |
| `modalidad_de_contratacion` | string | ❌ | Contracting modality |
| `precio_base` | number | ❌ | Base price |
| `duracion` | string | ❌ | Duration |
| `unidad_de_duracion` | string | ❌ | Duration unit |
| `fecha_de_recepcion_de` | string | ❌ | Reception date |
| `estado_del_procedimiento` | string | ❌ | Procedure status |
| `adjudicado` | string | ❌ | Awarded indicator |
| `nombre_del_proveedor` | string | ❌ | Provider name |
| `valor_total_adjudicacion` | number | ❌ | Total award value |
| `urlproceso` | string | ❌ | Process URL |
| `codigo_principal_de_categoria` | string | ❌ | Main category code |
| `tipo_de_contrato` | string | ❌ | Contract type |
| `estado_de_apertura_del_proceso` | string | ❌ | Opening status |
| `estado_resumen` | string | ❌ | Summary status |
| `proveedores_invitados` | string | ❌ | Invited providers count |
| `proveedores_que_manifestaron` | string | ❌ | Providers who expressed interest |
| `respuestas_al_procedimiento` | string | ❌ | Response count |
| `numero_de_lotes` | integer | ❌ | Number of lots |

### Example Request

```json
  {
    "entidad": "DEPARTAMENTO ADMINISTRATIVO NACIONAL DE ESTADISTICA (DANE)",
    "nit_entidad": "899999027",
    "departamento_entidad": "Distrito Capital de Bogotá",
    "ciudad_entidad": "Bogotá",
    "ordenentidad": "Nacional",
    "id": "CO1.REQ.2577563",
    "referencia_del_proceso": "EDP-545-2022",
    "nombre_del_procedimiento": "Convocatoria DSNFT-0001-FEEC-2024",
    "descripci_n_del_procedimiento": "Aunar esfuerzos entre las partes para ejecutar el proyecto...",
    "fase": "Presentación de oferta",
    "fecha_de_publicacion_del": "2024-02-19T00:00:00.000",
    "fecha_de_ultima_publicaci": "2024-03-12T00:00:00.000",
    "modalidad_de_contratacion": "Licitación pública",
    "precio_base": "500000000.00",
    "duracion": "90",
    "unidad_de_duracion": "Días",
    "fecha_de_recepcion_de": "2024-04-15T17:00:00.000",
    "estado_del_procedimiento": "Publicado",
    "adjudicado": "No",
    "nombre_del_proveedor": null,
    "valor_total_adjudicacion": null,
    "urlproceso": "https://community.secop.gov.co/...",
    "codigo_principal_de_categoria": "72101500",
    "tipo_de_contrato": "Servicios",
    "estado_de_apertura_del_proceso": "Abierto",
    "estado_resumen": "Recepcion de ofertas",
    "proveedores_invitados": "0",
    "proveedores_que_manifestaron": "0",
    "respuestas_al_procedimiento": "0",
    "numero_de_lotes": "1"
  }
```

### Response — `201 Created`

Returns the saved (or updated) proposal object.

```json
{
    "id": "CO1.REQ.2577563",
    "entidad": "DEPARTAMENTO ADMINISTRATIVO NACIONAL DE ESTADISTICA (DANE)",
    "nit_entidad": "899999027",
    "departamento_entidad": "Distrito Capital de Bogotá",
    "ciudad_entidad": "Bogotá",
    "ordenentidad": "Nacional",
    "referencia_del_proceso": "EDP-545-2022",
    "nombre_del_procedimiento": "Convocatoria DSNFT-0001-FEEC-2024",
    "descripci_n_del_procedimiento": "Aunar esfuerzos entre las partes para ejecutar el proyecto...",
    "fase": "Presentación de oferta",
    "fecha_de_publicacion_del": "2024-02-19T00:00:00.000",
    "fecha_de_ultima_publicaci": "2024-03-12T00:00:00.000",
    "modalidad_de_contratacion": "Licitación pública",
    "precio_base": "500000000.00",
    "duracion": "90",
    "unidad_de_duracion": "Días",
    "fecha_de_recepcion_de": "2024-04-15T17:00:00.000",
    "estado_del_procedimiento": "Publicado",
    "adjudicado": "No",
    "nombre_del_proveedor": null,
    "valor_total_adjudicacion": null,
    "urlproceso": "https://community.secop.gov.co/...",
    "codigo_principal_de_categoria": "72101500",
    "tipo_de_contrato": "Servicios",
    "estado_de_apertura_del_proceso": "Abierto",
    "estado_resumen": "Recepcion de ofertas",
    "proveedores_invitados": "0",
    "proveedores_que_manifestaron": "0",
    "respuestas_al_procedimiento": "0",
    "numero_de_lotes": "1"
}
```

### Error Responses

| Status | Condition | Body |
|--------|-----------|------|
| 401 Unauthorized | Missing/invalid/expired token | `{"error":"Credentials invalid"}` |
| 400 Bad Request | Invalid JSON body | `{"error":"invalid request body"}` |
| 400 Bad Request | Missing `id` field | `{"error":"id is required"}` |
| 500 Internal Server Error | Storage failure | `{"error":"..."}` |

---

## 10. Logout

**Protocol:** `POST /logout`

**Auth:** `Authorization: Bearer *** — invalidates the current JWT session. After logout, the same token cannot be used to access protected endpoints.

### Input

| Header | Value |
|--------|-------|
| `Authorization` | `Bearer <jwt-token>` |

No body parameters required.

### Response - 200 OK

```json
{
  "message": "Sesion cerrada exitosamente"
}
```

### Error Responses

| Status | Condition | Body |
|--------|-----------|------|
| 401 Unauthorized | Missing/invalid/expired token | `{"error":"Credentials invalid"}` |

---

## Authentication Summary

```
POST /user-create        → Public (no token)
POST /login              → Public (no token) → returns JWT
GET  /health             → Public (no token)

GET  /user-info          → Requires Bearer token
POST /user-modify        → Requires Bearer token
GET  /public-proposals   → Requires Bearer token
GET  /saved_proposals    → Requires Bearer token
POST /saved-proposals    → Requires Bearer token
POST /proposal-save      → Requires Bearer token
POST /logout             → Requires Bearer token
```

### How to use the token

```bash
# 1. Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"MiClave"}'

# Response includes "token" — use it in subsequent requests

# 2. Call a protected endpoint
curl -X GET http://localhost:8080/public-proposals?limit=5 \
  -H "Authorization: Bearer <token>"
```

### Token expiry

JWT tokens expire **1 hour** after issuance. After expiry, the endpoint responds with:

```json
{"error":"Credentials invalid"}
```

Obtain a new token by calling `POST /login` again.

---

## Error Response Format

All errors follow the same structure:

```json
{
  "error": "description of the error"
}
```

### HTTP Status Codes

| Code | Meaning |
|------|---------|
| `200` | Success (GET/POST) |
| `201` | Created (POST) |
| `400` | Bad request (missing/invalid fields) |
| `401` | Unauthorized (missing/invalid token) |
| `404` | Not found |
| `409` | Conflict (duplicate document/email) |
| `500` | Internal server error |
