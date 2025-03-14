package customstore

// StoreInterface defines a custom store

type StoreInterface interface {
	// AutoMigrate migrates the tables
	AutoMigrate() error

	// EnableDebug - enables the debug option
	EnableDebug(debug bool)

	// RecordCreate creates a new record
	RecordCreate(record *Record) error

	// RecordFindByID finds a record by ID
	RecordFindByID(id string) (*Record, error)

	// RecordUpdate updates a record
	RecordUpdate(record *Record) error
}
