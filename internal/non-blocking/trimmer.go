package nonblocking

import (
	"log"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
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
		executionStart := time.Now()
		t.trimChunk(connection, trimChunk)
		executionElapsed := time.Since(executionStart)
		log.Printf("[Trimmer #%d] Trimming chunk: %d - %d (%s)", t.trimmerID, trimChunk.StartInvervalID, trimChunk.EndIntervalID, executionElapsed)
	}
}

// trimChunk
func (t *Trimmer) trimChunk(connection *Connection, trimChunk TrimChunk) {
	// retry deleting chunks on locks
	for try := 0; try < 5; try++ {
		err := t.deleteChunk(connection, trimChunk)

		if err == nil {
			return
		}

		if err.Number != 1205 && err.Number != 1213 {
			// unknown error
			panic(err.Error())
		}

		log.Printf("[Trimmer #%d] Trimming chunk: %d : %d - Retry", t.trimmerID, trimChunk.StartInvervalID, trimChunk.EndIntervalID)
	}
}

// deleteChunk - trim planned data chunks
func (t *Trimmer) deleteChunk(connection *Connection, trimChunk TrimChunk) *mysql.MySQLError {
	var err error

	if trimChunk.EndIntervalID == 0 {
		// last chunk trimming
		_, err = connection.driver.Exec("DELETE FROM catalog_product_entity WHERE row_id >= ?", trimChunk.StartInvervalID)
	} else {
		_, err = connection.driver.Exec("DELETE FROM catalog_product_entity WHERE row_id >= ? AND row_id < ?", trimChunk.StartInvervalID, trimChunk.EndIntervalID)
	}

	if err == nil {
		return nil
	}

	return err.(*mysql.MySQLError)
}
