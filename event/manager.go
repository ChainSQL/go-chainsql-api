package event

import (
	"encoding/hex"
	"log"
	"strings"
	"sync"

	"github.com/ChainSQL/go-chainsql-api/common"
	"github.com/ChainSQL/go-chainsql-api/data"
	"github.com/ChainSQL/go-chainsql-api/export"
	"github.com/ChainSQL/go-chainsql-api/util"
	"github.com/buger/jsonparser"
)

// Manager manages the subscription
type Manager struct {
	txCache          map[string]export.Callback
	contractCache    map[string]bool
	ctrEventCache    map[string]chan *data.Log
	tableCache       map[string]export.Callback
	ledgerCloseCache []export.Callback
	muxTx            *sync.Mutex
	muxTable         *sync.Mutex
}

// NewEventManager is constructor for EventManager
func NewEventManager() *Manager {
	return &Manager{
		txCache:          make(map[string]export.Callback),
		tableCache:       make(map[string]export.Callback),
		contractCache:    make(map[string]bool),
		ctrEventCache:    make(map[string]chan *data.Log),
		ledgerCloseCache: make([]export.Callback, 0, 10),
		muxTx:            new(sync.Mutex),
		muxTable:         new(sync.Mutex),
	}
}

// SubscribeTable subscribe a table and set a callback function
func (e *Manager) SubscribeTable(name string, owner string, callback export.Callback) {
	e.muxTable.Lock()
	e.tableCache[name+owner] = callback
	e.muxTable.Unlock()
}

// UnSubscribeTable cancel the subscription
func (e *Manager) UnSubscribeTable(name string, owner string) {
	e.muxTable.Lock()
	delete(e.tableCache, name+owner)
	e.muxTable.Unlock()
}

//SubscribeTx subscribe a transaction
func (e *Manager) SubscribeTx(hash string, callback export.Callback) {
	e.muxTx.Lock()
	e.txCache[hash] = callback
	e.muxTx.Unlock()
}

//UnSubscribeTx unsubscribe a transaction
func (e *Manager) UnSubscribeTx(hash string) {
	e.muxTx.Lock()
	delete(e.txCache, hash)
	e.muxTx.Unlock()
}

//SubscribeCtrAddr subscribe a contract
func (e *Manager) SubscribeCtrAddr(address string, ok bool) {
	e.muxTx.Lock()
	e.contractCache[address] = ok
	e.muxTx.Unlock()
}

//UnSubscribeCtrAddr unsubscribe a contract
func (e *Manager) UnSubscribeCtrAddr(address string) {
	e.muxTx.Lock()
	delete(e.contractCache, address)
	e.muxTx.Unlock()
}

func (e *Manager) RegisterCtrEvent(eventSign string, contractMsgCh chan *data.Log) {
	e.muxTx.Lock()
	e.ctrEventCache[strings.ToUpper(eventSign[2:])] = contractMsgCh
	e.muxTx.Unlock()
}

func (e *Manager) UnRegisterCtrEvent(eventSign string, contractMsgCh chan *data.Log) {
	e.muxTx.Lock()
	delete(e.ctrEventCache, eventSign)
	e.muxTx.Unlock()
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

func (e *Manager) OnContractMsg(msg string) {
	log.Println(msg)
	ctrEventTopic, err := jsonparser.GetString([]byte(msg), "ContractEventTopics", "[0]")
	if err != nil {
		return
	}

	logRaw := &data.Log{}
	contractAddr, err := jsonparser.GetString([]byte(msg), "ContractAddress")
	if err != nil {
		return
	}
	logRaw.Address = contractAddr
	ctrEventInfo, err := jsonparser.GetString([]byte(msg), "ContractEventInfo")
	if err != nil {
		return
	}
	ctrEventInfoHex, err := hex.DecodeString(ctrEventInfo)
	if err != nil {
		return
	}
	logRaw.Data = ctrEventInfoHex
	ctrEventTopics, _, _, err := jsonparser.Get([]byte(msg), "ContractEventTopics")
	if err != nil {
		return
	}
	_, _ = jsonparser.ArrayEach(ctrEventTopics, func(value []byte, dataType jsonparser.ValueType, offset int, errin error) {
		valueStr := string(value)
		valueHex, err := hex.DecodeString(valueStr)
		if err != nil {
			return
		}
		logRaw.Topics = append(logRaw.Topics, common.BytesToHash(valueHex))
	})

	e.muxTx.Lock()
	if _, ok := e.contractCache[contractAddr]; ok {
		if eventMsgCh, ok := e.ctrEventCache[ctrEventTopic]; ok {
			eventMsgCh <- logRaw
		}
	}
	e.muxTx.Unlock()
	// if (data.hasOwnProperty("ContractEventTopics")) {
	// 	data.ContractEventTopics.map(function (topic, index) {
	// 		data.ContractEventTopics[index] = "0x" + data.ContractEventTopics[index].toLowerCase();
	// 	});
	// }
	// if (data.hasOwnProperty("ContractEventInfo") && data.ContractEventInfo !== "") {
	// 	data.ContractEventInfo = "0x" + data.ContractEventInfo;
	// }
	// let key = data.ContractEventTopics[0];
	// if (that.cache[key]) {
	// 	let contractObj = that.cache[data.ContractAddress];
	// 	let currentEvent = contractObj.options.jsonInterface.find(function (json) {
	// 		return (json.type === 'event' && json.signature === '0x' + key.replace('0x', ''));
	// 	});
	// 	let output = contractObj._decodeEventABI(currentEvent, data);
	// 	that.cache[key](null, output);
	// 	// delete that.cache[key];
	// 	// let keyIndex = contractObj.registeredEvent.indexOf(key);
	// 	// contractObj.registeredEvent.splice(keyIndex,1);
	// }
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
	e.muxTx.Lock()
	if cb, ok := e.txCache[txid]; ok {
		cb(msg)
	}
	e.muxTx.Unlock()

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
			e.UnSubscribeTx(txid)
		}
	} else {
		e.UnSubscribeTx(txid)
	}
}

//OnTableMsg trigger the callback
func (e *Manager) OnTableMsg(msg string) {

}
