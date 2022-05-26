package tls

import (
	"bufio"
	"github.com/peersafe/gm-crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"testing"
)

var DefaultTLSCipherSuites = []uint16{
	GMTLS_SM2_WITH_SM4_SM3,
}

var (
	AuthType = RequestClientCert //fabric单证书模式
	//AuthType         = RequireAndVerifyClientCert //fabric双证书模式
	SERVER_CA        = "./gmcert/ca.crt"     //服务端根证书
	SERVER_SIGN_CERT = "./gmcert/server.crt" //服务端签名证书
	SERVER_SIGN_KEY  = "./gmcert/server.key" //服务端签名私钥
	SERVER_ENC_CERT  = "./gmcert/server.crt" //服务端加密证书
	SERVER_ENC_KEY   = "./gmcert/server.key" //服务端加密私钥
	CLIENT_CA        = "./gmcert/ca.crt"     //客户端根证书
	CLIENT_SIGN_CERT = "./gmcert/client.crt" //客户端签名证书
	CLIENT_SIGN_KEY  = "./gmcert/client.key" //客户端签名私钥
	CLIENT_ENC_CERT  = "./gmcert/client.crt" //客户端加密证书
	CLIENT_ENC_KEY   = "./gmcert/client.key" //客户端签名私钥
)

func TestTLS(t *testing.T) {
	signcert, err := LoadX509KeyPair(SERVER_SIGN_CERT, SERVER_SIGN_KEY)
	if err != nil {
		log.Println(err)
		return
	}
	//enccert, err := LoadX509KeyPair(SERVER_ENC_CERT, SERVER_ENC_KEY)
	//	//if err != nil {
	//	//	log.Println(err)
	//	//	return
	//	//}
	caPem, err := ioutil.ReadFile(SERVER_CA)
	if err != nil {
		log.Fatalf("Failed to load ca cert %v", err)
	}
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(caPem)
	c := &Config{
		GMSupport:    &GMSupport{},
		ClientAuth:   AuthType,
		Certificates: []Certificate{signcert /*, enccert*/},
		ClientCAs:    certpool,
		CipherSuites: DefaultTLSCipherSuites,
	}

	ln, err := Listen("tcp", ":8080", c)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Start to server address:%v  clientAuthType=%v\n", ln.Addr(),AuthType)
	go client()
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(msg)
		n, err := conn.Write([]byte("server pong!\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}

func client() {
	caPem, err := ioutil.ReadFile(CLIENT_CA)
	if err != nil {
		log.Fatal(err)
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(caPem) {
		log.Fatal("credentials: failed to append certificates")
	}
	signcert, err := LoadX509KeyPair(CLIENT_SIGN_CERT, CLIENT_SIGN_KEY)
	if err != nil {
		log.Fatal("Failed to Load client keypair")
	}
	//enccert, err := LoadX509KeyPair(CLIENT_ENC_CERT, CLIENT_ENC_KEY)
	//if err != nil {
	//	log.Fatal("Failed to Load client keypair")
	//}
	c := &Config{
		ServerName:   "tls.testserver.com",
		GMSupport:    &GMSupport{},
		Certificates: []Certificate{signcert /*, enccert*/},
		RootCAs:      cp,
		CipherSuites: DefaultTLSCipherSuites,
		//InsecureSkipVerify: true, // Client verifies server's cert if false, else skip.
	}
	serverAddress := "127.0.0.1:8080"
	log.Printf("start client connect %s\n", serverAddress)
	conn, err := Dial("tcp", serverAddress, c)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	n, err := conn.Write([]byte("client write ping!\n"))
	if err != nil {
		log.Println(n, err)
		return
	}
	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}
	log.Println(string(buf[:n]))
}
