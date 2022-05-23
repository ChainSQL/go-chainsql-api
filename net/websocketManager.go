package net

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/ChainSQL/go-chainsql-api/util"
	"github.com/gorilla/websocket"
)

type Reconnected func()

// WebsocketManager is a websocket client manager
type WebsocketManager struct {
	conn              *websocket.Conn
	url               string
	sendMsgChan       chan string
	recvMsgChan       chan string
	isAlive           bool
	timeout           int // used for reconnecting
	muxRead           *sync.Mutex
	muxWrite          *sync.Mutex
	muxConnect        *sync.Mutex
	onReconnected     Reconnected
	first_connect     bool
	serverName        string
	tlsRootCertPath   string
	tlsClientKeyPath  string
	tlsClientCertPath string
}

// NewWsClientManager is a constructor
func NewWsClientManager(url, tlsRootCertPath, tlsClientCertPath, tlsClientKeyPath, serverName string, timeout int) *WebsocketManager {
	var sendChan = make(chan string, 1024)
	var recvChan = make(chan string, 1024)
	var conn *websocket.Conn
	return &WebsocketManager{
		url:               url,
		conn:              conn,
		sendMsgChan:       sendChan,
		recvMsgChan:       recvChan,
		isAlive:           false,
		timeout:           timeout,
		muxRead:           new(sync.Mutex),
		muxWrite:          new(sync.Mutex),
		muxConnect:        new(sync.Mutex),
		onReconnected:     nil,
		first_connect:     true,
		serverName:        serverName,
		tlsRootCertPath:   tlsRootCertPath,
		tlsClientKeyPath:  tlsClientKeyPath,
		tlsClientCertPath: tlsClientCertPath,
	}
}

// 链接服务端
func (wsc *WebsocketManager) dail() error {
	var err error
	log.Printf("connecting to %s", wsc.url)
	websocket.DefaultDialer.HandshakeTimeout = util.DIAL_TIMEOUT * time.Second

	// tls config
	tlsConfig := &tls.Config{
		MaxVersion: tls.VersionTLS12,
	}

	if wsc.tlsRootCertPath != "" {
		var caPem []byte
		caPem, err = ioutil.ReadFile(wsc.tlsRootCertPath)
		if err != nil {
			return fmt.Errorf("failed to load tls root cert(path = %s), err = %v", wsc.tlsRootCertPath, err)
		}
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(caPem) {
			return fmt.Errorf("credentials: failed to append certificates")
		}

		tlsConfig.ServerName = wsc.serverName
		tlsConfig.RootCAs = caPool

		if wsc.tlsClientCertPath != "" {
			clientCert, err := tls.LoadX509KeyPair(wsc.tlsClientCertPath, wsc.tlsClientKeyPath)
			tlsConfig.Certificates = []tls.Certificate{clientCert}
			if err != nil {
				return fmt.Errorf("failed to load tls client (cert path = %s; key path = %s), err = %v", wsc.tlsClientCertPath, wsc.tlsClientKeyPath, err)
			}
		}

		dailer := &websocket.Dialer{TLSClientConfig: tlsConfig}
		wsc.conn, _, err = dailer.Dial(wsc.url, nil)
	} else {
		wsc.conn, _, err = websocket.DefaultDialer.Dial(wsc.url, nil)
	}

	if err != nil {
		log.Printf("connecting to %s failed,err:%s", wsc.url, err.Error())
		return err
	}
	wsc.isAlive = true
	log.Printf("connecting to %s success!", wsc.url)
	return nil
}

func (wsc *WebsocketManager) Disconnect() error {
	wsc.muxWrite.Lock()
	wsc.muxConnect.Lock()

	if wsc.conn != nil {
		err := wsc.conn.Close()
		if err != nil {
			return err
		}
		wsc.conn = nil
	}

	wsc.isAlive = false

	wsc.muxConnect.Unlock()
	wsc.muxWrite.Unlock()
	return nil
}

func (wsc *WebsocketManager) SetUrl(url string) {
	wsc.url = url
}

// 发送消息
func (wsc *WebsocketManager) sendMsgThread() {
	go func() {
		for {
			if wsc.isAlive {
				msg := <-wsc.sendMsgChan

				if wsc.conn != nil {
					wsc.muxWrite.Lock()
					// wsc.conn.SetWriteDeadline(time.Now().Add(time.Duration(wsc.timeout)))
					err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
					wsc.muxWrite.Unlock()
					if err != nil {
						wsc.close()
						log.Println("write:", err)
						wsc.sendMsgChan <- msg
						wsc.isAlive = false
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
					wsc.isAlive = false
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
			if !wsc.isAlive {
				log.Println("checkReconnect ws disconnected,reconnect!")
				wsc.muxConnect.Lock()
				err := wsc.connectAndRun()
				wsc.muxConnect.Unlock()
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
	return err
}

func (wsc *WebsocketManager) connectAndRun() error {
	err := wsc.dail()
	if err == nil && wsc.first_connect {
		wsc.sendMsgThread()
		wsc.readMsgThread()
		wsc.checkReconnect()
		wsc.first_connect = false
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
