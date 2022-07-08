module github.com/ChainSQL/go-chainsql-api

go 1.15

require (
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/buger/jsonparser v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/kr/text v0.2.0 // indirect
	github.com/peersafe/gm-crypto v1.0.2
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
)

replace github.com/peersafe/gm-crypto v1.0.2 => gitlab.peersafe.cn/fabric/gm-crypto v1.0.2
