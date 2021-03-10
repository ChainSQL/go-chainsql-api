package net

import (
	"sync"

	"github.com/buger/jsonparser"
)

//PrepareTable return the account sequence and table NameInDB
func PrepareTable(client *Client, address string, name string) (int, string, error) {
	w := new(sync.WaitGroup)
	w.Add(2)
	seq := 0
	nameInDB := ""
	err := error(nil)
	go func() {
		defer w.Done()
		info := client.GetAccountInfo(address)
		sequence, errTmp := jsonparser.GetInt([]byte(info), "result", "account_data", "Sequence")
		if errTmp != nil {
			err = errTmp
			return
		}
		seq = int(sequence)
	}()
	go func() {
		defer w.Done()
		nameInDBTmp, errTmp := client.GetNameInDB(address, name)
		if errTmp != nil {
			err = errTmp
			return
		}
		nameInDB = nameInDBTmp
	}()
	w.Wait()
	return seq, nameInDB, err
}
