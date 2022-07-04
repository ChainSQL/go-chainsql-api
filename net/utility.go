package net

import (
	"sync"

	"github.com/buger/jsonparser"
)

const LastLedgerSeqOffset = 20

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

func PrepareLastLedgerSeqAndFee(client *Client) (int64, uint32, error) {
	var fee int64 = 10
	var lastLedgerSeq uint32 = 10
	if client.ServerInfo.Updated {
		lastLedgerSeq = uint32(client.ServerInfo.LedgerIndex + LastLedgerSeqOffset)
		fee = int64(client.ServerInfo.ComputeFee())
	} else {
		ledgerIndex, err := client.GetLedgerVersion()
		if err != nil {
			return 0, 0, err
		}
		lastLedgerSeq = uint32(ledgerIndex + LastLedgerSeqOffset)

		fee = 50
	}
	return fee, lastLedgerSeq, nil
}
