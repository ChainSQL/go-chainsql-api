package common

import "fmt"

//Account define the account format
type Account struct {
	Address      string `json:"address"`
	PublicKey    string `json:"publicKey"`
	PublicKeyHex string `json:"publicKeyHex"`
	PrivateKey   string `json:"privateKey"`
}

// Auth is the type with ws connection infomations
type Auth struct {
	Address string
	Secret  string
	Owner   string
}

// TableName is sub-struct of Table
type TableName struct {
	TableName string
	NameInDB  string `json:"NameInDB,omitempty"`
}

//TableObj is sub-struct of Tables
type TableObj struct {
	Table TableName
}

// TableFields is common table operation fields
type TableFields struct {
	Tables  []TableObj
	Raw     string
	Account string
}

// TableTxFields contains fields that chainsql table transaction needs
type TableTxFields struct {
	TransactionType string
	OpType          int16
	TableFields
}

//NetFields contains fields that need request from network
type NetFields struct {
	Sequence           int
	Fee                string
	LastLedgerSequence int `json:"LastLedgerSequence,omitempty"`
}

//IRequest define interface for request
type IRequest interface {
	GetID() int64
}

//RequestBase contains fields that all requests will have
type RequestBase struct {
	Command string `json:"command"`
	ID      int64  `json:"id,omitempty"`
}

// GetID  return id for request
func (r *RequestBase) GetID() int64 {
	return r.ID
}

// FormatTables create the Tables json array in Chainsql transaction
func FormatTables(name string, nameInDB string) []TableObj {
	return []TableObj{
		{
			Table: TableName{
				TableName: fmt.Sprintf("%x", name),
				NameInDB:  nameInDB,
			},
		},
	}
}

func FormatTablesForGet(name string, nameInDB string) []TableObj {
	return []TableObj{
		{
			Table: TableName{
				TableName: name,
				NameInDB:  nameInDB,
			},
		},
	}
}
