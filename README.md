# customstore <a href="https://gitpod.io/#https://github.com/gouniverse/customstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/gouniverse/customstore/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/gouniverse/customstore/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/customstore)](https://goreportcard.com/report/github.com/gouniverse/customstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/customstore)](https://pkg.go.dev/github.com/gouniverse/customstore)

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