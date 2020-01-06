/*
 * go-mydumper
 * xelabs.org
 *
 * Copyright (c) XeLabs
 * GPL License
 *
 */

package nonblocking

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectionPool
type ConnectionPool struct {
	mu    sync.RWMutex
	conns chan *Connection
}

// Connection
type Connection struct {
	ID     int
	driver *sql.DB
}

// NewConnectionPool creates the new pool.
func NewConnectionPool(threadCount int, driverName string, dataSourceName string) (*ConnectionPool, error) {
	conns := make(chan *Connection, threadCount)

	for i := 0; i < threadCount; i++ {
		db, err := sql.Open(driverName, dataSourceName)

		if err != nil {
			return nil, err
		}

		// configure connection session

		db.Exec("SET @@session.triggers = OFF")

		conn := &Connection{ID: i, driver: db}

		conns <- conn
	}

	return &ConnectionPool{
		conns: conns,
	}, nil
}

// Get used to get one connection from the pool.
func (p *ConnectionPool) Get() *Connection {
	conns := p.getConnections()

	if conns == nil {
		return nil
	}

	conn := <-conns

	return conn
}

// Put used to put one connection to the pool.
func (p *ConnectionPool) Put(conn *Connection) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.conns == nil {
		return
	}

	p.conns <- conn
}

// Close used to close the pool and the connections.
func (p *ConnectionPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.conns)

	for conn := range p.conns {
		conn.driver.Close()
	}

	p.conns = nil
}

func (p *ConnectionPool) getConnections() chan *Connection {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.conns
}
