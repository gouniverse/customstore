package customstore

import (
	"database/sql"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gouniverse/uid"
)

// Store defines a session store
type Store struct {
	tableName          string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debug              bool
}

// StoreOption options for the vault store
type StoreOption func(*Store)

// NewStore creates a new entity store
func NewStore(opts ...StoreOption) (*Store, error) {
	store := &Store{}
	for _, opt := range opts {
		opt(store)
	}

	if store.tableName == "" {
		log.Panic("Custom store: tableName is required")
	}

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}

// DriverName finds the driver name from database
func (st *Store) DriverName(db *sql.DB) string {
	dv := reflect.ValueOf(db.Driver())
	driverFullName := dv.Type().String()
	if strings.Contains(driverFullName, "mysql") {
		return "mysql"
	}
	if strings.Contains(driverFullName, "postgres") || strings.Contains(driverFullName, "pq") {
		return "postgres"
	}
	if strings.Contains(driverFullName, "sqlite") {
		return "sqlite"
	}
	if strings.Contains(driverFullName, "mssql") {
		return "mssql"
	}
	return driverFullName
}

// AutoMigrate migrates the tables
func (st *Store) AutoMigrate() error {
	sql := st.SqlCreateTable()

	if st.debug {
		log.Println(sql)
	}

	_, err := st.db.Exec(sql)
	if err != nil {
		// log.Println(err)
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) {
	st.debug = debug
}

// WithAutoMigrate sets the table name for the cache store
func WithAutoMigrate(automigrateEnabled bool) StoreOption {
	return func(s *Store) {
		s.automigrateEnabled = automigrateEnabled
	}
}

// WithDb sets the database for the setting store
func WithDb(db *sql.DB) StoreOption {
	return func(s *Store) {
		s.db = db
		s.dbDriverName = s.DriverName(s.db)
	}
}

// WithDebug prints the SQL queries
func WithDebug(debug bool) StoreOption {
	return func(s *Store) {
		s.debug = debug
	}
}

// WithTableName sets the table name for the custom record
func WithTableName(tableName string) StoreOption {
	return func(s *Store) {
		s.tableName = tableName
	}
}

// RecordCreate creates a node
func (st *Store) RecordCreate(record *Record) (bool, error) {
	if record.ID == "" {
		record.ID = uid.HumanUid()
	}
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	// log.Println(record)
	// log.Println(st.debug)
	// log.Println(st.tableName)

	var sqlStr string
	sqlStr, _, errSQL := goqu.Dialect(st.dbDriverName).Insert(st.tableName).Rows(record).ToSQL()

	if errSQL != nil {
		return false, errSQL
	}

	// log.Println(errSQL)

	if st.debug {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		return false, err
	}

	return true, nil
}
