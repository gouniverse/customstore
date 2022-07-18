package customstore

import (
	"database/sql"
	"os"
	"testing"

	// "time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreCreate(t *testing.T) {
	db := InitDB("test_data_store_create.db")

	store, err := NewStore(WithDb(db), WithTableName("data_create"), WithAutoMigrate(true))

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	// isOk, err := store.Set("post", "1234567890", 5)

	// if err != nil {
	// 	t.Fatalf("Cache could not be created: " + err.Error())
	// }

	// if isOk == false {
	// 	t.Fatalf("Cache could not be created")
	// }
}

func TestRecordCreate(t *testing.T) {
	db := InitDB("test_data_store_record_create.db")

	store, err := NewStore(WithDb(db), WithTableName("data_record_create"), WithAutoMigrate(true))

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	record := Record{
		Type: "person",
	}
	isOk, err := store.RecordCreate(&record)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	if isOk == false {
		t.Fatalf("Record could not be created")
	}

	if len(record.ID) != 32 {
		t.Fatalf("Record ID != 3 but %s", record.ID)
	}
}

func TestRecordFindByID(t *testing.T) {
	db := InitDB("test_data_store_record_find.db")

	store, err := NewStore(WithDb(db), WithTableName("data_record_find"), WithDebug(false), WithAutoMigrate(true))

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	record := Record{
		Type: "person",
	}
	record.SetMap(map[string]interface{}{
		"name": "Jon",
	})
	isOk, err := store.RecordCreate(&record)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	if isOk == false {
		t.Fatalf("Record could not be created")
	}

	if len(record.ID) != 32 {
		t.Fatalf("Record ID != 3 but %s", record.ID)
	}

	retrievedRecord, errFind := store.RecordFindByID(record.ID)

	if errFind != nil {
		t.Fatalf("Record could not be found: " + errFind.Error())
	}

	if retrievedRecord == nil {
		t.Fatalf("Record must not be NIL")
	}

	// log.Println(retrievedRecord.GetMap())
}
