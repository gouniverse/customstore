package customstore

import "github.com/gouniverse/sb"

// SqlCreateUserTable returns a SQL string for creating the user table
func (store *storeImplementation) SqlCreateTable() string {
	sql := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
		Table(store.tableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_RECORD_TYPE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 100,
			Unique: true,
		}).
		Column(sb.Column{
			Name: COLUMN_RECORD_DATA,
			Type: sb.COLUMN_TYPE_LONGTEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:     COLUMN_DELETED_AT,
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		CreateIfNotExists()

	return sql
}
