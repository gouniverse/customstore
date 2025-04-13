package customstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayloadMapKey(t *testing.T) {
	record := NewRecord("test")

	// Test setting a key
	err := record.SetPayloadMapKey("name", "John")
	assert.NoError(t, err)

	// Test getting the key
	value, err := record.PayloadMapKey("name")
	assert.NoError(t, err)
	assert.Equal(t, "John", value)

	// Test getting a non-existent key
	value, err = record.PayloadMapKey("nonexistent")
	assert.NoError(t, err)
	assert.Nil(t, value)

	// Test setting multiple keys
	err = record.SetPayloadMapKey("age", 30)
	assert.NoError(t, err)

	value, err = record.PayloadMapKey("age")
	assert.NoError(t, err)
	assert.Equal(t, float64(30), value) // JSON numbers are decoded as float64

	// Verify the entire payload map
	payloadMap, err := record.PayloadMap()
	assert.NoError(t, err)
	assert.Equal(t, "John", payloadMap["name"])
	assert.Equal(t, float64(30), payloadMap["age"])
}
