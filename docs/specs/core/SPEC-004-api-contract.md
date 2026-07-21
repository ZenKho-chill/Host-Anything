# SPEC-004: API Contract

## Status
Approved

## Overview
This document defines the REST API exposed by the Host Anything core engine. The Web UI and any potential CLI clients use this API to interact with the system. 

## Motivation
To provide a stable, well-documented interface for system interaction, ensuring loose coupling between the frontend UI and the Go backend core.

## Scope
- Base URL and headers.
- Endpoints for Service Management.
- Endpoints for Template/Marketplace Management.
- System endpoints (Health, Auth).
- Standard JSON response formats.

## Out of Scope
- Internal Go function signatures representing these routes.

## Specification

**Base URL**: `/api/v1`

### Authentication Header
All endpoints (except `/auth/login` and `/health`) require a Bearer token:
`Authorization: Bearer <JWT_TOKEN>`

### Data Schemas (JSON)

#### Standard Error Response
```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "Invalid configuration provided",
    "details": {
      "REDIS_PASSWORD": "Does not match required format"
    }
  }
}
```

### Endpoints

#### 1. `GET /health`
Returns system status.
**Response**: `200 OK`
```json
{
  "status": "up",
  "version": "1.0.0",
  "runtimes": ["docker"]
}
```

#### 2. `GET /services`
List all installed services.
**Response**: `200 OK`
```json
[
  {
    "id": "1234-5678",
    "name": "My Redis Cache",
    "template": "redis",
    "state": "RUNNING",
    "uptime_seconds": 3600
  }
]
```

#### 3. `POST /services`
Deploy a new service from a template.
**Request Body**:
```json
{
  "name": "My Redis Cache",
  "template_name": "redis",
  "template_version": "1.0.0",
  "config": {
    "REDIS_PASSWORD": "securepassword123",
    "MAX_MEMORY": "512mb"
  }
}
```
**Response**: `202 Accepted` - Returns the `id` of the newly created service.

#### 4. `GET /services/{id}`
Get detailed service status, including full config and runtime metrics.
**Response**: `200 OK`

#### 5. `PATCH /services/{id}/config`
Update service configuration.
**Request Body**: Key-value map of updated variables.
**Response**: `202 Accepted`

#### 6. `POST /services/{id}/start` & `POST /services/{id}/stop`
Initiate state transitions.
**Response**: `202 Accepted`

#### 7. `DELETE /services/{id}`
Remove the service and all its associated data/volumes.
**Response**: `204 No Content`

#### 8. `GET /templates`
List installed templates.
**Response**: `200 OK`

#### 9. `GET /templates/{id}`
Get full template schema (useful for UI rendering dynamic forms).
**Response**: `200 OK`

#### 10. `POST /marketplace/search`
Search the GitHub marketplace for templates.
**Request Body**: `{"query": "database", "tags": ["cache"]}`
**Response**: `200 OK` (Array of template metadata).

## Error Handling
- `400 Bad Request`: Validation errors.
- `401 Unauthorized`: Missing or invalid JWT.
- `403 Forbidden`: User lacks permissions.
- `404 Not Found`: Resource ID does not exist.
- `500 Internal Server Error`: Unhandled core exception.

## Security
- API is protected by JWT authentication (see SPEC-030).
- Input is strictly validated against schemas to prevent injection.

## Testing Strategy
- Integration tests using Go's `httptest` framework for all endpoints.
- Fuzz testing against POST/PATCH bodies to ensure the API does not panic on malformed JSON.
