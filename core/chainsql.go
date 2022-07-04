package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	// client *net.Client
	SubmitBase
	op *ChainsqlTxInfo
}

//TxInfo is the opearting details
type ChainsqlTxInfo struct {
	//Signer Signer
	Raw    string
	TxType TransactionType
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
		// client: net.NewClient(),
		op: &ChainsqlTxInfo{
			Query: make([]interface{}, 0),
		},
	}
	chainsql.client = net.NewClient()
	// chainsql.SubmitBase.client = chainsql.client
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

//Table create a new table object
func (c *Chainsql) Table(name string) *Table {
	return NewTable(name, c.client)
}

//Connect is used to create a websocket connection
func (c *Chainsql) Connect(url, tlsRootCertPath, tlsClientCertPath, tlsClientKeyPath, serverName string) error {
	return c.client.Connect(url, tlsRootCertPath, tlsClientCertPath, tlsClientKeyPath, serverName)
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
/*func (c *Chainsql) GenerateAccount(args ...string) (string, error) {
	if len(args) == 0 {
		return crypto.GenerateAccount()
	} else {
		return crypto.GenerateAccount(args[0])
	}
}*/

func (c *Chainsql) GenerateAddress(options string) (string, error) {
	return crypto.GenerateAddress(options)
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
		c.client.Unsubscribe()
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
	c.op.TxType = SCHEMA_CREATE
	c.op.Raw = schemaInfo
	return c
}

func (c *Chainsql) createSchema() (Signer, error) {
	var schemaInfo = c.op.Raw
	isValid := strings.Contains(schemaInfo, "SchemaName") && strings.Contains(schemaInfo, "WithState") &&
		strings.Contains(schemaInfo, "Validators") && strings.Contains(schemaInfo, "PeerList")

	if !isValid {
		return nil, fmt.Errorf("Invalid schemaInfo parameter")
	}
	createSchema := &SchemaCreate{TxBase: TxBase{TransactionType: SCHEMA_CREATE}}
	var jsonObj CreateSchema
	err := json.Unmarshal([]byte(schemaInfo), &jsonObj)
	if err != nil {
		return nil, err
	}

	createSchema.SchemaName = VariableLength(jsonObj.SchemaName)
	if strings.Contains(schemaInfo, "SchemaAdmin") {
		account, err := NewAccountFromAddress(jsonObj.SchemaAdmin)
		if err != nil {
			return nil, fmt.Errorf("Invalid schemaInfo parameter: SchemaAdmin")
		}
		if account != nil {
			createSchema.SchemaAdmin = account
		}
	}

	if jsonObj.WithState {
		//继承主链的节点状态
		if strings.Contains(schemaInfo, "AnchorLedgerHash") {
			leadgerHash, errHash := NewHash256(jsonObj.AnchorLedgerHash)
			if errHash != nil {
				return nil, fmt.Errorf("Invalid schemaInfo parameter: AnchorLedgerHash")
			}
			if leadgerHash != nil {
				createSchema.AnchorLedgerHash = leadgerHash
			}
		}
		createSchema.SchemaStrategy = 2
	} else {
		// 不继承主链的节点状态
		createSchema.SchemaStrategy = 1
		if strings.Contains(schemaInfo, "AnchorLedgerHash") {
			return nil, fmt.Errorf("Field 'AnchorLedgerHash' is unnecessary")
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
	var signer Signer = createSchema
	return signer, nil
}
func (c *Chainsql) ModifySchema(schemaType string, schemaInfo string) *Chainsql {
	c.op.TxType = SCHEMA_MODIFY
	c.op.Raw = "{\"SchemaType\": \"" + schemaType + "\", \"SchemaInfo\":" + schemaInfo + "}"
	return c
}

func (c *Chainsql) modifySchema() (Signer, error) {
	schemaType, _ := jsonparser.GetString([]byte(c.op.Raw), "SchemaType")
	result, _, _, _ := jsonparser.Get([]byte(c.op.Raw), "SchemaInfo")
	schemaInfo := string(result)
	isValid := strings.Contains(schemaInfo, "SchemaID") && strings.Contains(schemaInfo, "Validators") && strings.Contains(schemaInfo, "PeerList")

	if !isValid {
		return nil, fmt.Errorf("Invalid schemaInfo parameter")
	}
	var jsonObj ModifySchema
	errUnmarshal := json.Unmarshal([]byte(schemaInfo), &jsonObj)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	schemaModify := &SchemaModify{TxBase: TxBase{TransactionType: SCHEMA_MODIFY}}
	if schemaType == util.SchemaDel {
		schemaModify.OpType = util.OpTypeSchemaDel
	} else {
		schemaModify.OpType = util.OpTypeSchemaAdd
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
		return nil, errHash
	}
	if schemaIdHash != nil {
		schemaModify.SchemaID = *schemaIdHash
	}
	//schemaModify.TransactionType = SCHEMA_MODIFY
	var signer Signer = schemaModify
	return signer, nil
}

func (c *Chainsql) DeleteSchema(schemaID string) *Chainsql {
	c.op.TxType = SCHEMA_DELETE
	c.op.Raw = schemaID
	return c
}

func (c *Chainsql) deleteSchema() (Signer, error) {
	var schemaID = c.op.Raw
	if schemaID == "" {
		return nil, fmt.Errorf("Invalid parameter")
	}
	schemaDelete := &SchemaDelete{TxBase: TxBase{TransactionType: SCHEMA_DELETE}}
	schemaIdHash, errHash := NewHash256(schemaID)
	if errHash != nil {
		return nil, errHash
	}
	if schemaIdHash != nil {
		schemaDelete.SchemaID = *schemaIdHash
	}
	var signer Signer = schemaDelete
	return signer, nil
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

func (c *Chainsql) SetSchema(schemaId string) {
	if c.client.SchemaID != schemaId {
		c.client.Unsubscribe()
		c.client.SchemaID = schemaId
		c.client.InitSubscription()
	}
}
func (c *Chainsql) GetSchemaId(hash string) (string, error) {
	response, _ := c.client.GetTransaction(hash)
	if response == "" {
		return "", fmt.Errorf("Transaction does not exist ")
	}
	schemaID := ""
	flag := false
	//LedgerEntryType, err := jsonparser.GetString([]byte(response), "result", "meta", "AffectedNodes", "[0]", "CreatedNode", "LedgerEntryType")
	jsonparser.ArrayEach([]byte(response), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		LedgerEntryType, err := jsonparser.GetString(value, "CreatedNode", "LedgerEntryType")
		if err == nil {
			if LedgerEntryType == "Schema" {
				schemaID, _ = jsonparser.GetString([]byte(value), "CreatedNode", "LedgerIndex")
				flag = true
			}
		}

	}, "result", "meta", "AffectedNodes")
	if flag {
		return schemaID, nil
	}
	return "", fmt.Errorf("Invalid parameter")
}
func (c *Chainsql) GetTransaction(hash string) (string, error) {
	return c.client.GetTransaction(hash)
}

func (c *Chainsql) GetTransactionResult(hash string) (string, error) {
	return c.client.GetTransactionResult(hash)
}

// PrepareTx prepare tx json for submit
func (c *Chainsql) PrepareTx() (Signer, error) {
	var tx Signer
	var err error
	switch c.op.TxType {
	case SCHEMA_CREATE:
		tx, err = c.createSchema()
		break
	case SCHEMA_MODIFY:
		tx, err = c.modifySchema()
		break
	case SCHEMA_DELETE:
		tx, err = c.deleteSchema()
		break
	default:
	}
	if err != nil {
		return nil, err
	}

	return c.prepareTxBase(tx)
}

func (c *Chainsql) prepareTxBase(tx Signer) (Signer, error) {

	//tx := c.op.Signer
	seq, err := net.PrepareRipple(c.client)
	if err != nil {
		return nil, err
	}

	var fee int64 = 10
	var last uint32
	if c.client.ServerInfo.Updated {
		last = uint32(c.client.ServerInfo.LedgerIndex + util.Seqinterval)
		fee = int64(c.client.ServerInfo.ComputeFee())
	} else {
		ledgerIndex, err := c.client.GetLedgerVersion()
		if err != nil {
			return nil, err
		}
		last = uint32(ledgerIndex + util.Seqinterval)
		fee = 50
	}

	if tx.GetRaw() != "" {
		fee += util.GetExtraFee(tx.GetRaw(), c.client.ServerInfo.DropsPerByte)
	} else if tx.GetStatements() != "" {
		fee += util.GetExtraFee(tx.GetStatements(), c.client.ServerInfo.DropsPerByte)
	}

	finalFee, err := NewNativeValue(fee)
	if err != nil {
		return nil, err
	}
	account, err := NewAccountFromAddress(c.client.Auth.Address)
	if err != nil {
		return nil, err
	}
	tx.SetTxBase(seq, *finalFee, &last, *account)
	return tx, nil
}
