package core

import (
	"encoding/json"
	"errors"
	"strings"

	. "github.com/ChainSQL/go-chainsql-api/data"
	"github.com/ChainSQL/go-chainsql-api/net"
	"github.com/ChainSQL/go-chainsql-api/util"
)

//OpInfo is the opearting details
type OpInfo struct {
	Raw   string
	Exec  uint16
	Query []interface{}
}

type TableGetJSON struct {
	Tables      []TableObjForGet
	Raw         string `json:"Raw,omitempty"`
	Account     string
	Owner       string
	LedgerIndex int
}

//Table is used to process insert/delete/update/get operation
type Table struct {
	name   string
	client *net.Client
	op     *OpInfo
	SubmitBase
}

//NewTable creates a Table object
func NewTable(name string, client *net.Client) *Table {
	table := &Table{
		name:   name,
		client: client,
		op: &OpInfo{
			Query: make([]interface{}, 0),
		},
	}
	table.SubmitBase.client = table.client
	table.SubmitBase.IPrepare = table
	return table
}

//Insert method insert data to a table
func (t *Table) Insert(value string) *Table {
	t.op.Exec = util.RInsert
	t.op.Raw = value
	return t
}

//Get is used to select data from table
//parameter raw is a json-object string like
// {"$and":[{ "id": 2},{ "name": "张三"}]}
func (t *Table) Get(raw string) *Table {
	t.op.Exec = util.RGet
	if raw != "" {
		var jsonObj interface{}
		json.Unmarshal([]byte(raw), &jsonObj)
		t.op.Query = append(t.op.Query, jsonObj)
	}
	return t
}

//Limit is used to limit the record count
//parameter is a json-object string like {"total":10, "index":0}
func (t *Table) Limit(limit string) *Table {
	type LimitObj struct {
		Limit interface{} `json:"$limit"`
	}
	var jsonObj interface{}
	json.Unmarshal([]byte(limit), &jsonObj)
	limitObj := LimitObj{
		Limit: jsonObj,
	}
	t.op.Query = append(t.op.Query, limitObj)
	return t
}

//Order is used to order the result record
//parameter is a json-array string like [{"id":1},{"name":-1}]
func (t *Table) Order(raw string) *Table {
	type OrderObj struct {
		Order []interface{} `json:"$order"`
	}
	var jsonObj []interface{}
	json.Unmarshal([]byte(raw), &jsonObj)
	orderObj := OrderObj{
		Order: jsonObj,
	}
	t.op.Query = append(t.op.Query, orderObj)
	return t
}

func (t *Table) WithFields(raw string) *Table {
	var jsonObj interface{}
	json.Unmarshal([]byte(raw), &jsonObj)
	t.op.Query = append([]interface{}{jsonObj}, t.op.Query...)
	return t
}

//Request is used to end a get operation and send request
func (t *Table) Request() (string, error) {
	if t.op.Exec != util.RGet {
		return "", errors.New("Not a get operation")
	}
	// check withFields
	addedWithFields := true
	if len(t.op.Query) == 0 {
		addedWithFields = false
	} else {
		// fmt.Printf("WithFields len:%d\n",len(t.op.Query))
		str, err := json.Marshal(t.op.Query[0])
		if err != nil {
			return "", err
		}
		brackets := strings.Index(string(str), "[")
		if brackets != 0 {
			addedWithFields = false
		}
	}
	if !addedWithFields {
		var withFields interface{} = []string{}
		t.op.Query = append([]interface{}{withFields}, t.op.Query...)
	}
	strQuery, err := json.Marshal(t.op.Query)
	if err != nil {
		return "", err
	}
	// fmt.Printf("Query string:%s\n",string(strQuery))

	data := &TableGetJSON{}
	nameInDB, err := t.client.GetNameInDB(t.client.Auth.Owner, t.name)
	if err != nil {
		return "", err
	}
	if t.client.ServerInfo.Updated {
		data.LedgerIndex = t.client.ServerInfo.LedgerIndex
	} else {
		ledgerIndex, err := t.client.GetLedgerVersion()
		if err != nil {
			return "", err
		}
		data.LedgerIndex = ledgerIndex
	}
	data.Tables = FormatTablesForGet(t.name, nameInDB)
	data.Raw = string(strQuery)
	data.Account = t.client.Auth.Address
	data.Owner = t.client.Auth.Owner
	return t.client.GetTableData(data, false)
}

//PrepareTx prepare tx json for submit
func (t *Table) PrepareTx() (Signer, error) {
	tx := &SQLStatement{}
	seq, nameInDB, err := net.PrepareTable(t.client, t.name)
	if err != nil {
		// log.Println(err)
		return nil, err
	}
	account, err := NewAccountFromAddress(t.client.Auth.Address)
	if err != nil {
		// log.Println(err)
		return nil, err
	}
	owner, err := NewAccountFromAddress(t.client.Auth.Owner)
	if err != nil {
		// log.Println(err)
		return nil, err
	}
	var valRaw = VariableLength(t.op.Raw) //fmt.Sprintf("%x", t.op.Raw)
	tx.TransactionType = SQLSTATEMENT
	tx.Tables = FormatTables(t.name, nameInDB)
	tx.OpType = t.op.Exec
	tx.Raw = &valRaw
	tx.Account = *account
	tx.Owner = *owner
	tx.Sequence = seq
	fee, lastLedgerSeq, err := net.PrepareLastLedgerSeqAndFee(t.client)
	if err != nil {
		return nil, err
	}
	tx.LastLedgerSequence = &lastLedgerSeq
	fee += util.GetExtraFee(t.op.Raw, t.client.ServerInfo.DropsPerByte)
	finalFee, err := NewNativeValue(fee)
	if err != nil {
		return nil, err
	}
	tx.Fee = *finalFee
	return tx, nil
}
