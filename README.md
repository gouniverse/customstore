# customstore <a href="https://gitpod.io/#https://github.com/gouniverse/customstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/gouniverse/customstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/gouniverse/customstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/customstore)](https://goreportcard.com/report/github.com/gouniverse/customstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/customstore)](https://pkg.go.dev/github.com/gouniverse/customstore)

**customstore** is a Go package that provides a flexible way to store and manage custom records in a database table. It simplifies common database operations like creating, retrieving, updating, and deleting records.

## Features

*   **Easy Setup:** Quickly integrate with your existing database.
*   **Customizable Records:** Define your own record types and data structures.
*   **Automatic Migration:** Automatically create the necessary database table.
*   **CRUD Operations:** Supports standard Create, Read, Update, and Delete operations.
*   **Flexible Queries:** Query records based on various criteria.
*   **Soft Deletes:** Option to soft delete records instead of permanent deletion.
*   **Payload Search:** Search for records based on content within the payload.
*   **Debug Mode:** Enable debug mode for detailed logging.

## Installation

```bash
go get -u github.com/gouniverse/customstore
```

## Setup

```go
// Example with SQLite
db, err := sql.Open("sqlite3", "mydatabase.db")
if err != nil {
    panic(err)
}
defer db.Close()

// Initialize the store
customStore, err := customstore.NewStore(customstore.NewStoreOptions{
	DB:                 db,
	TableName:          "my_custom_records",
	AutomigrateEnabled: true,
	DebugEnabled:       false,
})

if err != nil {
    panic(err)
}
```

## Methods

- AutoMigrate() error - automigrate (creates) the session table
- DriverName(db *sql.DB) string - finds the driver name from database
- EnableDebug(debug bool) - enables / disables the debug option
- RecordCreate(record *Record) error
- RecordFindByID(id string) (*Record, error)

## Core Concepts

### Records
A Record represents a single entry in your custom data store. Each record has:

Type: A string that categorizes the record (e.g., "user", "product", "order").
ID: A unique identifier for the record.
Payload: A JSON-encoded string containing the record's data.
CreatedAt: A timestamp indicating when the record was created.
UpdatedAt: A timestamp indicating when the record was last updated.
DeletedAt: A timestamp indicating when the record was soft-deleted (if applicable).

### Store
The Store is the main interface for interacting with your custom data store. It provides methods for:

Creating records.
Retrieving records by ID.
Updating records.
Deleting records (both hard and soft deletes).
Listing records based on various criteria.
Counting records.

### RecordQuery
The RecordQuery struct allows you to build complex queries to filter and retrieve records. You can specify:

Record type.
ID.
Limit and offset for pagination.
Order by clause.
Whether to include soft-deleted records.
Payload search terms.

## Usage Examples

### Creating a Record

```go
record := customstore.NewRecord("person")
record.SetPayloadMap(map[string]interface{}{
	"name": "John Doe",
	"age":  30,
})

err := store.RecordCreate(record)
if err != nil {
	panic(err)
}
```

### Finding a Record by ID

```go
record, err := store.RecordFindByID("1234567890")
if err != nil {
	panic(err)
}
```

### Updating a Record
```go
record, err := store.RecordFindByID("1234567890")
if err != nil {
	panic(err)
}

record.SetPayloadMap(map[string]interface{}{
	"name": "John Doe",
	"age":  30,
})

err = store.RecordUpdate(record)
if err != nil {
	panic(err)
}
```

### Deleting a Record (Hard Delete)
```go
record, err := store.RecordFindByID("1234567890")
if err != nil {
	panic(err)
}

err = store.RecordDelete(record)
if err != nil {
	panic(err)
}
```

### Soft Deleting a Record
```go
record, err := store.RecordFindByID("1234567890")
if err != nil {
	panic(err)
}

err = store.RecordSoftDelete(record)
if err != nil {
	panic(err)
}
```

### Listing Records
```go
query := customstore.RecordQuery().SetType("person").SetLimit(10)
list, err := store.RecordList(query)
if err != nil {
	panic(err)
}
```

### Counting Records
```go
query := customstore.RecordQuery().SetType("person")
count, err := store.RecordCount(query)
if err != nil {
	panic(err)
}
```

### Payload Search
```go
query := customstore.RecordQuery().SetType("person").
    AddPayloadSearch(`"status": "active"`).
	AddPayloadSearch(`"name": "John"`)
list, err := store.RecordList(query)
if err != nil {
	panic(err)
}
```

### Soft Deleted Records
```go
query := customstore.RecordQuery().SetType("person").SetSoftDeletedIncluded(true)
list, err := store.RecordList(query)
if err != nil {
	panic(err)
}
```


## API Reference

### Store Methods

NewStore(options NewStoreOptions) (*Store, error) - Creates a new store instance.
options: A NewStoreOptions struct containing the database connection, table name, and other configuration options.
AutoMigrate() error - Automigrates (creates) the session table.
DriverName(db *sql.DB) string - Finds the driver name from the database.
EnableDebug(debug bool) - Enables/disables the debug option.
RecordCreate(record *Record) error - Creates a new record.
RecordFindByID(id string) (*Record, error) - Finds a record by its ID.
RecordUpdate(record *Record) error - Updates an existing record.
RecordDelete(record *Record) error - Deletes a record.
RecordDeleteByID(id string) error - Deletes a record by its ID.
RecordSoftDelete(record *Record) error - Soft deletes a record.
RecordSoftDeleteByID(id string) error - Soft deletes a record by its ID.
RecordList(query *RecordQuery) ([]*Record, error) - Lists records based on a query.
RecordCount(query *RecordQuery) (int64, error) - Counts records based on a query.
RecordQuery Methods
SetID(id string) *RecordQuery - Sets the ID to search for.
SetType(recordType string) *RecordQuery - Sets the record type to search for.
SetLimit(limit int) *RecordQuery - Sets the maximum number of records to return.
SetOffset(offset int) *RecordQuery - Sets the offset for the records to return.
SetOrderBy(orderBy string) *RecordQuery - Sets the order by clause.
SetSoftDeletedIncluded(softDeletedIncluded bool) *RecordQuery - Sets whether to include soft deleted records.
AddPayloadSearch(payloadSearch string) *RecordQuery - Adds a payload search term.
Contributing
Contributions are welcome! Please feel free to submit a pull request.