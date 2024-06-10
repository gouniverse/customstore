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

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_create",
		AutomigrateEnabled: true,
	})

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

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	record := Record{
		Type: "person",
	}
	err = store.RecordCreate(&record)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	if len(record.ID) != 32 {
		t.Fatalf("Record ID != 3 but %s", record.ID)
	}
}

func TestRecordFindByID(t *testing.T) {
	db := InitDB("test_data_store_record_find.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_find",
		AutomigrateEnabled: true,
	})

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
	err = store.RecordCreate(&record)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
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

func TestRecordUpdate(t *testing.T) {
	db := InitDB("test_data_store_record_update.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_update",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	record := NewRecord(`person`).SetMap(map[string]any{
		`first_name`: `John`,
		`last_name`:  `Doe`,
	})

	err = store.RecordCreate(record)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	retrievedRecord, errFind := store.RecordFindByID(record.ID)

	if errFind != nil {
		t.Fatalf("Record could not be found: " + errFind.Error())
	}

	if retrievedRecord == nil {
		t.Fatalf("Record must not be NIL")
	}

	if retrievedRecord.Data != `{"first_name":"John","last_name":"Doe"}` {
		t.Fatal("Record data must be", record.Data, " found: ", retrievedRecord.Data)
	}

	retrievedRecord.SetMap(map[string]any{
		`first_name`: `Jane`,
		`last_name`:  `Smith`,
		`country`:    `GB`,
	})

	err = store.RecordUpdate(retrievedRecord)

	if err != nil {
		t.Fatalf("Record could not be updated: " + err.Error())
	}

	retrievedRecord2, errFind := store.RecordFindByID(record.ID)

	if errFind != nil {
		t.Fatalf("Record could not be found: " + errFind.Error())
	}

	if retrievedRecord2 == nil {
		t.Fatalf("Record must not be NIL")
	}

	if retrievedRecord2.Data != `{"country":"GB","first_name":"Jane","last_name":"Smith"}` {
		t.Fatal("Record data must be", retrievedRecord.Data, " found: ", retrievedRecord2.Data)
	}

}
