package proc

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anoideaopen/ifc/blocks"
	"github.com/anoideaopen/ifc/utils"
	"github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

const (
	defaultTrysCount = 3
)

var (
	methods = []string{
		"changeMultisigPublicKey", "addMultisigWithBase58Signature", "addMultisig",
		"delFromList", "addToList", "removeAdditionalKey", "addAdditionalKey", "setAccountInfo",
		"removeAddressFromNominee", "addAddressForNominee", "removeRights", "addRights",
		"addUser", "addUserWithPublicKeyType", "setkyc", "changePublicKeyWithBase58Signature",
		"changePublicKey",
	}
	mapMethods map[string]bool
)

func init() {
	mapMethods = make(map[string]bool)
	for _, m := range methods {
		mapMethods[m] = true
	}
}

func TxsProc(
	actions []*peer.TransactionAction,
	chClient *channel.Client,
) error {
	for _, action := range actions {
		capl, _, err := blocks.GetPayloads(action)
		if err != nil {
			continue
		}

		cpp := &peer.ChaincodeProposalPayload{}
		if err = proto.Unmarshal(capl.GetChaincodeProposalPayload(), cpp); err != nil {
			continue
		}

		cis := &peer.ChaincodeInvocationSpec{}
		if err = proto.Unmarshal(cpp.GetInput(), cis); err != nil {
			continue
		}

		if len(cis.GetChaincodeSpec().GetInput().GetArgs()) == 0 {
			continue
		}

		if !mapMethods[string(cis.GetChaincodeSpec().GetInput().GetArgs()[0])] {
			continue
		}

		var resp channel.Response
		for j := 0; j < defaultTrysCount; j++ {
			if resp, err = chClient.Execute(channel.Request{
				ChaincodeID: utils.ChaincodeACL,
				Fcn:         string(cis.GetChaincodeSpec().GetInput().GetArgs()[0]),
				Args:        cis.GetChaincodeSpec().GetInput().GetArgs()[1:],
			},
				channel.WithTimeout(fab.Execute, time.Second*30),
			); err == nil {
				break
			}

			if strings.Contains(err.Error(), "already exists") ||
				strings.Contains(err.Error(), "incorrect nonce") {
				err = nil
				break
			}

			time.Sleep(time.Second)
		}

		if err != nil {
			err = fmt.Errorf("couldn't send state: %w", err)
			panic(err)
		}

		if resp.ChaincodeStatus != http.StatusOK && resp.ChaincodeStatus != 0 {
			err = fmt.Errorf("invalid response status: %d", resp.ChaincodeStatus)
			panic(err)
		}
	}
	return nil
}
