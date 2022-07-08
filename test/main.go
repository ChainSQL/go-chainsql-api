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

var root1 = Account{
	address: "zPHxWCKNZjpbQHV5DLpBy8rSR8HdXpPDzi",
	secret:  "xp6FwxZP1rrmPy2GDTobvHTgnZnrC",
}

//zPHxWCKNZjpbQHV5DLpBy8rSR8HdXpPDzi
//zEdmwPS5BoEukn1ech2LHjuV6STkq3QYkM
var root = Account{
	address: "zHb9CJAWyB4zj91VRWn96DkukG4bwdtyTh",
	secret:  "xnoPBzXtMeMyMHUVTgbuqAfg1SUTb",
}
var user1 = Account{
	address: "zBonp9s7isAaDUPcfrFfYjNnhgeznoBHxF",
	secret:  "xn2FhQLRQqhKJeNhpgMzp2PGAYbdw",
}
var user2 = Account{
	address: "zKXfeKXkTtLSTkEzaJyu2cRmRBFRvTW2zc",
	secret:  "xhtBo8BLBZtTgc3LHnRspaFro5P4H",
}
var smRoot = Account{
	address: "zN7TwUjJ899xcvNXZkNJ8eFFv2VLKdESsj",
	secret:  "p97evg5Rht7ZB7DbEpVqmV3yiSBMxR3pRBKJyLcRWt7SL5gEeBb",
}
var smUser1 = Account{
	secret:  "pwRdHmA4cSUKKtFyo4m2vhiiz5g6ym58Noo9dTsUU97mARNjevj",
	address: "zMXMtS2C36J1p3uhTxRFWV8pEhHa8AMMSL",
}
var tableName = "hello2"

func main() {
	c := core.NewChainsql()

	serverName := "peer0.org1.example.com"
	address := "wss://127.0.0.1:6006"
	rootPath := "./certs/root/ca.crt"
	clientCertPath := "./certs/client/client.crt"
	clientKeyPath := "./certs/client/client.key"

	serverName = ""
	rootPath = ""
	clientCertPath = ""
	clientKeyPath = ""

	err := c.Connect(address, rootPath, clientCertPath, clientKeyPath, serverName)
	log.Println("IsConnected:", c.IsConnected())
	if err != nil {
		log.Println(err)
		return
	}
	// var root = Account{
	// 	address: "zHb9CJAWyB4zj91VRWn96DkukG4bwdtyTh",
	// 	secret:  "xnoPBzXtMeMyMHUVTgbuqAfg1SUTb",
	// }
	// c.As(root.address, root.secret)
	c.As(smRoot.address, smRoot.secret)
	//c.SetSchema("44C2C733C17335C11B01BCB0B55340EA422F37307188FF84E6127F8BEBBF0C60")
	//GenerateKey(rand.Reader)
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
	//testGetAccountInfo(c)
	testGetServerInfo(c)
	// testPay(c)
	//testSchemaCreate(c) //创建子链
	//testSchemaModify(c) // 修改子链
	//testGetSchemaList(c) //获取子链列表
	//testGetSchemaInfo(c) //依据子链id获取子链信息
	//testStopSchema(c) //
	//testStartSchema(c)

	//testGetTransaction(c)
	//testGetSchemaId(c)
	//testGenerateAddress(c)
	//testDeleteSchema(c)
	// testGetTransactionResult(c)
	for {
		time.Sleep(time.Second * 10)
	}
}

/*func testGenerateAccount(c *core.Chainsql) {
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
}*/

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
	ret := c.Pay(user2.address, 3000000).Submit("validate_success")
	log.Println(ret)
}

func testSchemaCreate(c *core.Chainsql) {
	schemaInfo := "{\"SchemaName\":\"hello\",\"WithState\":true,\"SchemaAdmin\":\"zHb9CJAWyB4zj91VRWn96DkukG4bwdtyTh\",\"Validators\":[{\"Validator\":{\"PublicKey\":\"02E273A1B6C5D8427F60F490FC8E03D52F35D1A57B8755E19DB0F3B683D9A81239\"}},{\"Validator\":{\"PublicKey\":\"03946B1514971B645F22EE9C96843C5705D388D65C563B7DDEB8B56E29A5283D3C\"}},{\"Validator\":{\"PublicKey\":\"026C3AB1DE848906C5F42316D2B92C3382EEE68FC08F738A2BE9D8C0BC4954936D\"}},{\"Validator\":{\"PublicKey\":\"030F9602B680A71D962111CC3EF9D27601EB9FEC7C3BA24BB4323D94C3A1CF9A04\"}}],\"PeerList\":[{\"Peer\":{\"Endpoint\":\"10.100.0.78:25413\"}},{\"Peer\":{\"Endpoint\":\"10.100.0.78:25414\"}},{\"Peer\":{\"Endpoint\":\"10.100.0.78:25415\"}},{\"Peer\":{\"Endpoint\":\"10.100.0.104:5510\"}}]}"
	//schemaInfo := "{"SchemaName":"子链1","WithState":false,"SchemaAdmin":"zKTqp9kqJBag59YGmL7imH9RfukG6qVtfS","Validators":[{"Validator":{"PublicKey":"037B2D1B1C97A996B44A2FA25765DE5D937247840C960AC6E84D0E3AA8A718F96E"}},{"Validator":{"PublicKey":"038C4245389C8AB8C7665CA4002AEE75EF5D7EEB51A4410D48797BC74F275E9CC3"}},{"Validator":{"PublicKey":"0237788307F53E50D9F799F0D0ABD48258BC41D9418638BD51C481D1848E005443"}}],"PeerList":[{"Peer":{"Endpoint":"192.168.0.242:12260"}},{"Peer":{"Endpoint":"192.168.0.242:12264"}},{"Peer":{"Endpoint":"192.168.0.242:12269"}}]}"
	//schemaInfo := "{\"SchemaName\":\"hello\",\"WithState\":true,\"SchemaAdmin\":\"zN7TwUjJ899xcvNXZkNJ8eFFv2VLKdESsj\",\"Validators\":[{\"Validator\":{\"PublicKey\":\"47F7288B41B45F49342FAC6B65EC529B5ED52F3DDD35140C53BB54A3A7D03F3E9166B0FD574F098F2F9E30526EC8293CE95D4956AD8EC02B34060F0709DCDEA3C5\"}},{\"Validator\":{\"PublicKey\":\"47594A1F76382A89A811B485E3B3414F18967C55A9A2BB90DF7EF36FFF5FDCB915B9C495D66ADEA79DAD97C897596F6FE093C7CDADF90BDD0C91B99D8D014C1B05\"}},{\"Validator\":{\"PublicKey\":\"47C45A7D125E49FDFF1DE6C08F738122FFDC7171E0F5AFA794D02198E35F7F1B1F07CE271CDBF9BA4DB94AA087BE4F59F2A15A60868BE4ACFA86D13B448CD06038\"}}],\"PeerList\":[{\"Peer\":{\"Endpoint\":\"192.168.177.109:5432\"}},{\"Peer\":{\"Endpoint\":\"192.168.177.109:5433\"}},{\"Peer\":{\"Endpoint\":\"192.168.177.109:5441\"}}]}"
	ret := c.CreateSchema(schemaInfo).Submit("validate_success")
	log.Println(ret)
}

func testSchemaModify(c *core.Chainsql) {
	schemaInfo := "{\"SchemaID\":\"2A4C77CF90F1EE9495733D95549B6E0C2ECD61A49343D04C75FE2FD366366A13\",\"Validators\":[{\"Validator\":{\"PublicKey\":\"471D5247096A08552746B5E321E925639A43D6A3ED5F48E48C678E1630F6B92F88EBE20579DCEF85371C43D7305787CAA9AADF7D705BDD1523BBCCF9865FEB34A4\"}}],\"PeerList\":[{\"Peer\":{\"Endpoint\":\"10.100.0.104:5410\"}}]}"
	//schemaInfo := "{\"SchemaName\":\"hello\",\"WithState\":false,\"SchemaAdmin\":\"zBonp9s7isAaDUPcfrFfYjNnhgeznoBHxF\",\"Validators\":\"fhfhhfhfhfhfh\",\"PeerList\":[{\"Peer\":{\"Endpoint\":\"127.0.0.1:15125\"}},{\"Peer\":{\"Endpoint\":\"127.0.0.1:25125\"}},{\"Peer\":{\"Endpoint\":\"127.0.0.1:35125\"}}]}"
	ret := c.ModifySchema("schema_add", schemaInfo).Submit("validate_success")
	log.Println(ret)
}

func testGetSchemaList(c *core.Chainsql) {
	//param := "{\"running\":false}"
	param := ""
	ret, err := c.GetSchemaList(param)
	log.Println(ret)
	log.Println(err)
}

func testGetSchemaInfo(c *core.Chainsql) {
	schemaID := "44C2C733C17335C11B01BCB0B55340EA422F37307188FF84E6127F8BEBBF0C60"
	ret, err := c.GetSchemaInfo(schemaID)
	log.Println(ret)
	log.Println(err)
}

func testStopSchema(c *core.Chainsql) {
	schemaID := "2A4C77CF90F1EE9495733D95549B6E0C2ECD61A49343D04C75FE2FD366366A13"
	ret, err := c.StopSchema(schemaID)
	log.Println(ret)
	log.Println(err)
}

func testStartSchema(c *core.Chainsql) {
	schemaID := "2A4C77CF90F1EE9495733D95549B6E0C2ECD61A49343D04C75FE2FD366366A13"
	ret, err := c.StartSchema(schemaID)
	log.Println(ret)
	log.Println(err)
}

func testGetTransaction(c *core.Chainsql) {
	txHash := "52F7B8FEBC0FEDB56AAA26FA3A8F8C3D507C2E0706BF522644B308AD0190DB88"
	ret, err := c.GetTransaction(txHash)
	log.Println(ret)
	log.Println(err)
}

func testGetSchemaId(c *core.Chainsql) {
	txHash := "1E1EA9E9936574D17646EE9801B72B106DB35D13923FE4357746AD2DD2135C78"
	ret, err := c.GetSchemaId(txHash)
	log.Println(ret)
	log.Println(err)
}

func testGenerateAddress(c *core.Chainsql) {
	//option := ""
	option := "{\"algorithm\":\"softGMAlg\", \"secret\":\"pwRdHmA4cSUKKtFyo4m2vhiiz5g6ym58Noo9dTsUU97mARNjevj\"}"
	//xp6FwxZP1rrmPy2GDTobvHTgnZnrC
	//xnoPBzXtMeMyMHUVTgbuqAfg1SUTb
	//option := "{\"algorithm\":\"secp256k1\", \"secret\":\"xp6FwxZP1rrmPy2GDTobvHTgnZnrC\"}"
	ret, err := c.GenerateAddress(option)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(ret)
	}
}

func testDeleteSchema(c *core.Chainsql) {
	schemaID := "362C1B00328F08A3CE093FE64B0BC88951D364DC6D2B44601268E7F273354B4A"
	ret := c.DeleteSchema(schemaID).Submit("validate_success")
	log.Println(ret)
}

func testGetTransactionResult(c *core.Chainsql) {
	txHash := "48AF80C9169790BE33CBD8C5F32151C9D4483D855B195D188854636EB6551AF5"
	ret, err := c.GetTransactionResult(txHash)
	if err != nil {
		log.Println(err)
	} else {
		// {"id":2,"result":{"ledger_index":28,"transaction_result":"tesSUCCESS","tx_hash":"48AF80C9169790BE33CBD8C5F32151C9D4483D855B195D188854636EB6551AF5","tx_status":"validated"},"status":"success","type":"response"}
		log.Println(ret)
	}
}
