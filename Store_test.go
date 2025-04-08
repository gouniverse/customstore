package customstore

import (
	"database/sql"
	"os"
	"testing"

	"github.com/gouniverse/uid"
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

func TestNewStore(t *testing.T) {
	db := InitDB("test_data_store_new.db")

	// Test with valid options
	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_new",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	// Test with missing table name
	_, err = NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "",
		AutomigrateEnabled: true,
	})

	if err == nil {
		t.Fatalf("Store should not be created without table name")
	}

	// Test with missing database
	_, err = NewStore(NewStoreOptions{
		DB:                 nil,
		TableName:          "data_new",
		AutomigrateEnabled: true,
	})

	if err == nil {
		t.Fatalf("Store should not be created without database")
	}
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

	record := NewRecord("person")
	err = store.RecordCreate(record)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	if len(record.ID()) != 32 {
		t.Fatalf("Record ID != 32 but %s", record.ID())
	}

	if record.CreatedAt() == "" {
		t.Fatalf("Record CreatedAt is empty")
	}

	if record.UpdatedAt() == "" {
		t.Fatalf("Record UpdatedAt is empty")
	}

	if record.IsSoftDeleted() {
		t.Fatalf("Record should not be soft deleted")
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

	record := NewRecord("person")
	record.SetPayloadMap(map[string]interface{}{
		"name": "Jon",
	})
	err = store.RecordCreate(record)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	if len(record.ID()) != 32 {
		t.Fatalf("Record ID != 3 but %s", record.ID())
	}

	retrievedRecord, errFind := store.RecordFindByID(record.ID())

	if errFind != nil {
		t.Fatalf("Record could not be found: " + errFind.Error())
	}

	if retrievedRecord == nil {
		t.Fatalf("Record must not be NIL")
	}

	if retrievedRecord.Payload() != `{"name":"Jon"}` {
		t.Fatalf("Record payload must be {\"name\":\"Jon\"} but found %s", retrievedRecord.Payload())
	}

	// Test with non-existent ID
	retrievedRecord, errFind = store.RecordFindByID(uid.HumanUid())

	if errFind != nil {
		t.Fatalf("Record could not be found: " + errFind.Error())
	}

	if retrievedRecord != nil {
		t.Fatalf("Record must be NIL")
	}
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

	record := NewRecord(`person`)
	record.SetPayloadMap(map[string]any{
		`first_name`: `John`,
		`last_name`:  `Doe`,
	})

	err = store.RecordCreate(record)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	retrievedRecord, errFind := store.RecordFindByID(record.ID())

	if errFind != nil {
		t.Fatalf("Record could not be found: " + errFind.Error())
	}

	if retrievedRecord == nil {
		t.Fatalf("Record must not be NIL")
	}

	if retrievedRecord.Payload() != `{"first_name":"John","last_name":"Doe"}` {
		t.Fatal("Record data must be", record.Payload(), " found: ", retrievedRecord.Payload())
	}

	retrievedRecord.SetPayloadMap(map[string]any{
		`first_name`: `Jane`,
		`last_name`:  `Smith`,
		`country`:    `GB`,
	})

	err = store.RecordUpdate(retrievedRecord)

	if err != nil {
		t.Fatalf("Record could not be updated: " + err.Error())
	}

	retrievedRecord2, errFind := store.RecordFindByID(record.ID())

	if errFind != nil {
		t.Fatalf("Record could not be found: " + errFind.Error())
	}

	if retrievedRecord2 == nil {
		t.Fatalf("Record must not be NIL")
	}

	if retrievedRecord2.Payload() != `{"country":"GB","first_name":"Jane","last_name":"Smith"}` {
		t.Fatal("Record data must be", retrievedRecord.Payload(), " found: ", retrievedRecord2.Payload())
	}

	if retrievedRecord2.UpdatedAt() == record.UpdatedAt() {
		t.Fatal("Record UpdatedAt must be different")
	}
}

func TestRecordDelete(t *testing.T) {
	db := InitDB("test_data_store_record_delete.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	record := NewRecord("person")
	err = store.RecordCreate(record)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	err = store.RecordDelete(record)

	if err != nil {
		t.Fatalf("Record could not be deleted: " + err.Error())
	}

	retrievedRecord, errFind := store.RecordFindByID(record.ID())

	if errFind != nil {
		t.Fatalf("Record could not be found: " + errFind.Error())
	}

	if retrievedRecord != nil {
		t.Fatalf("Record must be NIL")
	}

	// Test with non-existent ID
	err = store.RecordDeleteByID(uid.HumanUid())

	if err != nil {
		t.Fatalf("Record could not be deleted: " + err.Error())
	}
}

func TestRecordSoftDelete(t *testing.T) {
	db := InitDB("test_data_store_record_soft_delete.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_soft_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	record := NewRecord("person")
	err = store.RecordCreate(record)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	err = store.RecordSoftDelete(record)

	if err != nil {
		t.Fatalf("Record could not be soft deleted: " + err.Error())
	}

	retrievedRecord, errFind := store.RecordFindByID(record.ID())

	if errFind != nil {
		t.Fatalf("Record could not be found: " + errFind.Error())
	}

	if retrievedRecord != nil {
		t.Fatalf("Record must be NIL")
	}

	// Test with non-existent ID
	err = store.RecordSoftDeleteByID(uid.HumanUid())

	if err != nil {
		t.Fatalf("Record could not be soft deleted: " + err.Error())
	}

	// Test with soft deleted included
	query := RecordQuery().SetSoftDeletedIncluded(true).SetID(record.ID())
	list, err := store.RecordList(query)

	if err != nil {
		t.Fatalf("Record could not be listed: " + err.Error())
	}

	if len(list) != 1 {
		t.Fatalf("Record list must be 1")
	}
}

func TestRecordList(t *testing.T) {
	db := InitDB("test_data_store_record_list.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_list",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	record1 := NewRecord("person")
	record1.SetPayloadMap(map[string]any{
		"name": "Jon",
	})
	err = store.RecordCreate(record1)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	record2 := NewRecord("person")
	record2.SetPayloadMap(map[string]any{
		"name": "Jane",
	})
	err = store.RecordCreate(record2)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	list, err := store.RecordList(RecordQuery())

	if err != nil {
		t.Fatalf("Record could not be listed: " + err.Error())
	}

	if len(list) != 2 {
		t.Fatalf("Record list must be 2")
	}

	if list[0].Payload() != `{"name":"Jon"}` && list[0].Payload() != `{"name":"Jane"}` {
		t.Fatalf("Record payload must be {\"name\":\"Jon\"} or {\"name\":\"Jane\"} but found %s", list[0].Payload())
	}

	if list[1].Payload() != `{"name":"Jon"}` && list[1].Payload() != `{"name":"Jane"}` {
		t.Fatalf("Record payload must be {\"name\":\"Jon\"} or {\"name\":\"Jane\"} but found %s", list[1].Payload())
	}
}

func TestRecordCount(t *testing.T) {
	db := InitDB("test_data_store_record_count.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_count",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	record1 := NewRecord("person")
	record1.SetPayloadMap(map[string]any{
		"name": "Jon",
	})
	err = store.RecordCreate(record1)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	record2 := NewRecord("person")
	record2.SetPayloadMap(map[string]any{
		"name": "Jane",
	})
	err = store.RecordCreate(record2)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	count, err := store.RecordCount(RecordQuery())

	if err != nil {
		t.Fatalf("Record could not be counted: " + err.Error())
	}

	if count != 2 {
		t.Fatalf("Record count must be 2")
	}
}

func TestRecordQuery(t *testing.T) {
	db := InitDB("test_data_store_record_query.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_query",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	record1 := NewRecord("person")
	record1.SetPayloadMap(map[string]any{
		"name": "Jon",
	})
	err = store.RecordCreate(record1)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	record2 := NewRecord("company")
	record2.SetPayloadMap(map[string]any{
		"name": "Acme",
	})
	err = store.RecordCreate(record2)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	// Test with type
	query := RecordQuery().SetType("person")
	list, err := store.RecordList(query)

	if err != nil {
		t.Fatalf("Record could not be listed: " + err.Error())
	}

	if len(list) != 1 {
		t.Fatalf("Record list must be 1")
	}

	if list[0].Payload() != `{"name":"Jon"}` {
		t.Fatalf("Record payload must be {\"name\":\"Jon\"} but found %s", list[0].Payload())
	}

	// Test with limit
	query = RecordQuery().SetLimit(1)
	list, err = store.RecordList(query)

	if err != nil {
		t.Fatalf("Record could not be listed: " + err.Error())
	}

	if len(list) != 1 {
		t.Fatalf("Record list must be 1")
	}

	// Test with offset
	query = RecordQuery().SetOffset(1)
	list, err = store.RecordList(query)

	if err != nil {
		t.Fatalf("Record could not be listed: " + err.Error())
	}

	if len(list) != 1 {
		t.Fatalf("Record list must be 1")
	}

	// Test with order by
	query = RecordQuery().SetOrderBy(COLUMN_CREATED_AT)
	list, err = store.RecordList(query)

	if err != nil {
		t.Fatalf("Record could not be listed: " + err.Error())
	}

	if len(list) != 2 {
		t.Fatalf("Record list must be 2")
	}
}

func TestRecordCreateWithEmptyID(t *testing.T) {
	db := InitDB("test_data_store_record_create_empty_id.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_create_empty_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	record := NewRecord("person")
	record.SetID("")
	err = store.RecordCreate(record)

	if err == nil {
		t.Fatalf("Record should not be created with empty ID")
	}

	if err.Error() != "record ID is required" {
		t.Fatal("Error should be 'record ID is required' but got:", err.Error())
	}
}

func TestRecordUpdateWithEmptyID(t *testing.T) {
	db := InitDB("test_data_store_record_update_empty_id.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_update_empty_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	record := NewRecord("person")
	record.SetID("")
	err = store.RecordUpdate(record)

	if err == nil {
		t.Fatalf("Record should not be updated with empty ID")
	}

	if err.Error() != "record id is required" {
		t.Fatal("Error should be 'record id is required' but got:", err.Error())
	}
}

func TestRecordDeleteWithEmptyID(t *testing.T) {
	db := InitDB("test_data_store_record_delete_empty_id.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_delete_empty_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	err = store.RecordDeleteByID("")

	if err == nil {
		t.Fatalf("Record should not be deleted with empty ID")
	}

	if err.Error() != "record id is empty" {
		t.Fatal("Error should be 'record id is empty' but got:", err.Error())
	}
}

func TestRecordFindByIDWithEmptyID(t *testing.T) {
	db := InitDB("test_data_store_record_find_empty_id.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_find_empty_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	_, err = store.RecordFindByID("")

	if err == nil {
		t.Fatalf("Record should not be found with empty ID")
	}

	if err.Error() != "record id is empty" {
		t.Fatal("Error should be 'record id is empty' but got:", err.Error())
	}
}

func TestRecordSoftDeleteWithEmptyID(t *testing.T) {
	db := InitDB("test_data_store_record_soft_delete_empty_id.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_soft_delete_empty_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	err = store.RecordSoftDeleteByID("")

	if err == nil {
		t.Fatalf("Record should not be soft deleted with empty ID")
	}

	if err.Error() != "record id is empty" {
		t.Fatal("Error should be 'record id is empty' but got:", err.Error())
	}
}

func TestRecordQueryPayloadSearch(t *testing.T) {
	db := InitDB("test_data_store_record_query_payload_search.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "data_record_query_payload_search",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	record1 := NewRecord("person")
	record1.SetPayloadMap(map[string]any{
		"name":    "Jon Doe",
		"country": "US",
	})
	err = store.RecordCreate(record1)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	record2 := NewRecord("person")
	record2.SetPayloadMap(map[string]any{
		"name":    "Jane Smith",
		"country": "GB",
	})
	err = store.RecordCreate(record2)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	record3 := NewRecord("company")
	record3.SetPayloadMap(map[string]any{
		"name":    "Acme Corp",
		"country": "US",
	})
	err = store.RecordCreate(record3)

	if err != nil {
		t.Fatalf("Record could not be created: " + err.Error())
	}

	// Test with payload search
	query := RecordQuery().AddPayloadSearch("Jon")
	list, err := store.RecordList(query)

	if err != nil {
		t.Fatalf("Record could not be listed: " + err.Error())
	}

	if len(list) != 1 {
		t.Fatalf("Record list must be 1")
	}

	if list[0].Payload() != `{"country":"US","name":"Jon Doe"}` {
		t.Fatalf("Record payload must be {\"country\":\"US\",\"name\":\"Jon Doe\"} but found %s", list[0].Payload())
	}

	// Test with payload search
	query = RecordQuery().AddPayloadSearch("US")
	list, err = store.RecordList(query)

	if err != nil {
		t.Fatalf("Record could not be listed: " + err.Error())
	}

	if len(list) != 2 {
		t.Fatalf("Record list must be 2")
	}

	// Test with payload search
	query = RecordQuery().AddPayloadSearch("Jane")
	list, err = store.RecordList(query)

	if err != nil {
		t.Fatalf("Record could not be listed: " + err.Error())
	}

	if len(list) != 1 {
		t.Fatalf("Record list must be 1")
	}

	if list[0].Payload() != `{"country":"GB","name":"Jane Smith"}` {
		t.Fatalf("Record payload must be {\"country\":\"GB\",\"name\":\"Jane Smith\"} but found %s", list[0].Payload())
	}

	// Test with payload search
	query = RecordQuery().AddPayloadSearch("Acme")
	list, err = store.RecordList(query)

	if err != nil {
		t.Fatalf("Record could not be listed: " + err.Error())
	}

	if len(list) != 1 {
		t.Fatalf("Record list must be 1")
	}

	if list[0].Payload() != `{"country":"US","name":"Acme Corp"}` {
		t.Fatalf("Record payload must be {\"country\":\"US\",\"name\":\"Acme Corp\"} but found %s", list[0].Payload())
	}

	// Test with payload search
	query = RecordQuery().AddPayloadSearch("Corp")
	list, err = store.RecordList(query)

	if err != nil {
		t.Fatalf("Record could not be listed: " + err.Error())
	}

	if len(list) != 1 {
		t.Fatalf("Record list must be 1")
	}

	if list[0].Payload() != `{"country":"US","name":"Acme Corp"}` {
		t.Fatalf("Record payload must be {\"country\":\"US\",\"name\":\"Acme Corp\"} but found %s", list[0].Payload())
	}
}
