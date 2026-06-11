# ADR-0002: Runtime State Is Managed Through Reconciliation

## Status

Accepted

## Context

The service exposes HTTP endpoints that allow users to create, activate, deactivate, and delete Telegram sessions.

Directly controlling runtime components from HTTP handlers couples API requests to process lifecycle management and makes recovery difficult.

A mechanism is required to restore runtime state after restarts and node failures.

## Decision

HTTP requests modify only the desired state stored in PostgreSQL.

A background reconciliation process continuously compares:

* desired state stored in the database
* actual state observed in runtime

The reconciler starts, stops, or recovers runtime components as necessary.

## Consequences

### Positive

* Automatic recovery after restart
* Idempotent management operations
* Simplified runtime lifecycle
* Better support for distributed deployment

### Negative

* State changes are eventually consistent
* Additional implementation complexity

### Alternatives Considered

* Direct runtime control from HTTP handlers
* Event-driven runtime management