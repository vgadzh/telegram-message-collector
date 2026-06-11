# ADR-0001: PostgreSQL Is the Source of Truth

## Status

Accepted

## Context

The service manages Telegram sessions, runtime subscribers, authorization state, and update synchronization state.

The system is designed to run multiple service replicas. Runtime state may be lost during restarts, deployments, crashes, or node failures.

A consistent source of truth is required to recover the desired system state after failures.

## Decision

PostgreSQL is the single source of truth for all persistent state.

The database stores:

* Telegram sessions
* Authorization status
* Desired runtime state
* Session ownership information
* Telegram update synchronization state
* Processed messages

Runtime components must be reconstructed from database state during startup and reconciliation.

## Consequences

### Positive

* Predictable recovery after restarts
* Simpler operational model
* No dependency on external coordination systems
* Easier debugging and auditing

### Negative

* Increased dependency on PostgreSQL availability
* Additional database load from reconciliation and ownership management

### Alternatives Considered

* Runtime as source of truth
* Redis-based state management
* External coordination services