package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ChainSQL/go-chainsql-api/core"
	"github.com/buger/jsonparser"
	"github.com/gorilla/websocket"
)

// Account struct
type Account struct {
	address string
	secret  string
}

var user1 = Account{
	address: "zBonp9s7isAaDUPcfrFfYjNnhgeznoBHxF",
	secret:  "xn2FhQLRQqhKJeNhpgMzp2PGAYbdw",
}
var user2 = Account{
	address: "zKXfeKXkTtLSTkEzaJyu2cRmRBFRvTW2zc",
	secret:  "xhtBo8BLBZtTgc3LHnRspaFro5P4H",
}
var tableName = "hello2"

func main() {
	c := core.NewChainsql()
	//err := c.Connect("wss://zxlm-fgm.peersafe.cn/ws-zhu")
	//err := c.Connect("ws://10.100.0.78:25510")
	err := c.Connect("ws://localhost:5510")
	// log.Println("IsConnected:", c.IsConnected())
	// if err != nil {
	log.Println(err)
	// 	return
	// }
	// var root = Account{
	// 	address: "zHb9CJAWyB4zj91VRWn96DkukG4bwdtyTh",
	// 	secret:  "xnoPBzXtMeMyMHUVTgbuqAfg1SUTb",
	// }

	c.As(user1.address, user1.secret)
	// c.As(root.address, root.secret)
	// c.Use(root.address)

	// // testSubLedger(c)
	//testGenerateAccount(c)
	//testInsert(c)
	// testGetLedger(c)
	// testSignPlainText(c)

	// testGetTableData(c)

	// testGetBySqlUser(c)
	// testWebsocket()
	// testTickerGet(c)\
	//testValidationCreate(c)
	//	testGetAccountInfo(c)
	//testGetServerInfo(c)
	//testPay(c)
	//testSchemaCreate(c)  //创建子链
	//testSchemaModify(c)  // 修改子链
	//testGetSchemaList(c)  //获取子链列表
	//testGetSchemaInfo(c)  //依据子链id获取子链信息
	testStopSchema(c) //
	//testStartSchema(c)
	for {
		time.Sleep(time.Second * 10)
	}
}

func testGenerateAccount(c *core.Chainsql) {
	accStr, err := c.GenerateAccount()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(accStr)

	//Recreate account using the privateKey
	privateKey, err := jsonparser.GetString([]byte(accStr), "privateKey")
	accStr, err = c.GenerateAccount(string(privateKey))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(accStr)
}

func testInsert(c *core.Chainsql) {
	var data = []byte(`[{"id":1,"name":"echo","age":18}]`)
	ret := c.Table("gmTest50").Insert(string(data)).Submit("db_success")
	log.Println(ret)
}

func testGetTableData(c *core.Chainsql) {
	//Test withfields
	log.Println("IsConnected:", c.IsConnected())
	var counts = []byte(`[ "COUNT(*) as count" ]`)
	ret, err := c.Table(tableName).Get("").WithFields(string(counts)).Request()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(ret)
	count, err := jsonparser.GetInt([]byte(ret), "lines", "[0]", "count")
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Total record number:%d\n", count)

	//Test limit order
	type Limit struct {
		Total int `json:"total"`
		Index int `json:"index"`
	}
	var getRaw = []byte(`{"name":"echo"}`)
	var order = []byte(`[{"id":1},{"name":-1}]`)
	limit := Limit{
		Total: 10,
		Index: 0,
	}
	limitStr, err := json.Marshal(limit)
	ret, err = c.Table(tableName).Get(string(getRaw)).Limit(string(limitStr)).Order(string(order)).Request()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("The first 10 records:%s\n", ret)
}

func testGetBySqlUser(c *core.Chainsql) {
	nameInDB, err := c.GetNameInDB("zHb9CJAWyB4zj91VRWn96DkukG4bwdtyTh", "hello")
	if err != nil {
		log.Println(err)
	}
	sql := fmt.Sprintf("select * from t_%s limit 0,10", nameInDB)
	ret, err := c.GetBySqlUser(sql)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("The first 10 records:%s\n", ret)
	for {
		time.Sleep(time.Second * 10)
	}
}

func testGetLedger(c *core.Chainsql) {
	for i := 20; i < 25; i++ {
		ledger := c.GetLedger(i)
		log.Printf("GetLedger %d:%s\n", i, ledger)
	}
}

func testSubLedger(c *core.Chainsql) {
	go func() {
		c.OnLedgerClosed(func(msg string) {
			log.Printf("OnLedgerClosed:%s\n", msg)
		})
	}()
}

func testSignPlainText(c *core.Chainsql) {
	signed, err := c.SignPlainData("xnoPBzXtMeMyMHUVTgbuqAfg1SUTb", "HelloWorld")
	if err != nil {
		log.Println("error:", err)
		return
	}
	log.Printf("signature:%s\n", signed)
}

func testWebsocket() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://192.168.29.105:6006", nil)
	if err != nil {
		log.Println("error:", err)
		return
	}

	go func() {
		subViewChange := []byte(`{
			"command":"subscribe",
			"streams":["view_change","ledger"]
		}`)

		conn.WriteMessage(websocket.TextMessage, subViewChange)
	}()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv:%s", message)
		}
	}()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	time.Sleep(50 * time.Second)
}

func testTickerGet(c *core.Chainsql) {
	ticker := time.NewTicker(2000 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				go testGetTableData(c)
			}
		}
	}()
	time.Sleep(50000 * time.Millisecond)
	ticker.Stop()
	done <- true
}

func testValidationCreate(c *core.Chainsql) {
	seedKey, err := c.ValidationCreate()
	if err != nil {
		log.Println(err)
	}
	log.Printf("seedKey %s\n", seedKey)
}

func testGetAccountInfo(c *core.Chainsql) {
	account, err := c.GetAccountInfo(user1.address)
	if err != nil {
		log.Println(err)
	}
	log.Printf("seedKey %s\n", account)
}

func testGetServerInfo(c *core.Chainsql) {
	serverInfo, err := c.GetServerInfo()
	if err != nil {
		log.Println(err)
	}
	log.Printf("seedKey %s\n", serverInfo)
}

func testPay(c *core.Chainsql) {
	ret := c.Pay(user2.address, 30).Submit("validate_success")
	log.Println(ret)
}

func testSchemaCreate(c *core.Chainsql) {
	schemaInfo := "{\"SchemaName\":\"hello\",\"WithState\":false,\"SchemaAdmin\":\"zBonp9s7isAaDUPcfrFfYjNnhgeznoBHxF\",\"Validators\":[{\"Validator\":{\"PublicKey\":\"032AE5321413947612478BD8DF609ACCBB7EB07930404AE38F8FF48721D82C0D45\"}},{\"Validator\":{\"PublicKey\":\"0313FEF8C100B25A62E1428D0B414FD945B00F369E6735AD52AFBE3DB88D80AAE9\"}},{\"Validator\":{\"PublicKey\":\"030F9602B680A71D962111CC3EF9D27601EB9FEC7C3BA24BB4323D94C3A1CF9A04\"}},{\"Validator\":{\"PublicKey\":\"02221A8AA9228AD1199BE05AD64FC6A4625525FE8C9780FD68223C85168B58E738\"}}],\"PeerList\":[{\"Peer\":{\"Endpoint\":\"10.100.0.78:25410\"}},{\"Peer\":{\"Endpoint\":\"10.100.0.78:25411\"}},{\"Peer\":{\"Endpoint\":\"10.100.0.104:5410\"}},{\"Peer\":{\"Endpoint\":\"10.100.0.78:25412\"}}]}"
	//schemaInfo := "{\"SchemaName\":\"hello\",\"WithState\":false,\"SchemaAdmin\":\"zBonp9s7isAaDUPcfrFfYjNnhgeznoBHxF\",\"Validators\":\"fhfhhfhfhfhfh\",\"PeerList\":[{\"Peer\":{\"Endpoint\":\"127.0.0.1:15125\"}},{\"Peer\":{\"Endpoint\":\"127.0.0.1:25125\"}},{\"Peer\":{\"Endpoint\":\"127.0.0.1:35125\"}}]}"
	ret := c.CreateSchema(schemaInfo).Submit("validate_success")
	log.Println(ret)
}

func testSchemaModify(c *core.Chainsql) {
	schemaInfo := "{\"SchemaID\":\"6BA63B86E5CE48283D03CC21D3BE5F4630CC6572CE7F54982E5AE687C998B7A3\",\"Validators\":[{\"Validator\":{\"PublicKey\":\"02BD87A95F549ECF607D6AE3AEC4C95D0BFF0F49309B4E7A9F15B842EB62A8ED1B\"}}],\"PeerList\":[{\"Peer\":{\"Endpoint\":\"192.168.29.108:5125\"}}]}"
	//schemaInfo := "{\"SchemaName\":\"hello\",\"WithState\":false,\"SchemaAdmin\":\"zBonp9s7isAaDUPcfrFfYjNnhgeznoBHxF\",\"Validators\":\"fhfhhfhfhfhfh\",\"PeerList\":[{\"Peer\":{\"Endpoint\":\"127.0.0.1:15125\"}},{\"Peer\":{\"Endpoint\":\"127.0.0.1:25125\"}},{\"Peer\":{\"Endpoint\":\"127.0.0.1:35125\"}}]}"
	ret := c.ModifySchema("schema_add", schemaInfo).Submit("validate_success")
	log.Println(ret)
}

func testGetSchemaList(c *core.Chainsql) {
	param := ""
	ret, err := c.GetSchemaList(param)
	log.Println(ret)
	log.Println(err)
}

func testGetSchemaInfo(c *core.Chainsql) {
	schemaID := "68AAF6D84D4D2F18E3B00475475011F56A52B2877DD77B3190803F4FF9EB2F6E"
	ret, err := c.GetSchemaInfo(schemaID)
	log.Println(ret)
	log.Println(err)
}

func testStopSchema(c *core.Chainsql) {
	schemaID := "E3EEFAEAEDBFFEC22DF4BBC602E5D7DA0DEF9308F23A76ADAE03840A00141036"
	ret, err := c.StopSchema(schemaID)
	log.Println(ret)
	log.Println(err)
}

func testStartSchema(c *core.Chainsql) {
	schemaID := "E3EEFAEAEDBFFEC22DF4BBC602E5D7DA0DEF9308F23A76ADAE03840A00141036"
	ret, err := c.StartSchema(schemaID)
	log.Println(ret)
	log.Println(err)
}
