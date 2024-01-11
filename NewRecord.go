package customstore

func NewRecord(recordType string) *Record {
	record := Record{
		Type: recordType,
	}

	return &record
}
