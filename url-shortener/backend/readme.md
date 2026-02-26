# URL Shortener Backend

## Overview

This project is a backend URL shortener written in **Go (v1.24.5)** using the standard `net/http` package. The goal is to understand how backend systems are built from first principles rather than relying heavily on frameworks.

The system follows a simple but practical pattern: fast in-memory storage for active data, and durable storage for long-term retention. It also includes authentication and request limiting to simulate real-world constraints.

---

## How It Works

A client sends a request to the API. Every request passes through middleware that validates a JWT and then enforces rate limits. If the request is allowed, the core handler processes it.

For URL resolution, the server first checks **Redis**. Since Redis is in-memory, lookups are fast and suitable for high-frequency reads. If the key does not exist in Redis (for example, after expiration), the system can retrieve it from **Amazon DynamoDB**, which serves as durable storage.

The response format returns shortened URLs under the `rilly/{uuid}` path.

---

## Rate Limiting and Tracking

Rate limiting is enforced server-side using Redis counters. Each validated JWT has an associated Redis key. On every request:

* The counter is incremented atomically.
* A 24-hour expiration is applied if the key is new.
* If the counter exceeds 10 requests, the request is rejected.

IP addresses are also tracked using a similar key structure. This prevents simple token farming and allows basic abuse monitoring.

No in-memory Go maps are used for tracking. This ensures the system remains safe under concurrency and scalable across multiple instances.

---

## Data Lifecycle

Active URLs live in Redis with a 24-hour sliding expiration. Each access refreshes the TTL, keeping frequently used links in fast storage.

When entries expire or become inactive, a background process batches relevant data to DynamoDB. DynamoDB holds durable records and supports indexing for longer-term storage and potential analytics.

This hot-to-cold storage model keeps the request path fast while preserving important data.

---

## Storage Design

Redis is responsible for:

* Fast URL resolution
* Hit counting
* Rate limiting
* Short-lived active data

DynamoDB is responsible for:

* Durable URL records
* Indexed lookups
* Long-term storage

Local development uses Docker for both services. Production would use managed DynamoDB on AWS.

---

## Design Principles

The system avoids heavy frameworks and keeps the request flow explicit and readable. Middleware handles authentication and rate limiting. Core handlers focus only on business logic. Storage responsibilities are clearly separated.

The emphasis is on understanding how caching, TTL policies, atomic counters, and batching strategies work in practice.
