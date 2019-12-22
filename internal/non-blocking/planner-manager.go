package nonblocking

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type PlannerManager struct {
	threadCount      int
	chunkSize        int
	connectionPool   *ConnectionPool
	plannerWaitGroup *sync.WaitGroup
}

// NewPlannerManager
func NewPlannerManager(threadCount int, connectionPool *ConnectionPool, chunkSize int) *PlannerManager {
	var plannerWaitGroup sync.WaitGroup

	return &PlannerManager{
		threadCount,
		chunkSize,
		connectionPool,
		&plannerWaitGroup,
	}
}

// Execute - init trimmers
func (m *PlannerManager) Execute(trimChunkChannel chan TrimChunk) {
	log.Printf("Init Planners..")

	connection := m.connectionPool.Get()

	tableSize := m.getTableSize(connection.driver)
	log.Printf("Table Size: %d", tableSize)

	intervalSize := tableSize / m.threadCount
	log.Printf("Planner Inverval Size: %d", intervalSize)

	// planning first planner interval
	startPlannerIntervalID := m.getPlannerStartIntervalID(connection.driver)
	endPlannerIntervalID := m.getPlannerEndIntervalID(connection.driver, startPlannerIntervalID, intervalSize)

	for i := 0; i < m.threadCount-1; i++ {
		planner := NewPlanner(i+1, m.connectionPool, m.plannerWaitGroup, m.chunkSize)

		go planner.Execute(trimChunkChannel, startPlannerIntervalID, endPlannerIntervalID)

		log.Printf("Init Planner #%d - %d : %d", i+1, startPlannerIntervalID, endPlannerIntervalID)

		// planning next planner interval
		startPlannerIntervalID = endPlannerIntervalID + 1
		endPlannerIntervalID = m.getPlannerEndIntervalID(connection.driver, startPlannerIntervalID, intervalSize)
	}

	// planning last planner interval
	endPlannerIntervalID = m.getLastIntervalID(connection.driver)

	planner := NewPlanner(m.threadCount, m.connectionPool, m.plannerWaitGroup, m.chunkSize)
	go planner.Execute(trimChunkChannel, startPlannerIntervalID, endPlannerIntervalID)

	log.Printf("Init Planner #%d - %d : %d", m.threadCount, startPlannerIntervalID, endPlannerIntervalID)

	m.connectionPool.Put(connection)
}

// WaitForPlanners
func (m *PlannerManager) WaitForPlanners() {
	m.plannerWaitGroup.Wait()
}

// getPlannerStartIntervalID
func (m *PlannerManager) getPlannerStartIntervalID(db *sql.DB) int {
	var startIntervalID int

	err := db.QueryRow("SELECT MIN(entity_id) FROM catalog_product_entity").Scan(&startIntervalID)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return startIntervalID
}

// getLastIntervalID
func (m *PlannerManager) getLastIntervalID(db *sql.DB) int {
	var startIntervalID int

	err := db.QueryRow("SELECT MAX(entity_id) FROM catalog_product_entity").Scan(&startIntervalID)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return startIntervalID
}

// GetEndIntervalID
func (m *PlannerManager) getPlannerEndIntervalID(db *sql.DB, startIntervalID int, chunkSize int) int {
	var endIntervalID int

	err := db.QueryRow("SELECT entity_id FROM catalog_product_entity WHERE entity_id >= ? ORDER BY entity_id LIMIT ?,1", startIntervalID, chunkSize).Scan(&endIntervalID)

	if err == sql.ErrNoRows {
		return 0
	}

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return endIntervalID
}

// GetTableSize
func (m *PlannerManager) getTableSize(db *sql.DB) int {
	var count int

	err := db.QueryRow("SELECT COUNT(entity_id) FROM catalog_product_entity").Scan(&count)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return count
}
