package core

import (
	"encoding/json"
	"strings"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-chainsql-api/common"
	"github.com/go-chainsql-api/net"
	"github.com/go-chainsql-api/util"
)

//OpInfo is the opearting details
type OpInfo struct {
	Raw  string
	Exec int16
	Query []interface{}
}

//TableJSON specifies the table operation format
type TableJSON struct {
	common.TableTxFields
	common.NetFields
	Owner         string
	AutoFillField string `json:"AutoFillField,omitempty"`
}

type TableGetJSON struct{
	common.TableFields
	Owner         string
	LedgerIndex   int
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
		op:     &OpInfo{
			Query: make([]interface{},0),
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
//parameter raw is a json-object string like {"$and":[{ "id": 2},{ "name": "张三"}]}
func (t *Table) Get(raw string) *Table{
	t.op.Exec = util.RGet
	if raw != ""{
		var jsonObj interface{}
		json.Unmarshal([]byte(raw), &jsonObj)
		t.op.Query = append(t.op.Query,jsonObj)
	}
	return t
}

//Limit is used to limit the record count
//parameter is a json-object string like {"total":10, "index":0}
func (t *Table)Limit(limit string) *Table{
	type LimitObj struct{
		Limit interface{} `json:"$limit"`
	}
	var jsonObj interface{}
	json.Unmarshal([]byte(limit), &jsonObj)
	limitObj := LimitObj{
		Limit:jsonObj,
	}
	t.op.Query = append(t.op.Query,limitObj)
	return t
}
//Order is used to order the result record
//parameter is a json-array string like [{"id":1},{"name":-1}]
func (t *Table)Order(raw string) *Table{
	type OrderObj struct{
		Order []interface{} `json:"$order"`
	}
	var jsonObj []interface{}
	json.Unmarshal([]byte(raw), &jsonObj)
	orderObj := OrderObj{
		Order:jsonObj,
	}
	t.op.Query = append(t.op.Query,orderObj)
	return t
}

func (t *Table)WithFields(raw string) *Table{
	var jsonObj interface{}
	json.Unmarshal([]byte(raw), &jsonObj)
	t.op.Query = append([]interface{}{jsonObj},t.op.Query...)
	return t
}

//Request is used to end a get operation and send request
func (t *Table)Request() (string,error){
	if t.op.Exec != util.RGet{
		return "",errors.New("Not a get operation")
	}
	// check withFields
	addedWithFields := true
	if len(t.op.Query) == 0 {
		addedWithFields = false
	}else{
		// fmt.Printf("WithFields len:%d\n",len(t.op.Query))
		str,err := json.Marshal(t.op.Query[0])
		if err != nil{
			return "",err
		}
		brackets := strings.Index(string(str), "[")
		if brackets != 0 {
			addedWithFields = false
		}
	}
	if !addedWithFields {
		var withFields interface{} = []string{}
		t.op.Query = append([]interface{}{withFields},t.op.Query...)
	}
	strQuery,err := json.Marshal(t.op.Query)
	if err != nil{
		return "",err
	}
	// fmt.Printf("Query string:%s\n",string(strQuery))

	data := &TableGetJSON{}
	nameInDB, err := t.client.GetNameInDB(t.client.Auth.Address, t.name)
	if err != nil {
		return "",err
	}
	if t.client.ServerInfo.Updated {
		data.LedgerIndex = t.client.ServerInfo.LedgerIndex
	} else {
		ledgerIndex, err := t.client.GetLedgerVersion()
		if err != nil{
			return "",err
		}
		data.LedgerIndex = ledgerIndex
	}
	data.Tables = common.FormatTablesForGet(t.name, nameInDB)
	data.Raw = string(strQuery)
	data.Account = t.client.Auth.Address
	data.Owner = t.client.Auth.Owner
	return t.client.GetTableData(data)
}

//PrepareTx prepare tx json for submit
func (t *Table) PrepareTx() (TxJSON, error) {
	tx := &TableJSON{}
	seq, nameInDB, err := net.PrepareTable(t.client, t.client.Auth.Address, t.name)
	if err != nil {
		// log.Println(err)
		return nil, err
	}

	tx.TransactionType = util.SQLStatement
	tx.Tables = common.FormatTables(t.name, nameInDB)
	tx.OpType = t.op.Exec
	tx.Raw = fmt.Sprintf("%x", t.op.Raw)
	tx.Account = t.client.Auth.Address
	tx.Owner = t.client.Auth.Owner
	tx.Sequence = seq
	fee := 10
	if t.client.ServerInfo.Updated {
		tx.LastLedgerSequence = t.client.ServerInfo.LedgerIndex + 20
		fee = t.client.ServerInfo.ComputeFee()
	} else {
		ledgerIndex, err := t.client.GetLedgerVersion()
		if err != nil {
			return nil, err
		}
		tx.LastLedgerSequence = ledgerIndex + 20

		fee = 50
	}

	fee += util.GetExtraFee(t.op.Raw, t.client.ServerInfo.DropsPerByte)
	tx.Fee = strconv.Itoa(fee)

	return tx, nil
}