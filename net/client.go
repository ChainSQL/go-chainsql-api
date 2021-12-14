package net

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ChainSQL/go-chainsql-api/common"
	"github.com/ChainSQL/go-chainsql-api/crypto"
	"github.com/ChainSQL/go-chainsql-api/event"
	"github.com/ChainSQL/go-chainsql-api/export"
	"github.com/ChainSQL/go-chainsql-api/util"

	"github.com/buger/jsonparser"
)

//ReconnectInterval is the interval to reconnect when ws socket is disconnected
const ReconnectInterval = 10

// Client is used to send and recv websocket msg
type Client struct {
	cmdIDs      int64
	SchemaID    string
	wm          *WebsocketManager
	sendMsgChan chan string
	recvMsgChan chan string
	requests    map[int64]*Request
	mutex       *sync.RWMutex
	Auth        *common.Auth
	ServerInfo  *ServerInfo
	Event       *event.Manager
	inited      bool
}

//NewClient is constructor
func NewClient() *Client {
	return &Client{
		cmdIDs:     0,
		requests:   make(map[int64]*Request),
		mutex:      new(sync.RWMutex),
		Auth:       &common.Auth{},
		ServerInfo: NewServerInfo(),
		Event:      event.NewEventManager(),
		inited:     false,
		SchemaID:   "",
	}
}

//Connect is used to create a websocket connection
func (c *Client) Connect(url string) error {
	if c.wm != nil {
		return c.reConnect(url)
	}

	c.wm = NewWsClientManager(url, ReconnectInterval)
	err := c.wm.Start()
	if err != nil {
		return err
	}

	c.init()
	return nil
}

func (c *Client) reConnect(url string) error {
	err := c.wm.Disconnect()
	if err != nil {
		return err
	}
	c.wm.SetUrl(url)
	err = c.wm.Start()
	if err != nil {
		return err
	}
	if !c.inited {
		c.init()
	} else {
		//connect changed,only subscribe
		c.InitSubscription()
	}
	return nil
}

func (c *Client) init() {
	c.sendMsgChan = c.wm.WriteChan()
	c.recvMsgChan = c.wm.ReadChan()

	go c.processMessage()
	go c.checkReconnection()
	c.InitSubscription()
	c.inited = true
}

func (c *Client) checkReconnection() {
	c.wm.OnReconnected(func() {
		c.InitSubscription()
	})
}

func (c *Client) GetWebocketManager() *WebsocketManager {
	return c.wm
}

func (c *Client) InitSubscription() {
	type Subscribe struct {
		common.RequestBase
		Streams []string `json:"streams"`
	}
	c.cmdIDs++
	subCmd := &Subscribe{
		RequestBase: common.RequestBase{
			Command: "subscribe",
			ID:      c.cmdIDs,
		},
		Streams: []string{"ledger", "server"},
	}
	if c.SchemaID != "" {
		subCmd.SchemaID = c.SchemaID
	}
	request := c.syncRequest(subCmd)

	result, _, _, err := jsonparser.Get([]byte(request.Response.Value), "result")
	if err != nil {
		fmt.Printf("initSubscription error:%s\n", err)
		return
	}
	c.ServerInfo.Update(string(result))
}

func (c *Client) processMessage() {
	for msg := range c.recvMsgChan {
		go c.handleClientMsg(msg)
	}
}

func (c *Client) handleClientMsg(msg string) {
	// log.Printf("handleClientMsg: %s", msg)
	msgType, err := jsonparser.GetString([]byte(msg), "type")
	if err != nil {
		fmt.Printf("handleClientMsg error:%s\n", err)
	}
	// fmt.Println(msgType)

	switch msgType {
	case "response":
		c.onResponse(msg)
	case "serverStatus":
		c.ServerInfo.Update(msg)
	case "ledgerClosed":
		c.ServerInfo.Update(msg)
		c.onLedgerClosed(msg)
	case "singleTransaction":
		c.onSingleTransaction(msg)
	case "table":
		c.onTableMsg(msg)
	default:
		log.Printf("Unhandled message %s", msg)
	}
}

func (c *Client) onResponse(msg string) {
	id, err := jsonparser.GetInt([]byte(msg), "id")
	if err != nil {
		// fmt.Println(err)
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	request, ok := c.requests[id]
	if !ok {
		log.Printf("onResponse:Request with id %d not exist\n", id)
		return
	}

	defer request.Wait.Done()
	delete(c.requests, id)

	request.Response = &Response{
		Value:   msg,
		Request: request,
	}
}

func (c *Client) onLedgerClosed(msg string) {
	c.Event.OnLedgerClosed(msg)
}

func (c *Client) onSingleTransaction(msg string) {
	c.Event.OnSingleTransaction(msg)
}

func (c *Client) onTableMsg(msg string) {
	c.Event.OnTableMsg(msg)
}

// GetLedger request for ledger data
func (c *Client) GetLedger(seq int) string {
	type getLedger struct {
		common.RequestBase
		LedgerIndex int `json:"ledger_index"`
	}
	c.cmdIDs++
	ledgerReq := &getLedger{
		RequestBase: common.RequestBase{
			Command: "ledger",
			ID:      c.cmdIDs,
		},
		LedgerIndex: seq,
	}
	request := c.syncRequest(ledgerReq)

	return request.Response.Value
}

// GetLedgerVersion request for ledger version
func (c *Client) GetLedgerVersion() (int, error) {
	type Request struct {
		common.RequestBase
	}
	c.cmdIDs++
	ledgerReq := &Request{
		RequestBase: common.RequestBase{
			Command: "ledger_current",
			ID:      c.cmdIDs,
		},
	}
	request := c.syncRequest(ledgerReq)
	err := c.parseResponseError(request)
	if err != nil {
		log.Println("GetLedgerVersion:", err)
		return 0, err
	}
	ledgerIndex, err := jsonparser.GetInt([]byte(request.Response.Value), "result", "ledger_current_index")
	if err != nil {
		return 0, err
	}
	return int(ledgerIndex), nil
}

func (c *Client) parseResponseError(request *Request) error {
	status, err := jsonparser.GetString([]byte(request.Response.Value), "status")
	if err != nil {
		return err
	}
	if status == "error" {
		errMsg, _ := jsonparser.GetString([]byte(request.Response.Value), "error_message")
		return fmt.Errorf("%s", errMsg)
	}
	return nil
}

// GetAccountInfo request for account_info
func (c *Client) GetAccountInfo(address string) (string, error) {
	type getAccount struct {
		common.RequestBase
		Account string `json:"account"`
	}
	c.cmdIDs++
	accountReq := &getAccount{}
	accountReq.ID = c.cmdIDs
	accountReq.Command = "account_info"
	accountReq.Account = address

	request := c.syncRequest(accountReq)

	err := c.parseResponseError(request)
	if err != nil {
		return "", err
	}

	return request.Response.Value, nil
}

// GetNameInDB request for table nameInDB
func (c *Client) GetNameInDB(address string, tableName string) (string, error) {
	type Request struct {
		common.RequestBase
		Account   string `json:"account"`
		TableName string `json:"tablename"`
	}
	c.cmdIDs++
	req := &Request{}
	req.ID = c.cmdIDs
	req.Command = "g_dbname"
	req.Account = address
	req.TableName = tableName

	request := c.syncRequest(req)

	err := c.parseResponseError(request)
	if err != nil {
		return "", err
	}
	nameInDB, err := jsonparser.GetString([]byte(request.Response.Value), "result", "nameInDB")
	if err != nil {
		return "", err
	}
	return nameInDB, nil
}

//Submit submit a signed transaction
func (c *Client) Submit(blob string) string {
	type Request struct {
		common.RequestBase
		TxBlob string `json:"tx_blob"`
	}
	c.cmdIDs++
	req := &Request{}
	req.ID = c.cmdIDs
	req.Command = "submit"
	req.TxBlob = blob
	if c.SchemaID != "" {
		req.SchemaID = c.SchemaID
	}

	request := c.syncRequest(req)

	return request.Response.Value
}

//SubscribeTx subscribe a transaction by hash
func (c *Client) SubscribeTx(hash string, callback export.Callback) {
	c.Event.SubscribeTx(hash, callback)

	type Request struct {
		common.RequestBase
		TxHash string `json:"transaction"`
	}
	req := Request{}
	req.Command = "subscribe"
	req.TxHash = hash
	if c.SchemaID != "" {
		req.SchemaID = c.SchemaID
	}
	c.asyncRequest(req)
}

//UnSubscribeTx subscribe a transaction by hash
func (c *Client) UnSubscribeTx(hash string) {
	c.Event.UnSubscribeTx(hash)

	type Request struct {
		common.RequestBase
		TxHash string `json:"transaction"`
	}
	req := Request{}
	req.Command = "unsubscribe"
	req.TxHash = hash
	if c.SchemaID != "" {
		req.SchemaID = c.SchemaID
	}
	c.asyncRequest(req)
}

func (c *Client) GetTableData(dataJSON interface{}, bSql bool) (string, error) {
	type Request struct {
		common.RequestBase
		PublicKey   string      `json:"publicKey"`
		Signature   string      `json:"signature"`
		SigningData string      `json:"signingData"`
		TxJSON      interface{} `json:"tx_json"`
	}
	c.cmdIDs++
	req := &Request{}
	req.ID = c.cmdIDs
	req.Command = "r_get"
	if bSql {
		req.Command = "r_get_sql_user"
	}
	req.TxJSON = dataJSON
	accStr, err := crypto.GenerateAccount(c.Auth.Secret)
	if err != nil {
		return "", err
	}
	publicKey, err := jsonparser.GetString([]byte(accStr), "publicKeyHex")
	if err != nil {
		return "", err
	}
	jsonStr, err := json.Marshal(dataJSON)
	if err != nil {
		return "", err
	}
	signature, err := util.SignPlainData(c.Auth.Secret, string(jsonStr))
	if err != nil {
		return "", err
	}

	req.PublicKey = publicKey
	req.SigningData = string(jsonStr)
	req.Signature = string(signature)

	request := c.syncRequest(req)

	err = c.parseResponseError(request)
	if err != nil {
		if err.Error() == "Invalid field 'LedgerIndex'." {
			c.ServerInfo.Updated = false
		}
		return "", err
	}

	result, _, _, err := jsonparser.Get([]byte(request.Response.Value), "result")
	if err != nil {
		// log.Printf("Cleint::GetTableData %s\n",err)
		return "", err
	}
	// log.Printf("type:%T\n",result)
	return string(result), nil
}

func (c *Client) syncRequest(v common.IRequest) *Request {
	if c.SchemaID != "" {
		v.SetSchemaID(c.SchemaID)
	}
	data, _ := json.Marshal(v)
	request := NewRequest(v.GetID(), string(data))
	request.Wait.Add(1)
	c.sendRequest(request)
	done := make(chan struct{})
	go func() {
		request.Wait.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(util.REQUEST_TIMEOUT * time.Second):
		{
			timeOutMsg := string(`{
				"status":"error",
				"error_message":"request timeout"
			}`)
			request.Response = &Response{
				Value:   timeOutMsg,
				Request: request,
			}
		}
	}
	return request
}

func (c *Client) sendRequest(request *Request) {
	c.mutex.Lock()
	c.requests[request.ID] = request
	c.mutex.Unlock()

	// log.Printf("sendRequest %s\n", request.JSON)
	c.sendMsgChan <- request.JSON
}

func (c *Client) asyncRequest(v interface{}) {
	data, _ := json.Marshal(v)
	c.sendMsgChan <- string(data)
}

// GetServerInfo request for ServerInfo
func (c *Client) GetServerInfo() (string, error) {
	type getServerInfo struct {
		common.RequestBase
	}
	c.cmdIDs++
	accountReq := &getServerInfo{}
	accountReq.ID = c.cmdIDs
	accountReq.Command = "server_info"

	request := c.syncRequest(accountReq)

	err := c.parseResponseError(request)
	if err != nil {
		return "", err
	}

	return request.Response.Value, nil
}

func (c *Client) GetSchemaList(params string) (string, error) {
	account, _ := jsonparser.GetString([]byte(params), "account")
	running, _ := jsonparser.GetBoolean([]byte(params), "running")
	type getSchemaList struct {
		common.RequestBase
		Account string `json:"account,omitempty"`
		Running bool   `json:"running,omitempty"`
	}
	c.cmdIDs++
	schemaListReq := &getSchemaList{}
	schemaListReq.ID = c.cmdIDs
	schemaListReq.Command = "schema_list"
	if account != "" {
		schemaListReq.Account = account
	}
	if strings.Contains(params, "running") {
		schemaListReq.Running = running
	}
	request := c.syncRequest(schemaListReq)
	err := c.parseResponseError(request)
	if err != nil {
		return "", err
	}
	return request.Response.Value, nil
}
func (c *Client) GetSchemaInfo(schemaID string) (string, error) {
	if schemaID == "" {
		panic("Invalid parameter")
	}
	type getSchemaInfo struct {
		common.RequestBase
		Schema string `json:"schema"`
	}
	c.cmdIDs++
	schemaInfoReq := &getSchemaInfo{}
	schemaInfoReq.ID = c.cmdIDs
	schemaInfoReq.Command = "schema_info"
	schemaInfoReq.Schema = schemaID
	request := c.syncRequest(schemaInfoReq)

	err := c.parseResponseError(request)
	if err != nil {
		return "", err
	}
	return request.Response.Value, nil
}

func (c *Client) StopSchema(schemaID string) (string, error) {
	if schemaID == "" {
		panic("Invalid parameter")
	}
	type StopSchema struct {
		common.RequestBase
		Schema string `json:"schema"`
	}
	c.cmdIDs++
	StopSchemaReq := &StopSchema{}
	StopSchemaReq.ID = c.cmdIDs
	StopSchemaReq.Command = "stop"
	StopSchemaReq.Schema = schemaID
	request := c.syncRequest(StopSchemaReq)

	err := c.parseResponseError(request)
	if err != nil {
		return "", err
	}
	return request.Response.Value, nil
}

func (c *Client) StartSchema(schemaID string) (string, error) {
	if schemaID == "" {
		panic("Invalid parameter")
	}
	type StartSchema struct {
		common.RequestBase
		Schema string `json:"schema"`
	}
	c.cmdIDs++
	StartSchemaReq := &StartSchema{}
	StartSchemaReq.ID = c.cmdIDs
	StartSchemaReq.Command = "schema_start"
	StartSchemaReq.Schema = schemaID
	request := c.syncRequest(StartSchemaReq)

	err := c.parseResponseError(request)
	if err != nil {
		return "", err
	}
	return request.Response.Value, nil
}

func (c *Client) Unsubscribe(schemaID string) (string, error) {
	type Unsubscribe struct {
		common.RequestBase
		Streams []string `json:"streams"`
	}
	c.cmdIDs++
	unsubCmd := &Unsubscribe{
		RequestBase: common.RequestBase{
			Command: "unsubscribe",
			ID:      c.cmdIDs,
		},
		Streams: []string{"ledger", "server"},
	}
	request := c.syncRequest(unsubCmd)

	err := c.parseResponseError(request)
	if err != nil {
		return "", err
	}
	return request.Response.Value, nil
}

func (c *Client) GetTransaction(hash string) (string, error) {
	type getTransaction struct {
		common.RequestBase
		Transaction string `json:"transaction"`
	}

	c.cmdIDs++
	getTransactionReq := &getTransaction{}
	getTransactionReq.ID = c.cmdIDs
	getTransactionReq.Command = "tx"
	getTransactionReq.Transaction = hash
	request := c.syncRequest(getTransactionReq)

	err := c.parseResponseError(request)
	if err != nil {
		return "", err
	}
	return request.Response.Value, nil
}
