# ADR-0003: Session Ownership Is Coordinated Through PostgreSQL Leases

## Status

Accepted

## Context

The service is expected to run multiple replicas.

A Telegram session must be processed by exactly one runtime instance at a time.

Running the same session on multiple nodes may lead to:

* duplicate processing
* session corruption
* update state conflicts

A distributed ownership mechanism is required.

## Decision

Session ownership is managed through lease records stored in PostgreSQL.

Each active session is assigned to a single node.

Ownership includes:

* owner node identifier
* lease acquisition timestamp
* lease expiration timestamp

Nodes periodically renew leases while processing sessions.

Expired leases may be acquired by another node.

## Consequences

### Positive

* No additional infrastructure required
* Simple operational model
* Automatic recovery from node failures

### Negative

* Additional database activity
* Ownership acquisition logic must be implemented carefully

### Alternatives Considered

* Redis locks
* Consul sessions
* Kubernetes leader election