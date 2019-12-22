package nonblocking

import (
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type TrimChunk struct {
	StartInvervalID int
	EndIntervalID   int
}

type Trimmer struct {
	trimmerID      int
	connectionPool *ConnectionPool
	waitGroup      *sync.WaitGroup
}

// NewTrimmer
func NewTrimmer(trimmerID int, connectionPool *ConnectionPool, waitGroup *sync.WaitGroup) *Trimmer {

	return &Trimmer{
		trimmerID,
		connectionPool,
		waitGroup,
	}
}

// Execute - execute DB trim job
func (t *Trimmer) Execute(trimChunkChannel chan TrimChunk) {
	t.waitGroup.Add(1)
	defer t.waitGroup.Done()

	connection := t.connectionPool.Get()
	defer func() {
		t.connectionPool.Put(connection)
	}()

	for trimChunk := range trimChunkChannel {
		log.Printf("[Trimmer #%d] Trimming chunk: %d - %d", t.trimmerID, trimChunk.StartInvervalID, trimChunk.EndIntervalID)
		t.trimChunk(connection, trimChunk)
	}
}

// TrimChunk - trim planned data chunks
func (t *Trimmer) trimChunk(connection *Connection, trimChunk TrimChunk) {
	var err error

	if trimChunk.EndIntervalID == 0 {
		// last chunk trimming
		_, err = connection.driver.Exec("DELETE FROM catalog_product_entity WHERE entity_id >= ?", trimChunk.StartInvervalID)
	} else {
		_, err = connection.driver.Exec("DELETE FROM catalog_product_entity WHERE entity_id >= ? AND entity_id < ?", trimChunk.StartInvervalID, trimChunk.EndIntervalID)
	}

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}
