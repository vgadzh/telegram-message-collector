# Bounded Contexts

This document describes the primary bounded contexts of the system.

The goal is to separate business responsibilities and establish clear ownership of domain concepts.

## Identity Context

The Identity Context is responsible for authentication and authorization within the platform.

Responsibilities:

* user authentication
* JWT generation
* JWT validation
* access control

Core concepts:

* User
* Access Token
* Authentication

This context does not contain any Telegram-specific logic.

## Session Management Context

The Session Management Context is responsible for managing Telegram account connections.

Responsibilities:

* session creation
* authorization workflow
* session status management
* encrypted session storage

Core concepts:

* Session
* Authorization
* Session Status

This context owns the lifecycle of Telegram accounts connected to the platform.

## Runtime Context

The Runtime Context is responsible for executing active Telegram sessions.

Responsibilities:

* subscriber lifecycle management
* Telegram connectivity
* reconciliation
* ownership management
* state synchronization

Core concepts:

* Subscriber
* Desired State
* Actual State
* Session Lease
* Reconciliation
* State Store

This context contains all long-running background processes.

## Message Processing Context

The Message Processing Context is responsible for handling incoming messages.

Responsibilities:

* message filtering
* whitelist validation
* message processing
* message persistence

Core concepts:

* Message
* Message Processor
* Whitelist

This context defines business rules for incoming Telegram messages.

## Infrastructure Context

The Infrastructure Context provides technical capabilities required by other contexts.

Responsibilities:

* PostgreSQL integration
* migrations
* observability
* tracing
* configuration
* deployment

Core concepts:

* Repository
* Configuration
* Trace
* Metric
* Log

This context should not contain business logic.