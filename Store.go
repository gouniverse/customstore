package customstore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/sqlscan"
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

// NewStoreOptions define the options for creating a new session store
type NewStoreOptions struct {
	TableName          string
	DB                 *sql.DB
	DbDriverName       string
	TimeoutSeconds     int64
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new session store
func NewStore(opts NewStoreOptions) (*Store, error) {
	store := &Store{
		tableName:          opts.TableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debug:              opts.DebugEnabled,
	}

	if store.tableName == "" {
		return nil, errors.New("customstore store: tableName is required")
	}

	if store.db == nil {
		return nil, errors.New("session store: DB is required")
	}

	if store.dbDriverName == "" {
		store.dbDriverName = store.DriverName(store.db)
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

// RecordCreate creates a node
func (st *Store) RecordCreate(record *Record) (bool, error) {
	if record.ID == "" {
		record.ID = uid.HumanUid()
	}
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	var sqlStr string
	sqlStr, _, errSQL := goqu.Dialect(st.dbDriverName).Insert(st.tableName).Rows(record).ToSQL()

	if errSQL != nil {
		return false, errSQL
	}

	if st.debug {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		return false, err
	}

	return true, nil
}

// RecordFindByID finds a user by ID
func (st *Store) RecordFindByID(id string) (*Record, error) {
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).
		From(st.tableName).
		Where(goqu.C("id").Eq(id), goqu.C("deleted_at").IsNull()).
		Limit(1).
		Select().
		ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	var record Record
	err := sqlscan.Get(context.Background(), st.db, &record, sqlStr)

	if err != nil {
		if err == sql.ErrNoRows {
			// Looks like this is now outdated for sqlscan
			return nil, nil
		}
		if sqlscan.NotFound(err) {
			return nil, nil
		}
		log.Println("Failed to execute query: ", err)
		return nil, err
	}

	return &record, nil
}
