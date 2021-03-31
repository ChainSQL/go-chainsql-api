package net

import (
	"log"

	"github.com/buger/jsonparser"
)

// ServerInfo struct
type ServerInfo struct {
	FeeBase      int
	FeeRef       int
	DropsPerByte int
	LoadBase     int
	LoadFactor   int
	LedgerIndex  int
	TxnSuccess   int
	TxnFailure   int
	TxnCount     int
	Ledgerhash   string
	ServerStatus string
	Updated      bool
}

//NewServerInfo is constructor
func NewServerInfo() *ServerInfo {
	return &ServerInfo{
		Updated:      false,
		DropsPerByte: 976,
		FeeRef:       10,
		FeeBase:      10,
		LoadBase:     256,
		LoadFactor:   256,
	}
}

//Update update ServerInfo from json result
func (s *ServerInfo) Update(result string) {

	s.GetFieldInt(result, &s.FeeBase, "fee_base")
	s.GetFieldInt(result, &s.FeeRef, "fee_ref")
	s.GetFieldInt(result, &s.DropsPerByte, "drops_per_byte")
	s.GetFieldInt(result, &s.LoadBase, "load_base")
	s.GetFieldInt(result, &s.LoadFactor, "load_factor")
	s.GetFieldInt(result, &s.LedgerIndex, "ledger_index")
	s.GetFieldInt(result, &s.TxnSuccess, "txn_success")
	s.GetFieldInt(result, &s.TxnFailure, "txn_failure")
	s.GetFieldInt(result, &s.TxnCount, "txn_count")
	s.GetFieldString(result, &s.Ledgerhash, "ledger_hash")
	s.GetFieldString(result, &s.ServerStatus, "server_status")

	s.Updated = true
}

//GetFieldInt get value from json
func (s *ServerInfo) GetFieldInt(result string, field *int, fieldInJSON string) {
	nValue, err := jsonparser.GetInt([]byte(result), fieldInJSON)
	if err == nil {
		*field = int(nValue)
	} else {
		log.Printf("GetFieldInt error for field %s:%s\n", fieldInJSON, result)
	}
}

//GetFieldString get value from json
func (s *ServerInfo) GetFieldString(result string, field *string, fieldInJSON string) {
	sValue, err := jsonparser.GetString([]byte(result), fieldInJSON)
	if err == nil {
		*field = sValue
	} else {
		log.Printf("GetFieldString error for field %s:%s\n", fieldInJSON, result)
	}
}

// ComputeFee compute the basic transaction fee
func (s *ServerInfo) ComputeFee() int {
	if !s.Updated {
		return 0
	}
	feeUnit := float32(s.FeeBase) / float32(s.FeeRef)
	feeUnit *= float32(s.LoadFactor) / float32(s.LoadBase)
	fee := int(float32(s.FeeBase) * feeUnit * 1.1)
	return fee
}
