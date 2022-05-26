package net

import (
	"sync"

	"github.com/buger/jsonparser"
)

//PrepareTable return the account sequence and table NameInDB
func PrepareTable(client *Client, name string) (uint32, string, error) {
	w := new(sync.WaitGroup)
	w.Add(2)
	var seq uint32 = 0
	nameInDB := ""
	err := error(nil)
	go func() {
		defer w.Done()
		info, errTmp := client.GetAccountInfo(client.Auth.Address)
		if errTmp != nil {
			err = errTmp
			return
		}
		sequence, errTmp := jsonparser.GetInt([]byte(info), "result", "account_data", "Sequence")
		if errTmp != nil {
			err = errTmp
			return
		}
		seq = uint32(sequence)
	}()
	go func() {
		defer w.Done()
		nameInDBTmp, errTmp := client.GetNameInDB(client.Auth.Owner, name)
		if errTmp != nil {
			err = errTmp
			return
		}
		nameInDB = nameInDBTmp
	}()
	w.Wait()
	return seq, nameInDB, err
}

//PrepareTable return the account sequence and table NameInDB
func PrepareRipple(client *Client) (uint32, error) {

	var seq uint32 = 0
	err := error(nil)

	info, errTmp := client.GetAccountInfo(client.Auth.Address)
	if errTmp != nil {
		err = errTmp
		return 0, errTmp
	}
	sequence, errTmp := jsonparser.GetInt([]byte(info), "result", "account_data", "Sequence")
	if errTmp != nil {
		err = errTmp
		return 0, errTmp
	}
	seq = uint32(sequence)

	return seq, err
}
