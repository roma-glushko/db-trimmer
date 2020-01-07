package poc

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Table
type Table struct {
	Name        string
	TypeGroup   string
	SQL         string
	Fields      []string
	Constraints []Constraint
	Triggers    []Trigger
}

// Constraint
type Constraint struct {
	Name                 string
	ColumnName           string
	ReferencedTableName  string
	ReferencedColumnName string
}

// Trigger
type Trigger struct {
	Name  string
	Event string
	SQL   string
}

// CopySchema
type CopySchemaPoc struct {
	driverName string
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string
	dbName     string
}

// NewCopySchemaPoc
func NewCopySchemaPoc(
	dbDriverName string,
	dbUser string,
	dbPassword string,
	dbHost string,
	dbPort string,
	dbName string,
) *CopySchemaPoc {
	return &CopySchemaPoc{
		dbDriverName,
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	}
}

// Execute
func (p *CopySchemaPoc) Execute() {
	log.Printf("Copy Database Schema..")
	executionStart := time.Now()

	db, err := sql.Open(
		p.driverName,
		fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s`, p.dbUser, p.dbPassword, p.dbHost, p.dbPort, p.dbName),
	)

	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	var mirrorDbName = "db-trim-mirror"

	mirrorDb, err := sql.Open(
		p.driverName,
		fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s`, p.dbUser, p.dbPassword, p.dbHost, p.dbPort, mirrorDbName),
	)

	if err != nil {
		log.Print(err.Error())
	}

	// configure session
	mirrorDb.Exec(`SET FOREIGN_KEY_CHECKS=0`)

	tableList := p.getTableList(db)

	// apply table schemas
	for _, table := range tableList {
		_, err := mirrorDb.Exec(table.SQL)

		if err != nil {
			log.Print(err.Error())
		}

		// apply trigger schemas
		for _, trigger := range table.Triggers {
			_, err := mirrorDb.Exec(trigger.SQL)

			if err != nil {
				log.Print(err.Error())
			}
		}
	}

	// check INSERT FROM SELECT query

	var fieldList = ``
	var catalogProductEntityTable = `catalog_product_entity`

	mirrorDb.Exec(`INSERT INTO ` + mirrorDbName + `.` + catalogProductEntityTable + ` (` + fieldList + `) 
	SELECT ` + fieldList + ` FROM ` + p.dbName + `.` + catalogProductEntityTable + ` `)

	executionElapsed := time.Since(executionStart)
	log.Printf("DB Copy took - %s", executionElapsed)
}

// getTableList
func (p *CopySchemaPoc) getTableList(db *sql.DB) []Table {
	tableRows, err := db.Query(`SHOW FULL TABLES`)

	if err != nil {
		panic(err.Error())
	}

	var tableList = make([]Table, 0)

	for tableRows.Next() {
		var nextTable = Table{}
		err := tableRows.Scan(&nextTable.Name, &nextTable.TypeGroup)

		if err != nil {
			panic(err.Error())
		}

		// get SQL difinition
		nextTable.SQL = p.getTableSQL(db, nextTable.Name)

		// get field list
		nextTable.Fields = p.getTableFields(db, nextTable.Name)

		// get constraints
		nextTable.Constraints = p.getConstraintList(db, nextTable.Name, p.dbName)

		// get triggers
		nextTable.Triggers = p.getTriggerList(db, nextTable.Name)

		tableList = append(tableList, nextTable)
	}

	return tableList
}

// getTableSQL
func (p *CopySchemaPoc) getTableSQL(db *sql.DB, tableName string) string {
	var dummyTableName string
	var tableCreateSQL string

	err := db.QueryRow(`SHOW CREATE TABLE `+tableName).Scan(&dummyTableName, &tableCreateSQL)

	if err != nil {
		panic(err.Error())
	}

	return tableCreateSQL
}

// getConstraintList
func (p *CopySchemaPoc) getConstraintList(db *sql.DB, tableName string, databaseName string) []Constraint {
	constraintRows, err := db.Query(
		`SELECT CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME 
		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
		WHERE TABLE_NAME = '` + tableName + `' AND REFERENCED_TABLE_SCHEMA = '` + databaseName + `'`,
	)

	if err != nil {
		log.Fatal(err)
	}

	constraints := make([]Constraint, 0)

	for constraintRows.Next() {
		var constraint = Constraint{}

		err := constraintRows.Scan(
			&constraint.Name,
			&constraint.ColumnName,
			&constraint.ReferencedTableName,
			&constraint.ReferencedColumnName,
		)

		if err != nil {
			log.Fatal(err)
		}

		constraints = append(constraints, constraint)
	}

	return constraints
}

// getTriggerList
func (p *CopySchemaPoc) getTriggerList(db *sql.DB, tableName string) []Trigger {
	triggerRows, err := db.Query(`SHOW TRIGGERS LIKE '` + tableName + `'`)

	if err != nil {
		panic(err.Error())
	}

	var dummy string
	var triggerList = make([]Trigger, 0)

	for triggerRows.Next() {
		var nextTrigger = Trigger{}

		err := triggerRows.Scan(
			&nextTrigger.Name,
			&nextTrigger.Event,
			&dummy,
			&dummy,
			&dummy,
			&dummy,
			&dummy,
			&dummy,
			&dummy,
			&dummy,
			&dummy,
		)

		if err != nil {
			panic(err.Error())
		}

		nextTrigger.SQL = p.getTriggerSQL(db, nextTrigger.Name)

		triggerList = append(triggerList, nextTrigger)
	}

	return triggerList
}

// getTriggerSQL
func (p *CopySchemaPoc) getTriggerSQL(db *sql.DB, triggerName string) string {
	var dummy string
	var triggerCreateSQL string

	err := db.QueryRow(`SHOW CREATE TRIGGER `+triggerName).Scan(
		&dummy,
		&dummy,
		&triggerCreateSQL,
		&dummy,
		&dummy,
		&dummy,
		&dummy,
	)

	if err != nil {
		panic(err.Error())
	}

	return triggerCreateSQL
}

// getTableFields
func (p *CopySchemaPoc) getTableFields(db *sql.DB, tableName string) []string {
	fieldRows, err := db.Query(`SHOW COLUMNS FROM ` + tableName)

	if err != nil {
		panic(err.Error())
	}

	var fieldList = make([]string, 0)

	for fieldRows.Next() {
		var dummyStr sql.NullString
		var tableField string
		err := fieldRows.Scan(&tableField, &dummyStr, &dummyStr, &dummyStr, &dummyStr, &dummyStr)

		if err != nil {
			panic(err.Error())
		}

		fieldList = append(fieldList, tableField)
	}

	return fieldList
}
