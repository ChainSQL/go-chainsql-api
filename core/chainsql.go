package core

import (
	"log"

	"github.com/ChainSQL/go-chainsql-api/common"
	"github.com/ChainSQL/go-chainsql-api/export"
	"github.com/ChainSQL/go-chainsql-api/net"
	"github.com/ChainSQL/go-chainsql-api/util"
)

// Chainsql is the interface struct for this package
type Chainsql struct {
	client *net.Client
	SubmitBase
}

//TableListSetJSON specifies the table operation format
type TableListSetJSON struct {
	common.TableTxFields
	common.NetFields
	AutoFillField string
}

type TableGetSqlJSON struct {
	Account     string
	Sql         string
	LedgerIndex int
}

// NewChainsql is a function that create a chainsql object
func NewChainsql() *Chainsql {
	chainsql := &Chainsql{
		client: net.NewClient(),
	}
	chainsql.SubmitBase.client = chainsql.client
	chainsql.SubmitBase.IPrepare = chainsql
	return chainsql
}

// As specify the operating account
func (c *Chainsql) As(address string, secret string) {
	c.client.Auth.Address = address
	c.client.Auth.Secret = secret

	if c.client.Auth.Owner == "" {
		c.client.Auth.Owner = address
	}
}

// Use specify the table owner
func (c *Chainsql) Use(owner string) {
	c.client.Auth.Owner = owner
}

// PrepareTx prepare tx json for submit
func (c *Chainsql) PrepareTx() (TxJSON, error) {

	log.Println("Chainsql prepareTx")
	tx := &TableListSetJSON{}
	return tx, nil
}

//Table create a new table object
func (c *Chainsql) Table(name string) *Table {
	return NewTable(name, c.client)
}

//Connect is used to create a websocket connection
func (c *Chainsql) Connect(url string) error {
	return c.client.Connect(url)
}

// GetLedger request a ledger
func (c *Chainsql) GetLedger(seq int) string {
	return c.client.GetLedger(seq)
}

//OnLedgerClosed reponses in callback functor
func (c *Chainsql) OnLedgerClosed(callback export.Callback) {
	c.client.Event.SubscribeLedger(callback)
}

// GenerateAccount generate an account with the format:
// {
//		"address":"zxY4HEbEDSivZwouzwzqHQBA9QbJYdqDTg",
//		"publicKey":"cBPjenRgb2qzoYTnXmPV934kq5wpj2czHoz6rscHtzL34NqZN3KA",
//		"publicKeyHex":"02EA30B2A25844D4AFBAF6020DA9C9FED573AA0058791BFC8642E69888693CF8EA",
//		"privateKey":"xniMQKhxZTMbfWb8scjRPXa5Zv6HB",
// }
func (c *Chainsql) GenerateAccount(args ...string) (string, error) {
	if len(args) == 0 {
		return util.GenerateAccount()
	} else {
		return util.GenerateAccount(args[0])
	}
}

//SignPlainData sign a plain text and return the signature
func (c *Chainsql) SignPlainData(privateKey string, data string) (string, error) {
	return util.SignPlainData(privateKey, data)
}

//GetNameInDB request for table nameInDB
func (c *Chainsql) GetNameInDB(address string, tableName string) (string, error) {
	return c.client.GetNameInDB(address, tableName)
}

//GetBySqlUser is used to select from database by sql
func (c *Chainsql) GetBySqlUser(sql string) (string, error) {
	data := &TableGetSqlJSON{
		Account: c.client.Auth.Address,
		Sql:     sql,
	}
	if c.client.ServerInfo.Updated {
		data.LedgerIndex = c.client.ServerInfo.LedgerIndex
	} else {
		ledgerIndex, err := c.client.GetLedgerVersion()
		if err != nil {
			return "", err
		}
		data.LedgerIndex = ledgerIndex
	}
	return c.client.GetTableData(data, true)
}

func (c *Chainsql) IsConnected() bool {
	return c.client.GetWebocketManager().IsConnected()
}
