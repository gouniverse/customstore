package customstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/sqlscan"
)

// Record type
type Record struct {
	ID        string     `json:"id" db:"id"`                   // varchar(40) primary_key
	Type      string     `json:"record_type" db:"record_type"` // varchar(100) DEFAULT NULL
	Data      string     `json:"record_data" db:"record_data"` // longtext DEFAULT NULL
	CreatedAt time.Time  `json:"created_at" db:"created_at"`   // datetime NOT NULL
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`   // datetime NOT NULL
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`   // datetime DEFAULT NULL
}

func (r *Record) SetMap(metas map[string]interface{}) {
	jsonBytes, err := json.Marshal(metas)
	if err != nil {
		log.Panic(err.Error())
	}
	jsonString := string(jsonBytes)
	r.Data = jsonString
}

func (r *Record) GetMap() map[string]interface{} {
	var data map[string]interface{}

	if r.Data == "" {
		return data
	}

	err := json.Unmarshal([]byte(r.Data), &data)

	if err != nil {
		log.Panic(err.Error())
	}

	return data
}

// RecordFindByID finds a user by ID
func (st *Store) RecordFindByID(id string) (*Record, error) {
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).From(st.tableName).Where(goqu.C("id").Eq(id), goqu.C("deleted_at").IsNull()).Select().Limit(1).ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	var record Record
	err := sqlscan.Get(context.Background(), st.db, &record, sqlStr)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, nil
		}
		log.Println("Failed to execute query: ", err)
		return nil, err
	}

	return &record, nil
}

// UserUpdate creates a user
// func (st *Store) UserUpdate(user *User) bool {

// 	// result := st.db.Table(st.userTableName).Save(&user)

// 	// if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 	// 	return false
// 	// }

// 	return true
// }

// SqlCreateUserTable returns a SQL string for creating the user table
func (st *Store) SqlCreateTable() string {
	sqlMysql := `
	CREATE TABLE IF NOT EXISTS ` + st.tableName + ` (
	  id varchar(40) NOT NULL PRIMARY KEY,
	  record_type varchar(100) NOT NULL, 
	  record_data longtext DEFAULT NULL,
	  created_at datetime NOT NULL,
	  updated_at datetime NOT NULL,
	  deleted_at datetime DEFAULT NULL
	);
	`

	sqlPostgres := `
	CREATE TABLE IF NOT EXISTS "` + st.tableName + `" (
	  "id" varchar(40) NOT NULL PRIMARY KEY,
	  "record_type" varchar(100) NOT NULL,
	  "record_data" longtext NOT NULL,
	  "created_at" timestamptz(6) NOT NULL,
	  "updated_at" timestamptz(6) NOT NULL,
	  "deleted_at" timestamptz(6) DEFAULT NULL
	)
	`

	sqlSqlite := `
	CREATE TABLE IF NOT EXISTS "` + st.tableName + `" (
	  "id" varchar(40) NOT NULL PRIMARY KEY,
	  "record_type" varchar(100) NOT NULL,
	  "record_data" longtext DEFAULT NULL,
	  "created_at" datetime NOT NULL,
	  "updated_at" datetime NOT NULL,
	  "deleted_at" datetime DEFAULT NULL
	)
	`

	sql := "unsupported driver " + st.dbDriverName

	if st.dbDriverName == "mysql" {
		sql = sqlMysql
	}
	if st.dbDriverName == "postgres" {
		sql = sqlPostgres
	}
	if st.dbDriverName == "sqlite" {
		sql = sqlSqlite
	}

	return sql
}
