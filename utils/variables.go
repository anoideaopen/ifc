package utils

import (
	"log"

	"github.com/anoideaopen/glog/std"
)

// EnvPrefix environment prefix.
const (
	EnvPrefix = "APPIFC"

	ChannelACL = "acl"
	ChaincodeACL
)

var (
	// ConnectionFile1 - file for connection to HLF.
	ConnectionFile1 string
	// Org1 - organisation.
	Org1 string
	// User1 - user.
	User1 string

	// ConnectionFile2 - file for connection to HLF.
	ConnectionFile2 string
	// Org2 - organisation.
	Org2 string
	// User2 - user.
	User2 string

	StdLog = std.New(log.Default(), std.LevelTrace)

	ErrorNotFindHLFConf = "could not find HLF config section"
)
