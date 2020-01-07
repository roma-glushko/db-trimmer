package poc

import (
	"database/sql"
	nonblocking "db-trimmer/internal/non-blocking"
	"log"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// NonBlockingPoc
type NonBlockingPoc struct {
	dbDriverName       string
	dbDataSourceName   string
	ChunkSize          int
	PlannerThreadCount int
	TrimmerThreadCount int
}

// NewNonBlockingPoc
func NewNonBlockingPoc(dbDriverName string, dbDataSourceName string, chunkSize int, plannerThreadCount int, trimmerThreadCount int) *NonBlockingPoc {
	return &NonBlockingPoc{
		dbDriverName,
		dbDataSourceName,
		chunkSize,
		plannerThreadCount,
		trimmerThreadCount,
	}
}

// Execute - non-blocking trimming of a table by 1k chunks iterating though a numeric primary key including optimal chunk interval resolution (ID gaps prune)
func (p *NonBlockingPoc) Execute() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Printf("NonBlocking PoC")
	log.Printf("Num CPU: %d", runtime.NumCPU())
	log.Printf("Chunk Size: %d", p.ChunkSize)
	log.Printf("Planner Thread Count: %d", p.PlannerThreadCount)
	log.Printf("Trimmer Thread Count: %d", p.TrimmerThreadCount)

	executionStart := time.Now()

	var trimChunkChannel = make(chan nonblocking.TrimChunk, p.PlannerThreadCount*100)

	// init connection pool
	log.Printf("Init Connection Pool..")
	connectionPool, err := nonblocking.NewConnectionPool((p.PlannerThreadCount + p.TrimmerThreadCount + 1), p.dbDriverName, p.dbDataSourceName)

	if err != nil {
		log.Print(err.Error())
	}

	defer connectionPool.Close()

	// init trimmers
	trimmerManager := nonblocking.NewTrimmerManager(p.TrimmerThreadCount, connectionPool)
	trimmerManager.Execute(trimChunkChannel)

	// init planners
	plannerManager := nonblocking.NewPlannerManager(p.PlannerThreadCount, connectionPool, p.ChunkSize)
	plannerManager.Execute(trimChunkChannel)

	plannerManager.WaitForPlanners()
	close(trimChunkChannel)
	log.Printf("Chunk Planning has been finished")

	// wait till trimmers process planned chunks
	log.Printf("Waiting for trimmers to process planned chunks")
	trimmerManager.WaitForTrimmers()

	executionElapsed := time.Since(executionStart)
	log.Printf("Table Trimming took - %s", executionElapsed)
}

// GetTableSize
func (p *NonBlockingPoc) getTableSize(db *sql.DB) int {
	var count int

	err := db.QueryRow("SELECT COUNT(row_id) FROM catalog_product_entity").Scan(&count)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return count
}

// GetStartIntervalID
func (p *NonBlockingPoc) getStartIntervalID(db *sql.DB) int {
	var startIntervalID int

	err := db.QueryRow("SELECT MIN(row_id) FROM catalog_product_entity").Scan(&startIntervalID)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return startIntervalID
}

// GetEndIntervalID
func (p *NonBlockingPoc) getEndIntervalID(db *sql.DB, startIntervalID int, chunkSize int) int {
	var endIntervalID int

	err := db.QueryRow("SELECT row_id FROM catalog_product_entity WHERE row_id >= ? ORDER BY row_id LIMIT ?,1", startIntervalID, chunkSize).Scan(&endIntervalID)

	if err == sql.ErrNoRows {
		return 0
	}

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return endIntervalID
}
