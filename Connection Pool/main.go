package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
)

const connString = "postgres://utkarshshrivastava:Qwerty121!@localhost:5434/postgres?sslmode=disable"

type ConnPool struct {
	conns chan *pgx.Conn
}

func NewConnPool(ctx context.Context, size int) (*ConnPool, error) {
	ch := make(chan *pgx.Conn, size)

	for i := 0; i < size; i++ {
		conn, err := pgx.Connect(ctx, connString)
		if err != nil {
			// cleanup any already-open connections
			close(ch)
			for c := range ch {
				c.Close(ctx)
			}
			return nil, fmt.Errorf("failed to create connection %d: %w", i, err)
		}
		ch <- conn
	}

	return &ConnPool{conns: ch}, nil
}

func (p *ConnPool) Acquire() *pgx.Conn {
	return <-p.conns // blocks if none available
}

func (p *ConnPool) Release(conn *pgx.Conn) {
	p.conns <- conn
}

func (p *ConnPool) Close(ctx context.Context) {
	close(p.conns)
	for conn := range p.conns {
		conn.Close(ctx)
	}
}

func worker(ctx context.Context, wg *sync.WaitGroup, id int, pool *ConnPool) {
	defer wg.Done()

	// get a connection from the pool
	conn := pool.Acquire()
	defer pool.Release(conn)

	// simulate a slow query
	_, err := conn.Exec(ctx, "SELECT pg_sleep(5);")
	if err != nil {
		log.Printf("goroutine %d error: %v", id, err)
		return
	}
	fmt.Printf("Goroutine %d done\n", id)
}

func main() {
	ctx := context.Background()

	// create a pool of, say, 10 connections
	pool, err := NewConnPool(ctx, 10)
	if err != nil {
		log.Fatal("failed to create pool:", err)
	}
	defer pool.Close(ctx)

	var wg sync.WaitGroup

	// 100 goroutines, but at most 10 DB connections in use at a time
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go worker(ctx, &wg, i, pool)
	}

	wg.Wait()
	fmt.Println("All done")
}
