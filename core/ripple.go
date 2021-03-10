package core

import (
	"fmt"
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
}

func (r *Ripple) Say() {
	fmt.Println("Ripple")
}
