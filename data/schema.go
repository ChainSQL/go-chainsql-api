package data

type CreateSchema struct {
	SchemaName       string
	WithState        bool
	SchemaAdmin      string `json:"SchemaAdmin,omitempty"`
	AnchorLedgerHash string `json:"AnchorLedgerHash,omitempty"`
	PeerList         []Peer
	Validators       []Validator
}

type ModifySchema struct {
	OpType     uint16
	Validators []Validator
	PeerList   []Peer
	SchemaID   string
}

type Peer struct {
	Peer Endpoint
}
type PeerFormat struct {
	Peer EndpointFormat
}
type Endpoint struct {
	Endpoint string
}

type EndpointFormat struct {
	Endpoint VariableLength
}

type Validator struct {
	Validator PublicKeyObj
}
type ValidatorFormat struct {
	Validator PublicKeyObjFormat
}
type PublicKeyObj struct {
	PublicKey string
}

type PublicKeyObjFormat struct {
	PublicKey VariableLength
}
