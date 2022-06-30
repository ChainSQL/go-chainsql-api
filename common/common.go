package common

type KeyType int

const (
	ECDSA     KeyType = 0
	Ed25519   KeyType = 1
	SoftGMAlg KeyType = 2
)

/*func (keyType KeyType) String() string {
	switch keyType {
	case secp256k1:
		return "secp256k1"
	case Ed25519:
		return "Ed25519"
	default:
		return "unknown key type"
	}
}*/

/*func (keyType KeyType) MarshalText() ([]byte, error) {

	return []byte(keyType.String()), nil
}*/

// Auth is the type with ws connection infomations
type Auth struct {
	Address string
	Secret  string
	Owner   string
}

//IRequest define interface for request
type IRequest interface {
	GetID() int64
	SetSchemaID(schemaID string) *RequestBase
}

//RequestBase contains fields that all requests will have
type RequestBase struct {
	Command  string `json:"command"`
	ID       int64  `json:"id,omitempty"`
	SchemaID string `json:"schema_id"`
}

// GetID  return id for request
func (r *RequestBase) GetID() int64 {
	return r.ID
}

func (r *RequestBase) SetSchemaID(schemaID string) *RequestBase {
	r.SchemaID = schemaID
	return r
}
