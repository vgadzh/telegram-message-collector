# Ubiquitous Language

This document defines the common vocabulary used throughout the project.

All documentation, code, architecture discussions, and pull requests should use these terms consistently.

## User

A platform user authenticated through the service API.

A user may connect multiple Telegram accounts as long as each account uses a different phone number.

## Session

A Telegram account connected to the platform.

A session contains authorization data, encrypted Telegram session storage, and runtime configuration.

A session belongs to exactly one user.

## Authorization

The process of authenticating a Telegram account.

Authorization may require a verification code and, in some cases, a two-factor authentication password.

A successfully authorized session may later be activated for message processing.

## Subscriber

A runtime component responsible for receiving Telegram updates for a specific session.

Each active session has at most one subscriber running in the cluster.

## Runtime

The in-memory part of the system responsible for maintaining Telegram connections and processing updates.

Runtime state is considered temporary and may be reconstructed from persistent state stored in the database.

## Desired State

The state requested by the application and stored in the database.

Desired state represents how the system should behave.

Examples:

* active
* inactive
* deleted

## Actual State

The state currently observed in runtime.

Actual state may temporarily differ from desired state.

## Reconciliation

The process of continuously comparing desired state with actual state and applying corrective actions.

Reconciliation is responsible for starting, stopping, and recovering subscribers.

## Session Lease

A temporary ownership record that assigns a session to a specific service instance.

A lease prevents multiple nodes from processing the same session simultaneously.

## State Store

Persistent storage of Telegram synchronization state.

The state store contains update offsets required for recovery after restarts.

Examples:

* PTS
* QTS
* SEQ
* DATE

## Runtime Event

A domain event emitted by runtime components.

Runtime events are used to report lifecycle changes, failures, and operational information.

## Message

A Telegram message received by a subscriber.

Only messages that pass business filtering rules may be processed and stored.

## Message Processor

Business logic responsible for handling accepted messages.

The runtime layer delegates message processing to a message processor implementation.

## Whitelist

A collection of allowed Telegram senders.

Only messages originating from whitelisted senders may be processed.

The whitelist is loaded from the database and periodically refreshed in memory.

## Ownership

The relationship between a session and the service instance currently responsible for processing it.

Ownership is coordinated through session leases stored in PostgreSQL.

## Service Instance

A single running application process.

Multiple service instances may exist simultaneously in a Kubernetes deployment.