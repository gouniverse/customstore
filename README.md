# customstore <a href="https://gitpod.io/#https://github.com/gouniverse/customstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

Stores a custom record to a database table.

## Installation
```
go get -u github.com/gouniverse/customstore
```

## Setup

```go
customStore = customstore.NewStore(customstore.NewStoreOptions{
	DB:                 databaseInstance,
	TableName:          "my_custom_record",
	AutomigrateEnabled: true,
	DebugEnabled:       false,
})
```

## Methods

- AutoMigrate() error - automigrate (creates) the session table
- DriverName(db *sql.DB) string - finds the driver name from database
- EnableDebug(debug bool) - enables / disables the debug option
- RecordCreate(record *Record) error
- RecordFindByID(id string) (*Record, error)