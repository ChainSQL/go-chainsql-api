package core

import (
	"fmt"

	"github.com/ChainSQL/go-chainsql-api/net"
)

// Base is a struct
type Base struct {
}

func (b *Base) Say() {
	fmt.Println("base")
}

func (b *Base) Say2() {
	fmt.Println("base2")
}

// Ripple is aaa
type Ripple struct {
	*Base
	client *net.Client
	SubmitBase
}

func (r *Ripple) Say() {
	fmt.Println("Ripple")
}

func NewRipple() *Ripple {
	ripple := &Ripple{
		Base:   &Base{},
		client: net.NewClient(),
	}
	ripple.SubmitBase.client = ripple.client
	ripple.SubmitBase.IPrepare = ripple
	return ripple
}

func (r *Ripple) Pay(accountId string, value string) *Ripple {

	return r
}
