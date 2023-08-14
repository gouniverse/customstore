package customstore

import (
	"encoding/json"
	"log"
	"time"
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
