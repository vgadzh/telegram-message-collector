# ADR-0004: Message Processing Uses At-Least-Once Delivery with Deduplication

## Status

Accepted

## Context

Telegram updates may be redelivered after:

* service restarts
* network interruptions
* runtime recovery
* update synchronization

The system processes financial notifications and must avoid storing duplicate messages.

Implementing true exactly-once delivery is not practical in a distributed system.

## Decision

The system uses at-least-once message delivery semantics.

Message persistence is protected by database-level deduplication.

Messages are uniquely identified by:

* sender_type
* sender_id
* telegram_message_id

A unique database constraint prevents duplicate message storage.

Message handlers must be idempotent.

## Consequences

### Positive

* Reliable recovery after failures
* Simple implementation
* Well-understood operational model

### Negative

* Message handlers must tolerate retries
* Duplicate delivery remains possible before persistence

### Alternatives Considered

* Exactly-once processing
* Distributed transaction coordination
* In-memory deduplication only