package net

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Reconnected func()

// WebsocketManager is a websocket client manager
type WebsocketManager struct {
	conn          *websocket.Conn
	url           string
	sendMsgChan   chan string
	recvMsgChan   chan string
	isAlive       bool
	timeout       int // used for reconnecting
	muxRead       *sync.Mutex
	muxWrite      *sync.Mutex
	onReconnected Reconnected
}

// NewWsClientManager is a constructor
func NewWsClientManager(url string, timeout int) *WebsocketManager {
	var sendChan = make(chan string, 1024)
	var recvChan = make(chan string, 1024)
	var conn *websocket.Conn
	return &WebsocketManager{
		url:           url,
		conn:          conn,
		sendMsgChan:   sendChan,
		recvMsgChan:   recvChan,
		isAlive:       false,
		timeout:       timeout,
		muxRead:       new(sync.Mutex),
		muxWrite:      new(sync.Mutex),
		onReconnected: nil,
	}
}

// 链接服务端
func (wsc *WebsocketManager) dail() error {
	var err error
	// log.Printf("connecting to %s", wsc.url)
	websocket.DefaultDialer.HandshakeTimeout = 10 * time.Second
	wsc.conn, _, err = websocket.DefaultDialer.Dial(wsc.url, nil)
	if err != nil {
		log.Printf("connecting to %s failed,err:%s", wsc.url, err.Error())
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
			if wsc.isAlive {
				msg := <-wsc.sendMsgChan

				if wsc.conn != nil {
					wsc.muxWrite.Lock()
					wsc.conn.SetWriteDeadline(time.Now().Add(time.Duration(wsc.timeout)))
					err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
					wsc.muxWrite.Unlock()
					if err != nil {
						wsc.close()
						log.Println("write:", err)
						wsc.sendMsgChan <- msg
						// break
					}
				}
			} else {
				time.Sleep(time.Second * 1)
			}
		}
	}()
}

func (wsc *WebsocketManager) OnReconnected(cb Reconnected) {
	wsc.onReconnected = cb
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
					wsc.close()
					log.Println("read:", err)
					// break
				} else {
					// log.Printf("recv: %s", message)
					// 需要读取数据，不然会阻塞
					wsc.recvMsgChan <- string(message)
				}
			} else {
				time.Sleep(time.Second * 1)
			}
		}
	}()
}

func (wsc *WebsocketManager) close() {
	if wsc.isAlive {
		defer wsc.conn.Close()
		wsc.isAlive = false
	}
}

func (wsc *WebsocketManager) checkReconnect() {
	go func() {
		for {
			if wsc.isAlive == false {
				log.Println("checkReconnect ws disconnected,reconnect!")
				err := wsc.connectAndRun()
				if err == nil && wsc.onReconnected != nil {
					wsc.onReconnected()
				}
			}
			time.Sleep(time.Second * time.Duration(wsc.timeout))
		}
	}()
}

//Start 开启服务并重连
func (wsc *WebsocketManager) Start() error {
	err := wsc.connectAndRun()
	wsc.checkReconnect()
	return err
}

func (wsc *WebsocketManager) connectAndRun() error {
	err := wsc.dail()
	if err == nil {
		wsc.sendMsgThread()
		wsc.readMsgThread()
	}
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

func (wsc *WebsocketManager) IsConnected() bool {
	return wsc.isAlive
}

// func main() {
//     wsc := NewWsClientManager("192.168.12.15", "10086", "/v1", 10)
//     wsc.start()
//     var w1 sync.WaitGroup
//     w1.Add(1)
//     w1.Wait()
//	}
