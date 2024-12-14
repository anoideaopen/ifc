package blocks

import (
	"github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
)

// UnmarshalEnvelope unmarshal bytes to an Envelope structure.
func UnmarshalEnvelope(encoded []byte) (*common.Envelope, error) {
	envelope := &common.Envelope{}
	err := proto.Unmarshal(encoded, envelope)

	return envelope, errors.Wrap(err, "error unmarshaling Envelope")
}

// ExtractPayload retrieves the payload of a given envelope and unmarshal it.
func ExtractPayload(envelope *common.Envelope) (*common.Payload, error) {
	payload := &common.Payload{}
	err := proto.Unmarshal(envelope.GetPayload(), payload)

	return payload, errors.Wrap(err, "no payload in envelope")
}

// UnmarshalChannelHeader returns a ChannelHeader from bytes.
func UnmarshalChannelHeader(bytes []byte) (*common.ChannelHeader, error) {
	chdr := &common.ChannelHeader{}
	err := proto.Unmarshal(bytes, chdr)

	return chdr, errors.Wrap(err, "error unmarshalling ChannelHeader")
}

// GetTransaction Get Transaction from bytes.
func GetTransaction(txBytes []byte) (*peer.Transaction, error) {
	tx := &peer.Transaction{}
	err := proto.Unmarshal(txBytes, tx)

	return tx, errors.Wrap(err, "error unmarshalling Transaction")
}

// GetPayloads gets the underlying payload objects in a TransactionAction.
func GetPayloads(
	txActions *peer.TransactionAction,
) (*peer.ChaincodeActionPayload, *peer.ChaincodeAction, error) {
	ccPayload, err := GetChaincodeActionPayload(txActions.GetPayload())
	if err != nil {
		return nil, nil, err
	}

	if ccPayload.GetAction().GetProposalResponsePayload() == nil {
		return nil, nil, err
	}

	pRespPayload, err := GetProposalResponsePayload(
		ccPayload.GetAction().GetProposalResponsePayload(),
	)
	if err != nil {
		return nil, nil, err
	}

	if pRespPayload.Extension == nil {
		return nil, nil, errors.New("response payload is missing extension")
	}

	respPayload, err := GetChaincodeAction(pRespPayload.GetExtension())
	if err != nil {
		return ccPayload, nil, err
	}

	return ccPayload, respPayload, nil
}

// GetChaincodeActionPayload Get ChaincodeActionPayload from bytes.
func GetChaincodeActionPayload(
	capBytes []byte,
) (*peer.ChaincodeActionPayload, error) {
	c := &peer.ChaincodeActionPayload{}
	err := proto.Unmarshal(capBytes, c)

	return c, errors.Wrap(err, "error unmarshaling ChaincodeActionPayload")
}

// GetProposalResponsePayload gets the proposal response payload.
func GetProposalResponsePayload(
	prpBytes []byte,
) (*peer.ProposalResponsePayload, error) {
	prp := &peer.ProposalResponsePayload{}
	err := proto.Unmarshal(prpBytes, prp)

	return prp, errors.Wrap(err, "error unmarshaling ProposalResponsePayload")
}

// GetChaincodeAction gets the ChaincodeAction given chaincode action bytes.
func GetChaincodeAction(caBytes []byte) (*peer.ChaincodeAction, error) {
	chaincodeAction := &peer.ChaincodeAction{}
	err := proto.Unmarshal(caBytes, chaincodeAction)

	return chaincodeAction, errors.Wrap(err, "error unmarshalling ChaincodeAction")
}
