package indexer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"indexer/pkg/models"
	"log"
	"time"

	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/rs/zerolog"
)

type BeaconChain struct {
	httpClient *http.Service
}

type EpochResult struct {
	Epoch *models.Epoch
	Error error
}

// Creates new instance of Ethereum http client service
func New(ctx context.Context, clientURL string) (*BeaconChain, error) {
	client, err := http.New(ctx, http.WithAddress(clientURL), http.WithLogLevel(zerolog.ErrorLevel))
	if err != nil {
		return nil, err
	}
	httpClient, ok := client.(*http.Service)
	if !ok {
		return nil, errors.New("invalid ethereum client")
	}
	return &BeaconChain{httpClient}, nil
}

func (b *BeaconChain) SubscribeToEpochs(ctx context.Context) <-chan EpochResult {

	var (
		slotPerEpoch     uint64
		slotDuration     time.Duration
		epochDuration    time.Duration
		epochStream      chan EpochResult = make(chan EpochResult)
		lastEpoch        uint64
		slots            []models.Slot = make([]models.Slot, 0, slotPerEpoch)
		err              error
		closeEpochStream bool
	)

	defer func() {
		if closeEpochStream {
			close(epochStream)
		}
	}()

	go func() {
		slotPerEpoch, err = b.httpClient.SlotsPerEpoch(ctx)
		if err != nil {
			closeEpochStream = true
			epochStream <- EpochResult{
				Epoch: nil,
				Error: fmt.Errorf("could not find SlotsPerEpoch, err: %v", err.Error()),
			}
			return
		}

		slotDuration, err = b.httpClient.SlotDuration(ctx)
		if err != nil {
			closeEpochStream = true
			epochStream <- EpochResult{
				Epoch: nil,
				Error: fmt.Errorf("could not find SlotDuration, err: %v", err.Error()),
			}
			return
		}

		epochDurStr := fmt.Sprintf("%ds", slotPerEpoch*uint64(slotDuration.Seconds()))
		epochDuration, err = time.ParseDuration(epochDurStr)
		if err != nil {
			closeEpochStream = true
			epochStream <- EpochResult{
				Epoch: nil,
				Error: fmt.Errorf("could not parse EpochDuration, err: %v", err.Error()),
			}
		}

		// subscribe to block event
		err = b.httpClient.Events(ctx, []string{"block"}, func(e *v1.Event) {
			// respect cancellation/unsubscription
			select {
			case <-ctx.Done():
				close(epochStream)
			default:
			}
			evtBytes, err := json.Marshal(e.Data)
			if err != nil {
				epochStream <- EpochResult{
					Epoch: nil,
					Error: fmt.Errorf("json.Marshal(e.Data) failed, err %v", err.Error()),
				}
				return
			}
			var blockEvent *v1.BlockEvent = &v1.BlockEvent{}
			err = blockEvent.UnmarshalJSON(evtBytes)
			if err != nil {
				epochStream <- EpochResult{
					Epoch: nil,
					Error: fmt.Errorf("blockEvent.UnmarshalJSON(evtBytes) failed, err %v", err.Error()),
				}
				return
			}

			epoch := uint64(blockEvent.Slot) / 32
			if lastEpoch == 0 {
				lastEpoch = epoch
			}

			aSlot := models.Slot{
				SlotNumber: uint64(blockEvent.Slot),
			}
			aBlock := models.Block{
				BlockRoot: blockEvent.Block.String(),
			}

			// if a signed beacon block for the block ID is not available this will return nil without an error.
			block, err := b.httpClient.SignedBeaconBlock(ctx, blockEvent.Block.String())
			if err != nil {
				epochStream <- EpochResult{
					Epoch: nil,
					Error: fmt.Errorf("httpClient.SignedBeaconBlock() failed, err: %v", err.Error()),
				}
				return
			}
			if block == nil {
				epochStream <- EpochResult{
					Epoch: nil,
					Error: fmt.Errorf("no signed beacon block for the block ID (%s)", blockEvent.Block.String()),
				}
				return
			}

			switch block.Version {
			case spec.DataVersionBellatrix:
				aBlock.BlockNumber = block.Bellatrix.Message.Body.ExecutionPayload.BlockNumber
				aBlock.GasLimit = block.Bellatrix.Message.Body.ExecutionPayload.GasLimit
				aBlock.GasUsed = block.Bellatrix.Message.Body.ExecutionPayload.GasUsed
				aBlock.NoOfTransactions = len(block.Bellatrix.Message.Body.ExecutionPayload.Transactions)
				stateRoot, err := block.StateRoot()
				if err != nil {
					epochStream <- EpochResult{
						Epoch: nil,
						Error: fmt.Errorf("block.StateRoot() failed, err: %v", err.Error()),
					}
					return
				}
				aBlock.StateRoot = stateRoot.String()
				aBlock.CreatedAt = time.Unix(int64(block.Bellatrix.Message.Body.ExecutionPayload.Timestamp), 0)
			case spec.DataVersionCapella:
				aBlock.BlockNumber = block.Capella.Message.Body.ExecutionPayload.BlockNumber
				aBlock.GasLimit = block.Capella.Message.Body.ExecutionPayload.GasLimit
				aBlock.GasUsed = block.Capella.Message.Body.ExecutionPayload.GasUsed
				aBlock.NoOfTransactions = len(block.Capella.Message.Body.ExecutionPayload.Transactions)
				stateRoot, err := block.StateRoot()
				if err != nil {
					epochStream <- EpochResult{
						Epoch: nil,
						Error: fmt.Errorf("block.StateRoot() failed, err: %v", err.Error()),
					}
					return
				}
				aBlock.StateRoot = stateRoot.String()
				aBlock.CreatedAt = time.Unix(int64(block.Capella.Message.Body.ExecutionPayload.Timestamp), 0)
			case spec.DataVersionDeneb:
				aBlock.BlockNumber = block.Deneb.Message.Body.ExecutionPayload.BlockNumber
				aBlock.GasLimit = block.Deneb.Message.Body.ExecutionPayload.GasLimit
				aBlock.GasUsed = block.Deneb.Message.Body.ExecutionPayload.GasUsed
				aBlock.NoOfTransactions = len(block.Deneb.Message.Body.ExecutionPayload.Transactions)
				stateRoot, err := block.StateRoot()
				if err != nil {
					epochStream <- EpochResult{
						Epoch: nil,
						Error: fmt.Errorf("block.StateRoot() failed, err: %v", err.Error()),
					}
					return
				}
				aBlock.StateRoot = stateRoot.String()
				aBlock.CreatedAt = time.Unix(int64(block.Deneb.Message.Body.ExecutionPayload.Timestamp), 0)
			}

			aBlock.SlotNumber = aSlot.SlotNumber
			aSlot.StartTime = aBlock.CreatedAt
			aSlot.EndTime = aSlot.StartTime.Add(slotDuration)
			aSlot.Block = aBlock
			aSlot.EpochNumber = lastEpoch
			slots = append(slots, aSlot)

			if lastEpoch != epoch && len(slots) > 0 {
				anEpoch := models.Epoch{
					EpochNumber: lastEpoch,
					Slots:       slots,
				}
				anEpoch.EpochNumber = lastEpoch
				anEpoch.StartTime = aBlock.CreatedAt
				anEpoch.StartTime = anEpoch.Slots[0].Block.CreatedAt
				anEpoch.EndTime = anEpoch.StartTime.Add(epochDuration)
				log.Println("new epoch", anEpoch.EpochNumber)
				slots = []models.Slot{}
				lastEpoch = epoch
				epochStream <- EpochResult{
					Epoch: &anEpoch,
					Error: nil,
				}
			}
		})
		if err != nil {
			closeEpochStream = true
			epochStream <- EpochResult{
				Epoch: nil,
				Error: fmt.Errorf("block subscription failed, err %v", err.Error()),
			}
		}
	}()

	return epochStream
}
