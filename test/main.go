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

var tableName = "hello2"

func main() {
	c := core.NewChainsql()
	// err := c.Connect("ws://127.0.0.1:6006")
	// log.Println("IsConnected:", c.IsConnected())
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// var root = Account{
	// 	address: "zHb9CJAWyB4zj91VRWn96DkukG4bwdtyTh",
	// 	secret:  "xnoPBzXtMeMyMHUVTgbuqAfg1SUTb",
	// }
	// var user = Account{
	// 	address: "zBonp9s7isAaDUPcfrFfYjNnhgeznoBHxF",
	// 	secret:  "xn2FhQLRQqhKJeNhpgMzp2PGAYbdw",
	// }
	// c.As(user.address, user.secret)
	// c.As(root.address, root.secret)
	// c.Use(root.address)

	// // testSubLedger(c)
	testGenerateAccount(c)
	// testInsert(c)
	// testGetLedger(c)
	// testSignPlainText(c)

	// testGetTableData(c)

	// testGetBySqlUser(c)
	// testWebsocket()
	// testTickerGet(c)
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
	ret := c.Table(tableName).Insert(string(data)).Submit("db_success")
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
