package nonblocking

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type Planner struct {
	plannerID      int
	connectionPool *ConnectionPool
	waitGroup      *sync.WaitGroup
	chunkSize      int
}

// NewPlanner
func NewPlanner(plannerID int, connectionPool *ConnectionPool, waitGroup *sync.WaitGroup, chunkSize int) *Planner {
	return &Planner{
		plannerID,
		connectionPool,
		waitGroup,
		chunkSize,
	}
}

// Execute - execute data chunk plan job
func (p *Planner) Execute(trimChunkChannel chan TrimChunk, plannerLeftLimitID int, plannerRightLimitID int) {
	p.waitGroup.Add(1)
	defer p.waitGroup.Done()

	connection := p.connectionPool.Get()
	defer func() {
		p.connectionPool.Put(connection)
	}()

	startIntervalID := plannerLeftLimitID
	endIntervalID := p.getEndIntervalID(connection.driver, startIntervalID, p.chunkSize, plannerRightLimitID)

	for endIntervalID != 0 {
		// chunk planning
		log.Printf("[Planner #%d] Planning chunk: %d - %d", p.plannerID, startIntervalID, endIntervalID)
		trimChunkChannel <- TrimChunk{
			startIntervalID,
			endIntervalID,
		}

		startIntervalID = endIntervalID + 1
		endIntervalID = p.getEndIntervalID(connection.driver, startIntervalID, p.chunkSize, plannerRightLimitID)
	}

	log.Printf("[Planner #%d] Planning last chunk: %d - %d", p.plannerID, startIntervalID, endIntervalID)
	trimChunkChannel <- TrimChunk{
		startIntervalID,
		plannerRightLimitID,
	}
}

// GetEndIntervalID
func (p *Planner) getEndIntervalID(db *sql.DB, startIntervalID int, chunkSize int, rightLimitID int) int {
	var endIntervalID int

	err := db.QueryRow("SELECT entity_id FROM catalog_product_entity WHERE entity_id >= ? AND entity_id <= ? ORDER BY entity_id LIMIT ?,1", startIntervalID, rightLimitID, chunkSize).Scan(&endIntervalID)

	if err == sql.ErrNoRows {
		return 0
	}

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return endIntervalID
}
