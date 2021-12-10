package core

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"

	"github.com/ChainSQL/go-chainsql-api/crypto"
	. "github.com/ChainSQL/go-chainsql-api/data"
	"github.com/ChainSQL/go-chainsql-api/export"
	"github.com/ChainSQL/go-chainsql-api/net"
	"github.com/ChainSQL/go-chainsql-api/util"
	"github.com/buger/jsonparser"
)

// Chainsql is the interface struct for this package
type Chainsql struct {
	client *net.Client
	SubmitBase
	op *ChainsqlTxInfo
}

//TxInfo is the opearting details
type ChainsqlTxInfo struct {
	Signer Signer
	Query  []interface{}
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
		op: &ChainsqlTxInfo{
			Query: make([]interface{}, 0),
		},
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
func (c *Chainsql) PrepareTx() (Signer, error) {

	log.Println("Chainsql prepareTx")
	tx := c.op.Signer
	seq, err := net.PrepareRipple(c.client)
	if err != nil {
		log.Println(err)
	}

	var fee int64 = 10
	var last uint32
	if c.client.ServerInfo.Updated {
		last = uint32(c.client.ServerInfo.LedgerIndex + 20)
		fee = int64(c.client.ServerInfo.ComputeFee())
	} else {
		ledgerIndex, err := c.client.GetLedgerVersion()
		if err != nil {
			log.Println("Chainsql prepareTx ", err)
		}
		last = uint32(ledgerIndex + 20)
		fee = 50
	}

	if c.op.Signer.GetRaw() != "" {
		fee += util.GetExtraFee(c.op.Signer.GetRaw(), c.client.ServerInfo.DropsPerByte)
	} else if c.op.Signer.GetStatements() != "" {
		fee += util.GetExtraFee(c.op.Signer.GetStatements(), c.client.ServerInfo.DropsPerByte)
	}

	finalFee, err := NewNativeValue(fee)
	if err != nil {
		log.Println("Chainsql prepareTx", err)
	}
	account, err := NewAccountFromAddress(c.client.Auth.Address)
	if err != nil {
		log.Println("Chainsql prepareTx", err)
		panic(err)
	}
	c.op.Signer.SetTxBase(seq, *finalFee, &last, *account)
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
		return crypto.GenerateAccount()
	} else {
		return crypto.GenerateAccount(args[0])
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
	if c.client.GetWebocketManager() != nil {
		return c.client.GetWebocketManager().IsConnected()
	}
	return false
}

func (c *Chainsql) Disconnect() {
	if c.client.GetWebocketManager() != nil {
		c.client.GetWebocketManager().Disconnect()
	}
}

func (c *Chainsql) ValidationCreate() (string, error) {
	return crypto.ValidationCreate()
}

func (c *Chainsql) GetServerInfo() (string, error) {
	return c.client.GetServerInfo()
}

func (c *Chainsql) GetAccountInfo(address string) (string, error) {
	return c.client.GetAccountInfo(address)
}

func (c *Chainsql) Pay(accountId string, value int64) *Ripple {
	r := NewRipple(c.client)
	return r.Pay(accountId, value)
}

func (c *Chainsql) CreateSchema(schemaInfo string) *Chainsql {
	StrContainers := strings.Contains(schemaInfo, "SchemaName") && strings.Contains(schemaInfo, "WithState") &&
		strings.Contains(schemaInfo, "Validators") && strings.Contains(schemaInfo, "PeerList")

	if !StrContainers {
		panic("Invalid schemaInfo parameter")
	}
	createSchema := &SchemaCreate{TxBase: TxBase{TransactionType: SCHEMA_CREATE}}
	var jsonObj CreateSchema
	err := json.Unmarshal([]byte(schemaInfo), &jsonObj)
	if err != nil {
		log.Println("CreateSchema ", err)
		panic(err)
	}

	createSchema.SchemaName = VariableLength(jsonObj.SchemaName)
	account, err := NewAccountFromAddress(jsonObj.SchemaAdmin)
	if err != nil {
		log.Println(err)
	}
	if account != nil {
		createSchema.SchemaAdmin = account
	}
	if jsonObj.WithState {
		//继承主链的节点状态
		leadgerHash, errHash := NewHash256(jsonObj.AnchorLedgerHash)
		if errHash != nil {
			log.Println("CreateSchema ", errHash)
		}
		if leadgerHash != nil {
			createSchema.AnchorLedgerHash = leadgerHash
		}
		createSchema.SchemaStrategy = 2
	} else {
		// 不继承主链的节点状态
		createSchema.SchemaStrategy = 1
		if strings.Contains(schemaInfo, "AnchorLedgerHash") {
			panic("Field 'AnchorLedgerHash' is unnecessary")
		}
	}

	validatorSlice := make([]ValidatorFormat, len(jsonObj.Validators))
	for i := 0; i < len(jsonObj.Validators); i++ {
		publicKeyHex := jsonObj.Validators[i].Validator.PublicKey
		publicKey, _ := hex.DecodeString(publicKeyHex)
		validatorSlice[i].Validator.PublicKey = VariableLength(publicKey)
	}
	createSchema.Validators = validatorSlice
	peerSlice := make([]PeerFormat, len(jsonObj.PeerList))
	for i := 0; i < len(jsonObj.PeerList); i++ {
		endpoint := jsonObj.PeerList[i].Peer.Endpoint
		peerSlice[i].Peer.Endpoint = VariableLength(endpoint)
	}

	createSchema.PeerList = peerSlice
	//createSchema.TransactionType = SCHEMA_CREATE
	c.op.Signer = createSchema
	return c
}

func (c *Chainsql) ModifySchema(schemaType string, schemaInfo string) *Chainsql {
	StrContainers := strings.Contains(schemaInfo, "SchemaID") && strings.Contains(schemaInfo, "Validators") && strings.Contains(schemaInfo, "PeerList")

	if !StrContainers {
		panic("Invalid schemaInfo parameter")
	}
	var jsonObj ModifySchema
	err := json.Unmarshal([]byte(schemaInfo), &jsonObj)
	if err != nil {
		log.Println("CreateSchema ", err)
		panic(err)
	}
	schemaModify := &SchemaModify{TxBase: TxBase{TransactionType: SCHEMA_MODIFY}}
	if schemaType == util.SchemaDel {
		schemaModify.OpType = 2
	} else {
		schemaModify.OpType = 1
	}
	validatorSlice := make([]ValidatorFormat, len(jsonObj.Validators))
	for i := 0; i < len(jsonObj.Validators); i++ {
		publicKeyHex := jsonObj.Validators[i].Validator.PublicKey
		publicKey, _ := hex.DecodeString(publicKeyHex)
		validatorSlice[i].Validator.PublicKey = VariableLength(publicKey)
	}
	schemaModify.Validators = validatorSlice

	peerSlice := make([]PeerFormat, len(jsonObj.PeerList))
	for i := 0; i < len(jsonObj.PeerList); i++ {
		endpoint := jsonObj.PeerList[i].Peer.Endpoint
		peerSlice[i].Peer.Endpoint = VariableLength(endpoint)
	}
	schemaModify.PeerList = peerSlice
	schemaIdHash, errHash := NewHash256(jsonObj.SchemaID)
	if errHash != nil {
		log.Println("CreateSchema ", errHash)
	}
	if schemaIdHash != nil {
		schemaModify.SchemaID = *schemaIdHash
	}
	//schemaModify.TransactionType = SCHEMA_MODIFY
	c.op.Signer = schemaModify
	return c
}

func (c *Chainsql) GetSchemaList(params string) (string, error) {
	return c.client.GetSchemaList(params)
}

func (c *Chainsql) GetSchemaInfo(schemaID string) (string, error) {
	return c.client.GetSchemaInfo(schemaID)
}

func (c *Chainsql) StopSchema(schemaID string) (string, error) {
	return c.client.StopSchema(schemaID)
}

func (c *Chainsql) StartSchema(schemaID string) (string, error) {
	return c.client.StartSchema(schemaID)
}

// func (c *Chainsql) setSchema(schemaId string) {
// 	if c.client.SchemaID == schemaId {

// 	}
// }
func (c *Chainsql) GetSchemaId(hash string) (string, error) {
	response, _ := c.client.GetTransaction(hash)
	LedgerEntryType, err := jsonparser.GetString([]byte(response), "result", "meta", "AffectedNodes", "[0]", "CreatedNode", "LedgerEntryType")
	if err != nil {
		return "", err
	}
	if LedgerEntryType == "Schema" {
		schemaID, err := jsonparser.GetString([]byte(response), "result", "meta", "AffectedNodes", "[0]", "CreatedNode", "LedgerIndex")
		if err != nil {
			return "", err
		}
		return schemaID, nil
	}
	panic("Invalid parameter")
}
func (c *Chainsql) GetTransaction(hash string) (string, error) {
	return c.client.GetTransaction(hash)
}
