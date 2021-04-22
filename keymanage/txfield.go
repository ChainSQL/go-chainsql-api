package keymanage

type SerializedTypeID int

const (
	STI_UNKNOWN    = -2
	STI_DONE       = -1
	STI_NOTPRESENT = 0
	//types (common)
	STI_UINT16  = 1
	STI_UINT32  = 2
	STI_UINT64  = 3
	STI_HASH128 = 4
	STI_HASH256 = 5
	STI_AMOUNT  = 6
	STI_VL      = 7
	STI_ACCOUNT = 8

	STI_ENTRY = 9
	// 9-13 are reserved
	STI_OBJECT = 14
	STI_ARRAY  = 15

	// types (uncommon)
	STI_UINT8     = 16
	STI_HASH160   = 17
	STI_PATHSET   = 18
	STI_VECTOR256 = 19

	// high level types
	// cannot be serialized inside other types
	STI_TRANSACTION = 10001
	STI_LEDGERENTRY = 10002
	STI_VALIDATION  = 10003
	STI_METADATA    = 10004
)

type TxField struct {
	fildID SerializedTypeID
	iOrder int
}

//GetFieldID return fildID
func (o *TxField) GetFieldID() SerializedTypeID {
	return o.fildID
}

//GetFieldID return fildID
func (o *TxField) GetOrder() int {
	return o.iOrder
}

//GetFieldID return fildID
func (o *TxField) GetCombinedKey() int {
	return o.iOrder<<16 + o.iOrder
}

//var MapFieldInitData map[string]TxField
//MapFieldInitData: global param for field

var MapFieldInitData map[string]TxField = map[string]TxField{
	"LedgerEntry": TxField{STI_LEDGERENTRY, 257},
	"Transaction": TxField{STI_TRANSACTION, 257},
	"Validation":  TxField{STI_VALIDATION, 257},
	"Metadata":    TxField{STI_METADATA, 257},
	"hash":        TxField{STI_HASH256, 257},
	"index":       TxField{STI_HASH256, 258},

	"CloseResolution":   TxField{STI_UINT8, 1},
	"Method":            TxField{STI_UINT8, 2},
	"TransactionResult": TxField{STI_UINT8, 3},
	"Deleted":           TxField{STI_UINT8, 50},

	"TickSize": TxField{STI_UINT8, 16},

	"LedgerEntryType": TxField{STI_UINT16, 1},
	"TransactionType": TxField{STI_UINT16, 2},
	"SignerWeight":    TxField{STI_UINT16, 3},
	"OpType":          TxField{STI_UINT16, 50},

	// 32-bit integers (common)
	"Flags":             TxField{STI_UINT32, 2},
	"SourceTag":         TxField{STI_UINT32, 3},
	"Sequence":          TxField{STI_UINT32, 4},
	"PreviousTxnLgrSeq": TxField{STI_UINT32, 5},
	"LedgerSequence":    TxField{STI_UINT32, 6},
	"CloseTime":         TxField{STI_UINT32, 7},
	"ParentCloseTime":   TxField{STI_UINT32, 8},
	"SigningTime":       TxField{STI_UINT32, 9},
	"Expiration":        TxField{STI_UINT32, 10},
	"TransferRate":      TxField{STI_UINT32, 11},
	"WalletSize":        TxField{STI_UINT32, 12},
	"OwnerCount":        TxField{STI_UINT32, 13},
	"DestinationTag":    TxField{STI_UINT32, 14},

	// 32-bit integers (uncommon)
	"HighQualityIn":       TxField{STI_UINT32, 16},
	"HighQualityOut":      TxField{STI_UINT32, 17},
	"LowQualityIn":        TxField{STI_UINT32, 18},
	"LowQualityOut":       TxField{STI_UINT32, 19},
	"QualityIn":           TxField{STI_UINT32, 20},
	"QualityOut":          TxField{STI_UINT32, 21},
	"StampEscrow":         TxField{STI_UINT32, 22},
	"BondAmount":          TxField{STI_UINT32, 23},
	"LoadFee":             TxField{STI_UINT32, 24},
	"OfferSequence":       TxField{STI_UINT32, 25},
	"FirstLedgerSequence": TxField{STI_UINT32, 26},
	"LastLedgerSequence":  TxField{STI_UINT32, 27},

	"TransactionIndex":  TxField{STI_UINT32, 28},
	"OperationLimit":    TxField{STI_UINT32, 29},
	"ReferenceFeeUnits": TxField{STI_UINT32, 30},
	"ReserveBase":       TxField{STI_UINT32, 31},
	"ReserveIncrement":  TxField{STI_UINT32, 32},
	"SetFlag":           TxField{STI_UINT32, 33},
	"ClearFlag":         TxField{STI_UINT32, 34},
	"SignerQuorum":      TxField{STI_UINT32, 35},
	"CancelAfter":       TxField{STI_UINT32, 36},
	"FinishAfter":       TxField{STI_UINT32, 37},
	"SignerListID":      TxField{STI_UINT32, 38},
	"SettleDelay":       TxField{STI_UINT32, 39},
	"TxnLgrSeq":         TxField{STI_UINT32, 50},
	"CreateLgrSeq":      TxField{STI_UINT32, 51},
	"NeedVerify":        TxField{STI_UINT32, 52},

	"IndexNext":     TxField{STI_UINT64, 1},
	"IndexPrevious": TxField{STI_UINT64, 2},
	"BookNode":      TxField{STI_UINT64, 3},
	"OwnerNode":     TxField{STI_UINT64, 4},
	"BaseFee":       TxField{STI_UINT64, 5},
	"ExchangeRate":  TxField{STI_UINT64, 6},
	"LowNode":       TxField{STI_UINT64, 7},
	"HighNode":      TxField{STI_UINT64, 8},

	// 128-bit
	"EmailHash": TxField{STI_HASH128, 1},

	// 160-bit (common)
	"TakerPaysCurrency": TxField{STI_HASH160, 1},
	"TakerPaysIssuer":   TxField{STI_HASH160, 2},
	"TakerGetsCurrency": TxField{STI_HASH160, 3},
	"TakerGetsIssuer":   TxField{STI_HASH160, 4},
	"NameInDB":          TxField{STI_HASH160, 50},
	// 256-bit (common)
	"LedgerHash":        TxField{STI_HASH256, 1},
	"ParentHash":        TxField{STI_HASH256, 2},
	"TransactionHash":   TxField{STI_HASH256, 3},
	"AccountHash":       TxField{STI_HASH256, 4},
	"PreviousTxnID":     TxField{STI_HASH256, 5},
	"LedgerIndex":       TxField{STI_HASH256, 6},
	"WalletLocator":     TxField{STI_HASH256, 7},
	"RootIndex":         TxField{STI_HASH256, 8},
	"AccountTxnID":      TxField{STI_HASH256, 9},
	"PrevTxnLedgerHash": TxField{STI_HASH256, 50},
	"TxnLedgerHash":     TxField{STI_HASH256, 51},
	"TxCheckHash":       TxField{STI_HASH256, 52},
	"CreatedLedgerHash": TxField{STI_HASH256, 53},
	"CreatedTxnHash":    TxField{STI_HASH256, 54},
	"CurTxHash":         TxField{STI_HASH256, 55},
	"FutureTxHash":      TxField{STI_HASH256, 56},

	// 256-bit (uncommon)
	"BookDirectory": TxField{STI_HASH256, 16},
	"InvoiceID":     TxField{STI_HASH256, 17},
	"Nickname":      TxField{STI_HASH256, 18},
	"Amendment":     TxField{STI_HASH256, 19},
	"TicketID":      TxField{STI_HASH256, 20},
	"Digest":        TxField{STI_HASH256, 21},
	"Channel":       TxField{STI_HASH256, 220},

	// currency amount (common)
	"Amount":      TxField{STI_AMOUNT, 1},
	"Balance":     TxField{STI_AMOUNT, 2},
	"LimitAmount": TxField{STI_AMOUNT, 3},
	"TakerPays":   TxField{STI_AMOUNT, 4},
	"TakerGets":   TxField{STI_AMOUNT, 5},
	"LowLimit":    TxField{STI_AMOUNT, 6},
	"HighLimit":   TxField{STI_AMOUNT, 7},
	"Fee":         TxField{STI_AMOUNT, 8},
	"SendMax":     TxField{STI_AMOUNT, 9},
	"DeliverMin":  TxField{STI_AMOUNT, 10},

	// currency amount (uncommon)
	"MinimumOffer":    TxField{STI_AMOUNT, 16},
	"RippleEscrow":    TxField{STI_AMOUNT, 17},
	"DeliveredAmount": TxField{STI_AMOUNT, 18},

	// variable length (common)
	"PublicKey":     TxField{STI_VL, 1},
	"SigningPubKey": TxField{STI_VL, 3},
	"Signature":     TxField{STI_VL, 6},
	"MessageKey":    TxField{STI_VL, 2},
	"TxnSignature":  TxField{STI_VL, 4},
	"Domain":        TxField{STI_VL, 7},
	"FundCode":      TxField{STI_VL, 8},
	"RemoveCode":    TxField{STI_VL, 9},
	"ExpireCode":    TxField{STI_VL, 10},
	"CreateCode":    TxField{STI_VL, 11},
	"MemoType":      TxField{STI_VL, 12},
	"MemoData":      TxField{STI_VL, 13},
	"MemoFormat":    TxField{STI_VL, 14},

	// variable length (uncommon)
	"Fulfillment":     TxField{STI_VL, 16},
	"Condition":       TxField{STI_VL, 17},
	"MasterSignature": TxField{STI_VL, 18},
	"Token":           TxField{STI_VL, 50},
	"TableName":       TxField{STI_VL, 51},
	"Raw":             TxField{STI_VL, 52},
	"TableNewName":    TxField{STI_VL, 53},
	"AutoFillField":   TxField{STI_VL, 54},
	"Statements":      TxField{STI_VL, 55},

	// account
	"Account":         TxField{STI_ACCOUNT, 1},
	"Owner":           TxField{STI_ACCOUNT, 2},
	"Destination":     TxField{STI_ACCOUNT, 3},
	"Issuer":          TxField{STI_ACCOUNT, 4},
	"Target":          TxField{STI_ACCOUNT, 7},
	"RegularKey":      TxField{STI_ACCOUNT, 8},
	"User":            TxField{STI_ACCOUNT, 50},
	"OriginalAddress": TxField{STI_ACCOUNT, 51},

	"Entry": TxField{STI_ENTRY, 1},
	// path set
	"Paths": TxField{STI_PATHSET, 1},

	// vector of 256-bit
	"Indexes":    TxField{STI_VECTOR256, 1},
	"Hashes":     TxField{STI_VECTOR256, 2},
	"Amendments": TxField{STI_VECTOR256, 3},
	// inner object
	// OBJECT/1 is reserved for end of object
	"TransactionMetaData": TxField{STI_OBJECT, 2},
	"CreatedNode":         TxField{STI_OBJECT, 3},
	"DeletedNode":         TxField{STI_OBJECT, 4},
	"ModifiedNode":        TxField{STI_OBJECT, 5},
	"PreviousFields":      TxField{STI_OBJECT, 6},
	"FinalFields":         TxField{STI_OBJECT, 7},
	"NewFields":           TxField{STI_OBJECT, 8},
	"TemplateEntry":       TxField{STI_OBJECT, 9},
	"Memo":                TxField{STI_OBJECT, 10},
	"SignerEntry":         TxField{STI_OBJECT, 11},
	"Table":               TxField{STI_OBJECT, 50},

	// inner object (uncommon)
	"Signer":   TxField{STI_OBJECT, 16},
	"Majority": TxField{STI_OBJECT, 18},

	"Signers":       TxField{STI_ARRAY, 3},
	"SignerEntries": TxField{STI_ARRAY, 4},
	"Template":      TxField{STI_ARRAY, 5},
	"Necessary":     TxField{STI_ARRAY, 6},
	"Sufficient":    TxField{STI_ARRAY, 7},
	"AffectedNodes": TxField{STI_ARRAY, 8},
	"Memos":         TxField{STI_ARRAY, 9},
	"TableEntries":  TxField{STI_ARRAY, 50},
	"Tables":        TxField{STI_ARRAY, 51},
	"Users":         TxField{STI_ARRAY, 52},

	// array of objects (uncommon)
	"Majorities": TxField{STI_ARRAY, 16},
}

type TxType int

const (
	ttINVALID = -1

	ttPAYMENT         = 0
	ttESCROW_CREATE   = 1
	ttESCROW_FINISH   = 2
	ttACCOUNT_SET     = 3
	ttESCROW_CANCEL   = 4
	ttREGULAR_KEY_SET = 5
	ttNICKNAME_SET    = 6 // open
	ttOFFER_CREATE    = 7
	ttOFFER_CANCEL    = 8
	no_longer_used    = 9
	ttTICKET_CREATE   = 10
	ttTICKET_CANCEL   = 11
	ttSIGNER_LIST_SET = 12
	ttPAYCHAN_CREATE  = 13
	ttPAYCHAN_FUND    = 14
	ttPAYCHAN_CLAIM   = 15

	ttTRUST_SET      = 20
	ttTABLELISTSET   = 21
	ttSQLSTATEMENT   = 22
	ttSQLTRANSACTION = 23
	ttAMENDMENT      = 100
	ttFEE            = 101
)
