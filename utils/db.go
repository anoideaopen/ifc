package utils

import (
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	PathLastTxID  = "tx_id"
	PathLastBlock = "block"
	DefaultDB     = "./db/state"
)

func Block(db *leveldb.DB) uint64 {
	bBytes, err := db.Get([]byte(PathLastBlock), nil)
	if err != nil {
		return 0
	}

	b, err := strconv.ParseUint(string(bBytes), 10, 64)
	if err != nil {
		return 0
	}

	return b
}

func SetBlock(db *leveldb.DB, b uint64) error {
	err := db.Put([]byte(PathLastBlock), []byte(strconv.FormatUint(b, 10)), nil)
	if err != nil {
		return err
	}
	err = db.Put([]byte(PathLastTxID), []byte(""), nil)
	if err != nil {
		return err
	}
	return nil
}

func TxID(db *leveldb.DB) string {
	txID, err := db.Get([]byte(PathLastTxID), nil)
	if err != nil {
		return ""
	}

	return string(txID)
}

func SetTxID(db *leveldb.DB, id string) error {
	err := db.Put([]byte(PathLastTxID), []byte(id), nil)
	if err != nil {
		return err
	}
	return nil
}
