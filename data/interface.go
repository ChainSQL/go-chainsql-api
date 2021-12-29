package data

import (
	"io"
)

type Hashable interface {
	GetType() string
	Prefix() HashPrefix
	GetHash() *Hash256
}

type Signer interface {
	Hashable
	InitialiseForSigning(kType KeyType)
	SigningPrefix() HashPrefix
	GetPublicKey() *PublicKey
	GetSignature() *VariableLength
	SetTxBase(seq uint32, fee Value, astLedgerSequence *uint32, account Account)
	GetRaw() string
	GetStatements() string
}

type Router interface {
	Hashable
	SuppressionId(keyType KeyType) Hash256
}

type Storer interface {
	Hashable
	Ledger() uint32
	NodeType() NodeType
	NodeId() *Hash256
}

type LedgerEntry interface {
	Storer
	GetLedgerEntryType() LedgerEntryType
	GetLedgerIndex() *Hash256
	GetPreviousTxnId() *Hash256
	Affects(Account) bool
}

type Transaction interface {
	Signer
	GetTransactionType() TransactionType
	GetBase() *TxBase
	PathSet() PathSet
}

type Wire interface {
	Unmarshal(Reader) error
	Marshal(io.Writer) error
}
