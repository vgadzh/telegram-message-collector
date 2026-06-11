# Roadmap v0.1 — Foundation

## Goal

Establish the technical foundation of the service and demonstrate successful communication with the Telegram API.

At the end of this milestone, the system should be deployable, observable, and capable of initiating Telegram account authorization by sending a verification code to a user phone number.

This version focuses on infrastructure, architecture, and operational readiness rather than business functionality.

---

## Scope

### Infrastructure

* PostgreSQL integration
* Database migration framework
* Configuration management
* Docker deployment
* Kubernetes manifests
* Local development environment

### Security

* Admin user bootstrap
* JWT authentication
* Protected API endpoints

### Observability

* Structured logging
* OpenTelemetry tracing
* Request correlation
* Health endpoints

### Telegram Integration

* Telegram client initialization
* Session storage abstraction
* Session creation workflow
* Verification code delivery

---

## Functional Requirements

### Authentication

The service must provide:

* Login endpoint
* JWT token generation
* JWT middleware
* Protected API routes

### Session Creation

The service must provide:

```
POST /api/v1/sessions
```

Request:
```json
{
  "phone": "+1234567890"
}
```

Expected behavior:

1. Validate request.
2. Create a new Telegram session record.
3. Initialize Telegram client.
4. Send verification code to the specified phone number.
5. Persist authorization state.
6. Return session identifier and current status.

### Health Monitoring

The service must expose:
```
GET /health/live
GET /health/ready
```

---

## Non-Functional Requirements

### Configuration

All configuration must be provided through environment variables.

Examples:

* database connection
* Telegram credentials
* JWT secrets
* tracing configuration

### Graceful Shutdown

The application must:

* react to termination signals
* stop accepting new requests
* release resources cleanly

### Logging

All requests must be logged with:

* timestamp
* request id
* execution duration
* response status

### Tracing

The application must generate traces for:

* HTTP requests
* database operations
* Telegram API calls

---

## Deliverables

The following artifacts must exist at the end of the milestone:

* PostgreSQL schema
* Migration scripts
* Docker configuration
* Kubernetes manifests
* Application configuration package
* Authentication module
* Telegram session module
* Session creation endpoint
* Tracing middleware
* Health endpoints

---

## Acceptance Criteria

The milestone is considered complete when:

* the service starts successfully
* migrations run automatically
* login endpoint returns a valid JWT
* protected endpoints require authentication
* a Telegram verification code can be sent to a valid phone number
* traces are generated
* health endpoints report correct status
* the application runs inside Docker
* the application can be deployed to Kubernetes

---

## Out of Scope

The following features are intentionally excluded from v0.1:

* verification code confirmation
* 2FA password handling
* runtime subscribers
* message processing
* whitelist filtering
* session ownership
* reconciliation loop
* update synchronization
* message persistence

These features will be implemented in later milestones.