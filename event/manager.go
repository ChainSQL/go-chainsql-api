package event

import (
	"log"

	"github.com/ChainSQL/go-chainsql-api/export"
	"github.com/ChainSQL/go-chainsql-api/util"
	"github.com/buger/jsonparser"
)

// Manager manages the subscription
type Manager struct {
	txCache          map[string]export.Callback
	tableCache       map[string]export.Callback
	ledgerCloseCache []export.Callback
}

// NewEventManager is constructor for EventManager
func NewEventManager() *Manager {
	return &Manager{
		txCache:          make(map[string]export.Callback),
		tableCache:       make(map[string]export.Callback),
		ledgerCloseCache: make([]export.Callback, 0, 10),
	}
}

// SubscribeTable subscribe a table and set a callback function
func (e *Manager) SubscribeTable(name string, owner string, callback export.Callback) {
	e.tableCache[name+owner] = callback
}

// UnSubscribeTable cancel the subscription
func (e *Manager) UnSubscribeTable(name string, owner string) {
	delete(e.tableCache, name+owner)
}

//SubscribeTx subscribe a transaction
func (e *Manager) SubscribeTx(hash string, callback export.Callback) {
	e.txCache[hash] = callback
}

//UnSubscribeTx unsubscribe a transaction
func (e *Manager) UnSubscribeTx(hash string) {
	delete(e.txCache, hash)
}

// SubscribeLedger subscribe ledgerClosed
func (e *Manager) SubscribeLedger(callback export.Callback) {
	e.ledgerCloseCache = append(e.ledgerCloseCache, callback)
}

// OnLedgerClosed trigger the callback
func (e *Manager) OnLedgerClosed(msg string) {
	for i := 0; i < len(e.ledgerCloseCache); i++ {
		e.ledgerCloseCache[i](msg)
	}
}

// OnSingleTransaction trigger the callback
func (e *Manager) OnSingleTransaction(msg string) {
	// log.Println(msg)
	txid, err := jsonparser.GetString([]byte(msg), "transaction", "hash")
	if err != nil {
		log.Printf("OnSingleTransaction error:%s\n", err)
		return
	}
	//trigger callback
	if cb, ok := e.txCache[txid]; ok {
		cb(msg)
	}

	txType, err := jsonparser.GetString([]byte(msg), "transaction", "TransactionType")
	if err != nil {
		// log.Printf("OnSingleTransaction error:%s\n", err)
		return
	}
	//remove subscription
	if util.IsChainsqlType(txType) {
		status, err := jsonparser.GetString([]byte(msg), "status")
		if err != nil {
			log.Printf("OnSingleTransaction error:%s\n", err)
			return
		}
		if util.ValidateSuccess != status {
			delete(e.txCache, txid)
		}
	} else {
		delete(e.txCache, txid)
	}
}

//OnTableMsg trigger the callback
func (e *Manager) OnTableMsg(msg string) {

}
