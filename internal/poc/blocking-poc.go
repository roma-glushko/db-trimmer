package poc

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// BlockingPoc
type BlockingPoc struct {
	dbDriverName     string
	dbDataSourceName string
	ChunkSize        int
}

// NewBlockingPoc
func NewBlockingPoc(dbDriverName string, dbDataSourceName string, chunkSize int) *BlockingPoc {
	return &BlockingPoc{
		dbDriverName,
		dbDataSourceName,
		chunkSize,
	}
}

// Execute - blocking trimming of a table by 1k chunks iterating though a numeric primary key including optimal chunk interval resolution (ID gaps prune)
func (p *BlockingPoc) Execute() {
	// Open up our database connection.
	db, err := sql.Open(p.dbDriverName, p.dbDataSourceName)

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	tableSize := p.getTableSize(db)

	log.Printf("Blocking PoC")
	log.Printf("Chunk Size: %d", p.ChunkSize)
	log.Printf("Table Size: %d", tableSize)

	executionStart := time.Now()

	startIntervalID := p.getStartIntervalID(db)
	endIntervalID := p.getEndIntervalID(db, startIntervalID, p.ChunkSize)

	for endIntervalID != 0 {
		p.deleteChunk(db, startIntervalID, endIntervalID)

		startIntervalID := endIntervalID
		endIntervalID = p.getEndIntervalID(db, startIntervalID, p.ChunkSize)
	}

	p.deleteLastChunk(db, startIntervalID)

	executionElapsed := time.Since(executionStart)
	log.Printf("Table Trimming took - %s", executionElapsed)
}

// GetTableSize
func (p *BlockingPoc) getTableSize(db *sql.DB) int {
	var count int

	err := db.QueryRow("SELECT COUNT(entity_id) FROM catalog_product_entity").Scan(&count)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return count
}

// GetStartIntervalID
func (p *BlockingPoc) getStartIntervalID(db *sql.DB) int {
	var startIntervalID int

	err := db.QueryRow("SELECT MIN(entity_id) FROM catalog_product_entity").Scan(&startIntervalID)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return startIntervalID
}

// GetEndIntervalID
func (p *BlockingPoc) getEndIntervalID(db *sql.DB, startIntervalID int, chunkSize int) int {
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

// DeleteChunk
func (p *BlockingPoc) deleteChunk(db *sql.DB, startIntervalID int, endIntervalID int) {
	_, err := db.Exec("DELETE FROM catalog_product_entity WHERE entity_id >= ? AND entity_id < ?", startIntervalID, endIntervalID)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

// DeleteLastChunk
func (p *BlockingPoc) deleteLastChunk(db *sql.DB, startIntervalID int) {
	_, err := db.Exec("DELETE FROM catalog_product_entity WHERE entity_id >= ?", startIntervalID)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}
