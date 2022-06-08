// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the Go wrapper for the constant-time, 64-bit assembly
// implementation of sm2 curve. The optimizations performed here are described in
// detail in:
// S.Gueron and V.Krasnov, "Fast prime field elliptic-curve cryptography with
//                          256-bit primes"
// http://link.springer.com/article/10.1007%2Fs13389-014-0090-x
// https://eprint.iacr.org/2013/816.pdf

package sm2

import (
	"crypto/elliptic"
	"math/big"
)

type (
	p256Curve struct {
		*elliptic.CurveParams
	}
)

var (
	p256            p256Curve
)

func initP256() {
	// See FIPS 186-3, section D.2.3
	p256.CurveParams = &elliptic.CurveParams{Name: "SM2-P-256"}
	p256.P, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)
	p256.N, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
	p256.B, _ = new(big.Int).SetString("28E9FA9E9D9F5E344D5A9E4BCF6509A7F39789F515AB8F92DDBCBD414D940E93", 16)
	p256.Gx, _ = new(big.Int).SetString("32C4AE2C1F1981195F9904466A39C9948FE30BBFF2660BE1715A4589334C74C7", 16)
	p256.Gy, _ = new(big.Int).SetString("BC3736A2F4F6779C59BDCEE36B692153D0A9877CC62A474002DF32E52139F0A0", 16)
	p256.BitSize = 256
}