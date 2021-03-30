package net

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebsocketManager is a websocket client manager
type WebsocketManager struct {
	conn        *websocket.Conn
	url         string
	sendMsgChan chan string
	recvMsgChan chan string
	isAlive     bool
	timeout     int // used for reconnecting
	muxRead     *sync.Mutex
	muxWrite    *sync.Mutex
}

// NewWsClientManager is a constructor
func NewWsClientManager(url string, timeout int) *WebsocketManager {
	var sendChan = make(chan string, 1024)
	var recvChan = make(chan string, 1024)
	var conn *websocket.Conn
	return &WebsocketManager{
		url:         url,
		conn:        conn,
		sendMsgChan: sendChan,
		recvMsgChan: recvChan,
		isAlive:     false,
		timeout:     timeout,
		muxRead:     new(sync.Mutex),
		muxWrite:    new(sync.Mutex),
	}
}

// 链接服务端
func (wsc *WebsocketManager) dail() error {
	var err error
	// log.Printf("connecting to %s", wsc.url)
	wsc.conn, _, err = websocket.DefaultDialer.Dial(wsc.url, nil)
	if err != nil {
		return err
	}
	wsc.isAlive = true
	log.Printf("connecting to %s success!", wsc.url)
	return nil
}

// 发送消息
func (wsc *WebsocketManager) sendMsgThread() {
	go func() {
		for {
			if wsc.conn != nil && wsc.isAlive {
				msg := <-wsc.sendMsgChan
				wsc.muxWrite.Lock()
				err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
				wsc.muxWrite.Unlock()
				if err != nil {
					defer wsc.conn.Close()
					log.Println("write:", err)
					wsc.isAlive = false
					break
				}
			}
		}
	}()
}

// 读取消息
func (wsc *WebsocketManager) readMsgThread() {
	go func() {
		for {
			if wsc.conn != nil && wsc.isAlive {
				wsc.muxRead.Lock()
				_, message, err := wsc.conn.ReadMessage()
				wsc.muxRead.Unlock()
				if err != nil {
					defer wsc.conn.Close()
					log.Println("read:", err)
					wsc.isAlive = false
					// 出现错误，退出读取，尝试重连
					break
				}
				// log.Printf("recv: %s", message)
				// 需要读取数据，不然会阻塞
				wsc.recvMsgChan <- string(message)
			}

		}
	}()
}

func (wsc *WebsocketManager) checkReconnect() {
	go func() {
		for {
			if wsc.isAlive == false {
				log.Println("checkReconnect ws disconnected,reconnect!")
				wsc.dail()
				wsc.sendMsgThread()
				wsc.readMsgThread()
			}
			time.Sleep(time.Second * time.Duration(wsc.timeout))
		}
	}()
}

//Start 开启服务并重连
func (wsc *WebsocketManager) Start() error {
	err := wsc.dail()
	wsc.sendMsgThread()
	wsc.readMsgThread()
	wsc.checkReconnect()
	return err
}

//Print print the channel buffer size
func (wsc *WebsocketManager) Print() {
	log.Println("read buffer size: %n", len(wsc.recvMsgChan))
	log.Println("write buffer size: %n ", len(wsc.sendMsgChan))
}

//WriteChan return the write channel
func (wsc *WebsocketManager) WriteChan() chan string {
	return wsc.sendMsgChan
}

//ReadChan return the channel used to read
func (wsc *WebsocketManager) ReadChan() chan string {
	return wsc.recvMsgChan
}

// func main() {
//     wsc := NewWsClientManager("192.168.12.15", "10086", "/v1", 10)
//     wsc.start()
//     var w1 sync.WaitGroup
//     w1.Add(1)
//     w1.Wait()
//	}
