package customstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
)

// RecordInterface represents an record for accessing the API
type RecordInterface interface {
	dataobject.DataObjectInterface

	IsSoftDeleted() bool

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string)

	ID() string
	SetID(id string)

	Type() string
	SetType(t string)

	Meta(name string) string
	SetMeta(name, value string) error

	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	Memo() string
	SetMemo(memo string)

	Payload() string
	SetPayload(payload string)

	PayloadMap() (map[string]any, error)
	SetPayloadMap(payloadMap map[string]any) error

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string)

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string)
}
