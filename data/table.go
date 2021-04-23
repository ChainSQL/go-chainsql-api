package data

import "encoding/hex"

type TableName struct {
	TableName VariableLength `json:"TableName,omitempty"`
	NameInDB  Hash160        `json:"NameInDB,omitempty"`
}

// TableFields defines the table struct
type TableObj struct {
	Table TableName
}

// TableName is sub-struct of Table
type TableNameForGet struct {
	TableName string
	NameInDB  string `json:"NameInDB,omitempty"`
}

//TableObj is sub-struct of Tables
type TableObjForGet struct {
	Table TableNameForGet
}

// FormatTables create the Tables json array in Chainsql transaction
func FormatTables(name string, nameInDB string) []TableObj {
	// var nameHex = VariableLength(fmt.Sprintf("%x", name))
	dbBytes, _ := hex.DecodeString(nameInDB)
	var h160 Hash160
	copy(h160[:], dbBytes)
	var valName = VariableLength(name)
	return []TableObj{
		{
			Table: TableName{
				TableName: valName,
				NameInDB:  h160,
			},
		},
	}
}

func FormatTablesForGet(name string, nameInDB string) []TableObjForGet {
	return []TableObjForGet{
		{
			Table: TableNameForGet{
				TableName: name,
				NameInDB:  nameInDB,
			},
		},
	}
}
