package customstore

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
