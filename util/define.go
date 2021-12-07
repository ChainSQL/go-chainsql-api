package util

const (
	RInsert = 6
	RUpdate = 8
	RDelete = 9
	RGet    = 7
)

const (
	Payment        = "Payment"
	TableListSet   = "TableListSet"
	SQLStatement   = "SQLStatement"
	SQLTransaction = "SQLTransaction"
	Contract       = "Contract"
)

const (
	SendSuccess     = "send_success"
	ValidateSuccess = "validate_success"
	DbSuccess       = "db_success"
)

const (
	SendError       = "send_error"
	ValidateError   = "validate_error"
	ValidateTimeout = "validate_timeout"
)

const (
	REQUEST_TIMEOUT = 5
	DIAL_TIMEOUT    = 2
)

const (
	SchemaAdd = "schema_add"
	SchemaDel = "schema_del"
)
