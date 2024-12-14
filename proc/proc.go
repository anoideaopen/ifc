package proc

import (
	"context"
	"fmt"

	"github.com/anoideaopen/ifc/blocks"
	"github.com/anoideaopen/ifc/utils"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	beginBlock    = 0
	sizeOfEvtChan = 200
)

func Process(
	ctx context.Context,
	conn1 string,
	org1 string,
	user1 string,
	conn2 string,
	org2 string,
	user2 string,
) {
	var (
		err error
		cnt int
	)

	hlfCfgTo := config.FromFile(conn2)

	sdkTo, err := fabsdk.New(hlfCfgTo)
	if err != nil {
		err = fmt.Errorf("couldn't initialize SDK: %w", err)
		panic(err)
	}
	defer sdkTo.Close()

	clientChannelContext := sdkTo.ChannelContext(utils.ChannelACL,
		fabsdk.WithUser(user2), fabsdk.WithOrg(org2))

	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		err = fmt.Errorf("couldn't create channel client: %w", err)
		panic(err)
	}

	for {
		hlfCfgFrom := config.FromFile(conn1)

		err = func() error {
			sdk, err := fabsdk.New(hlfCfgFrom)
			if err != nil {
				err = fmt.Errorf("couldn't initialize SDK: %w", err)
				panic(err)
			}
			defer sdk.Close()

			evnCli := eventClient(sdk, org1, user1, utils.ChannelACL, beginBlock)

			registration, eventChannel, err := evnCli.RegisterBlockEvent()
			if err != nil {
				err = fmt.Errorf("couldn't register chaincode event listener: %w", err)
				panic(err)
			}

			defer func() {
				go func() {
					for range eventChannel {
					}
				}()
				evnCli.Unregister(registration)
			}()

			// Start processing
			evtBuf := make(chan *fab.BlockEvent, sizeOfEvtChan)

			go listen(ctx, eventChannel, evtBuf)

			return rest(ctx, evtBuf, channelClient)
		}()
		if err == nil {
			return
		}

		cnt++
		if cnt >= 3 {
			utils.StdLog.Error(err)
			return
		}
	}
}

func eventClient(
	sdk *fabsdk.FabricSDK,
	org string,
	user string,
	channel string,
	beg uint64,
) (eventClient *event.Client) {
	var err error

	clientChannelContext := sdk.ChannelContext(
		channel,
		fabsdk.WithUser(user),
		fabsdk.WithOrg(org),
	)

	if beg == 1 || beg == 0 {
		eventClient, err = event.New(
			clientChannelContext,
			event.WithBlockEvents(),
			event.WithSeekType(seek.Oldest),
			event.WithEventConsumerTimeout(0),
		)
	} else {
		eventClient, err = event.New(
			clientChannelContext,
			event.WithBlockEvents(),
			event.WithSeekType(seek.FromBlock),
			event.WithBlockNum(beg),
			event.WithEventConsumerTimeout(0),
		)
	}

	if err != nil {
		err = fmt.Errorf("couldn't create event client: %w", err)
		panic(err)
	}

	return
}

func listen(
	ctx context.Context,
	evCh <-chan *fab.BlockEvent,
	evtBuf chan<- *fab.BlockEvent,
) {
	for {
		select {
		case <-ctx.Done():
			return

		case e, ok := <-evCh:
			if !ok {
				utils.StdLog.Info("event channel closed")

				return
			}

			select {
			case evtBuf <- e:
			case <-ctx.Done():
				return
			}
		}
	}
}

func rest(
	ctx context.Context,
	evtBuf <-chan *fab.BlockEvent,
	chClient *channel.Client,
) error {
	var breaking bool

	for {
		select {
		case <-ctx.Done():
			return nil

		case e := <-evtBuf:
			err := RequestBlock(e.Block, chClient)
			if err != nil {
				return err
			}

			if breaking {
				return nil
			}
		}
	}
}

func RequestBlock(b *common.Block, chClient *channel.Client) error {
	txFilter := b.GetMetadata().GetMetadata()[common.BlockMetadataIndex_TRANSACTIONS_FILTER]

	for j, tx := range b.GetData().GetData() {
		if len(txFilter) != 0 {
			if peer.TxValidationCode(txFilter[j]) != peer.TxValidationCode_VALID &&
				peer.TxValidationCode(txFilter[j]) != peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE {
				continue
			}
		}

		envelop, err := blocks.UnmarshalEnvelope(tx)
		if err != nil {
			continue
		}

		payload, err := blocks.ExtractPayload(envelop)
		if err != nil {
			continue
		}

		chHeader, err := blocks.UnmarshalChannelHeader(payload.GetHeader().GetChannelHeader())
		if err != nil {
			continue
		}

		if common.HeaderType(chHeader.GetType()) != common.HeaderType_ENDORSER_TRANSACTION {
			continue
		}

		pTx, err := blocks.GetTransaction(payload.GetData())
		if err != nil {
			continue
		}

		err = TxsProc(pTx.GetActions(), chClient)
		if err != nil {
			return err
		}
	}

	return nil
}
