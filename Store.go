package customstore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
)

// Store defines a session store
type storeImplementation struct {
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
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	store := &storeImplementation{
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
		store.dbDriverName = sb.DatabaseDriverName(store.db)
	}

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}

// AutoMigrate migrates the tables
func (st *storeImplementation) AutoMigrate() error {
	sql := st.SqlCreateTable()

	if st.debug {
		log.Println(sql)
	}

	_, err := st.db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debug = debug
}

// RecordCreate creates a record
func (st *storeImplementation) RecordCreate(record *Record) error {
	if record.ID == "" {
		record.ID = uid.HumanUid()
	}
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	var sqlStr string
	sqlStr, sqlParams, errSQL := goqu.Dialect(st.dbDriverName).
		Insert(st.tableName).
		Rows(record).
		Prepared(true).
		ToSQL()

	if errSQL != nil {
		return errSQL
	}

	if st.debug {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	return nil
}

// RecordFindByID finds a record by ID
func (st *storeImplementation) RecordFindByID(id string) (*Record, error) {
	sqlStr, sqlParams, _ := goqu.Dialect(st.dbDriverName).
		From(st.tableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id), goqu.C(COLUMN_DELETED_AT).IsNull()).
		Limit(1).
		Select().
		ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	var record Record
	err := sqlscan.Get(context.Background(), st.db, &record, sqlStr, sqlParams...)

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

// RecordUpdate updates a record
func (st *storeImplementation) RecordUpdate(record *Record) error {
	fields := map[string]interface{}{}
	fields[COLUMN_RECORD_DATA] = record.Data
	fields[COLUMN_UPDATED_AT] = time.Now()

	var sqlStr string

	sqlStr, sqlParams, errSQL := goqu.Dialect(st.dbDriverName).
		Update(st.tableName).
		Set(fields).
		Where(goqu.C("id").Eq(record.ID)).
		Prepared(true).
		ToSQL()

	if errSQL != nil {
		return errSQL
	}

	if st.debug {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	return nil
}
