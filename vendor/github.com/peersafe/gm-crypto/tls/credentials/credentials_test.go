package credentials

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"testing"

	"github.com/peersafe/gm-crypto/tls"
	"github.com/peersafe/gm-crypto/tls/credentials/echo"
	"github.com/peersafe/gm-crypto/x509"
	"google.golang.org/grpc"
)

const (
	port    = ":50051"
	address = "localhost:50051"
)

var DefaultTLSCipherSuites = []uint16{
	tls.GMTLS_SM2_WITH_SM4_SM3,
}

var end chan bool

type server struct{}

func (s *server) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{Result: req.Req}, nil
}

const cacrt = "../gmcert/ca.crt"
const cakey = "../gmcert/ca.key"
const servercrt = "../gmcert/server.crt"
const serverkey = "../gmcert/server.key"
const clientcrt = "../gmcert/client.crt"
const clientkey = "../gmcert/client.key"

func serverRun() {
	cert, err := tls.LoadX509KeyPair(servercrt, serverkey)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	cacert, err := ioutil.ReadFile(cacrt)
	if err != nil {
		log.Fatal(err)
	}
	certPool.AppendCertsFromPEM(cacert)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("fail to listen: %v", err)
	}
	creds := NewTLS(&tls.Config{
		GMSupport:    &tls.GMSupport{},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert,cert},
		ClientCAs:    certPool,
		CipherSuites: DefaultTLSCipherSuites,
	})

	s := grpc.NewServer(grpc.Creds(creds))
	echo.RegisterEchoServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Serve: %v", err)
	}
}

func clientRun() {
	cert, err := tls.LoadX509KeyPair(clientcrt, clientkey)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	cacert, err := ioutil.ReadFile(cacrt)
	if err != nil {
		log.Fatal(err)
	}
	certPool.AppendCertsFromPEM(cacert)
	creds := NewTLS(&tls.Config{
		GMSupport:    &tls.GMSupport{},
		ServerName:   "tls.testserver.com",
		Certificates: []tls.Certificate{cert,cert},
		RootCAs:      certPool,
		CipherSuites: DefaultTLSCipherSuites,
	})
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("cannot to connect: %v", err)
	}
	defer conn.Close()
	c := echo.NewEchoClient(conn)
	echoTest(c)
	end <- true
}

func echoTest(c echo.EchoClient) {
	r, err := c.Echo(context.Background(), &echo.EchoRequest{Req: "###hello###"})
	if err != nil {
		log.Fatalf("failed to echo: %v", err)
	}
	fmt.Printf("%s\n", r.Result)
}

func Test(t *testing.T) {
	end = make(chan bool, 64)
	go serverRun()
	go clientRun()
	<-end
}
