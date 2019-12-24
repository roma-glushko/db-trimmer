package nonblocking

import (
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type TrimmerManager struct {
	ThreadCount      int
	connectionPool   *ConnectionPool
	trimmerWaitGroup *sync.WaitGroup
}

// NewTrimmerManager
func NewTrimmerManager(ThreadCount int, connectionPool *ConnectionPool) *TrimmerManager {
	var trimmerWaitGroup sync.WaitGroup

	return &TrimmerManager{
		ThreadCount,
		connectionPool,
		&trimmerWaitGroup,
	}
}

// Execute - init trimmers
func (m *TrimmerManager) Execute(trimChunkChannel chan TrimChunk) {
	log.Printf("Init Trimmers..")
	

	for i := 0; i < m.ThreadCount; i++ {
		trimmer := NewTrimmer(i+1, m.connectionPool, m.trimmerWaitGroup)

		go trimmer.Execute(trimChunkChannel)
	}
}

// WaitForTrimmers
func (m *TrimmerManager) WaitForTrimmers() {
	m.trimmerWaitGroup.Wait()
}
