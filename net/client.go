package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/ChainSQL/go-chainsql-api/common"
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
	schemaID    string
	wm          *WebsocketManager
	sendMsgChan chan string
	recvMsgChan chan string
	requests    map[int64]*Request
	mutex       *sync.RWMutex
	Auth        *common.Auth
	ServerInfo  *ServerInfo
	Event       *event.Manager
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
	}
}

//Connect is used to create a websocket connection
func (c *Client) Connect(url string) error {
	c.wm = NewWsClientManager(url, ReconnectInterval)
	err := c.wm.Start()
	if err != nil {
		return err
	}
	c.sendMsgChan = c.wm.WriteChan()
	c.recvMsgChan = c.wm.ReadChan()

	go c.processMessage()
	c.initSubscription()
	return nil
}

func (c *Client) initSubscription() {
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
	request, ok := c.requests[id]
	if !ok {
		log.Printf("onResponse:Request with id %d not exist\n", id)
		return
	}
	defer request.Wait.Done()
	c.mutex.Lock()
	delete(c.requests, id)
	c.mutex.Unlock()
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

	ledgerIndex, err := jsonparser.GetInt([]byte(request.Response.Value), "result", "ledger_current_index")
	if err != nil {
		return 0, err
	}
	return int(ledgerIndex), nil
}

// GetAccountInfo request for account_info
func (c *Client) GetAccountInfo(address string) string {
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

	return request.Response.Value
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
	status, err := jsonparser.GetString([]byte(request.Response.Value), "status")
	if err != nil {
		return "", err
	}
	if status == "error" {
		errCode, _ := jsonparser.GetString([]byte(request.Response.Value), "error")
		return "", fmt.Errorf("%s", errCode)
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
	accStr, err := util.GenerateAccount(c.Auth.Secret)
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

	status, err := jsonparser.GetString([]byte(request.Response.Value), "status")
	if err != nil {
		return "", err
	}
	if status != "success" {
		errorMsg, _ := jsonparser.GetString([]byte(request.Response.Value), "error_message")
		return "", errors.New(errorMsg)
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
	data, _ := json.Marshal(v)
	request := NewRequest(v.GetID(), string(data))
	request.Wait.Add(1)
	c.sendRequest(request)
	request.Wait.Wait()
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
