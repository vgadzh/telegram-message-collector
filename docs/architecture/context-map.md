# Context Map

This document describes relationships between bounded contexts.

## Overview

```
Identity Context
    |
    v
Session Management Context
    |
    v
Runtime Context
    |
    v
Message Processing Context
Infrastructure Context
    ^
    |
supports all contexts
```

## Identity Context → Session Management Context

The Session Management Context depends on the Identity Context.

A user must be authenticated before performing session management operations.

Examples:

* create session
* activate session
* delete session

The Identity Context does not depend on Session Management.

## Session Management Context → Runtime Context

The Runtime Context depends on session information stored by the Session Management Context.

The Session Management Context defines the desired state of a session.

The Runtime Context is responsible for achieving that state through reconciliation.

Examples:

* activate session
* deactivate session
* recover session after restart

## Runtime Context → Message Processing Context

The Runtime Context receives Telegram updates and delegates accepted messages to the Message Processing Context.

The Runtime Context does not contain business-specific message handling logic.

Message Processing is implemented as a separate concern.

## Infrastructure Context

The Infrastructure Context supports all other contexts.

Examples:

* PostgreSQL repositories
* configuration providers
* tracing
* metrics
* logging

Business contexts depend on infrastructure abstractions rather than concrete implementations.

## Dependency Direction

The preferred dependency flow is:

```
Identity
    ↓
Session Management
    ↓
Runtime
    ↓
Message Processing
Infrastructure
    ↑
supports all layers
```

Dependencies should always point toward business capabilities and never in the opposite direction.

The Runtime Context must not depend on HTTP handlers.

The Message Processing Context must not depend on Telegram-specific implementation details.

The Infrastructure Context must not contain business rules.