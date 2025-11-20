# Go Connection Pool

A simple, thread-safe PostgreSQL connection pool implementation in Go from scratch. This project demonstrates how to manage database connections efficiently using Go channels to create a blocking queue.

## Overview

We've built a custom connection pool that:

- Manages a fixed set of database connections
- Uses a **buffered channel** as a blocking queue
- Handles concurrent access safely
- Blocks callers when no connections are available until one is released

## Components

1. **`ConnPool` Struct**: Holds the buffered channel of `*pgx.Conn`.
2. **`NewConnPool(size)`**: Initializes the pool by opening `size` connections.
3. **`Acquire()`**: Retreives a connection from the channel. Blocks if the pool is empty.
4. **`Release(conn)`**: Returns a connection back to the channel for reuse.
5. **Simulation**: A test in `main.go` that spawns 100 goroutines sharing a pool of only 10 connections.

## Prerequisites

- Go 1.21+
- PostgreSQL running locally

## Configuration

The connection string is currently hardcoded in `main.go`:

```go
const connString = "postgres://username:password@localhost:5434/postgres?sslmode=disable"
```

Update this string in `main.go` to match your local PostgreSQL configuration (host, port, user, password).

## How to Run

1. **Initialize the module** (if not already done):

   ```bash
   go mod tidy
   ```

2. **Run the program**:

   ```bash
   go run main.go
   ```

## Expected Output

You will see 100 goroutines trying to execute a query that sleeps for 5 seconds. Since the pool size is 10, they will execute in batches of 10.

```text
Goroutine 5 done
Goroutine 1 done
...
Goroutine 99 done
All done
```

The total execution time will be roughly `(100 queries / 10 concurrent) * 5 seconds = 50 seconds` (approx), proving the pool limits concurrency correctly.
