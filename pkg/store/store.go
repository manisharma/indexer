package store

import (
	"context"
	"fmt"
	"indexer/pkg/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	Create(context.Context, models.Epoch) error
	Get(context.Context) ([]models.Epoch, error)
	KeepOnlyTop5(context.Context, uint64) error
}

type Store struct {
	pool    *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

// Store implements Repository
var _ Repository = &Store{}

func New(pool *pgxpool.Pool) *Store {
	return &Store{
		pool:    pool,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (s *Store) Create(ctx context.Context, e models.Epoch) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	success := false
	defer func() {
		if success {
			tx.Commit(ctx)

		} else {
			tx.Rollback(ctx)
		}
	}()

	// insert epoch
	qry, args, err := s.builder.Insert("epochs").
		Columns("epoch_number", "start_time", "end_time").
		Values(e.EpochNumber, e.StartTime, e.EndTime).ToSql()
	if err != nil {
		return fmt.Errorf("epochs insert query prep failed, err: %v", err.Error())
	}
	_, err = tx.Exec(ctx, qry, args...)
	if err != nil {
		return fmt.Errorf("epochs insert query failed, err: %v", err.Error())
	}

	// insert slots & blocks
	slotsBldr := s.builder.Insert("slots").
		Columns("slot_number", "start_time", "end_time", "epoch_number")
	blocksBldr := s.builder.Insert("blocks").
		Columns("block_number", "block_root", "state_root", "slot_number", "gas_limit", "gas_used", "no_of_transactions", "created_at")
	for _, s := range e.Slots {
		slotsBldr = slotsBldr.Values(s.SlotNumber, s.StartTime, s.EndTime, s.EpochNumber)
		blocksBldr = blocksBldr.Values(s.Block.BlockNumber, s.Block.BlockRoot, s.Block.StateRoot,
			s.Block.SlotNumber, s.Block.GasLimit, s.Block.GasUsed, s.Block.NoOfTransactions, s.Block.CreatedAt)
	}
	qry, args, err = slotsBldr.ToSql()
	if err != nil {
		return fmt.Errorf("slots insert query prep failed, err: %v", err.Error())
	}
	_, err = tx.Exec(ctx, qry, args...)
	if err != nil {
		return fmt.Errorf("slots insert query failed, err: %v", err.Error())
	}

	qry, args, err = blocksBldr.ToSql()
	if err != nil {
		return fmt.Errorf("blocks insert query prep failed, err: %v", err.Error())
	}
	_, err = tx.Exec(ctx, qry, args...)
	if err != nil {
		return fmt.Errorf("blocks insert query failed, err: %v", err.Error())
	}

	success = true
	return nil
}

func (s *Store) Get(ctx context.Context) ([]models.Epoch, error) {
	qry, args, err := s.builder.Select("*").From("epochs").ToSql()
	if err != nil {
		return nil, fmt.Errorf("epochs select query prep failed, err: %v", err.Error())
	}
	var epochs []models.Epoch
	err = pgxscan.Select(ctx, s.pool, &epochs, qry, args...)
	if err != nil {
		return nil, fmt.Errorf("epochs select query failed, err: %v", err.Error())
	}

	qry, args, err = s.builder.Select("*").From("slots").ToSql()
	if err != nil {
		return nil, fmt.Errorf("slots select query prep failed, err: %v", err.Error())
	}
	var slots []models.Slot
	err = pgxscan.Select(ctx, s.pool, &slots, qry, args...)
	if err != nil {
		return nil, fmt.Errorf("slots select query failed, err: %v", err.Error())
	}

	qry, args, err = s.builder.Select("*").From("blocks").ToSql()
	if err != nil {
		return nil, fmt.Errorf("blocks select query prep failed, err: %v", err.Error())
	}
	var blocks []models.Block
	err = pgxscan.Select(ctx, s.pool, &blocks, qry, args...)
	if err != nil {
		return nil, fmt.Errorf("blocks select query failed, err: %v", err.Error())
	}

	for idxEpoch := range epochs {
		for idxSlot, slot := range slots {
			if slot.EpochNumber == epochs[idxEpoch].EpochNumber {
				for _, block := range blocks {
					if slot.SlotNumber == block.SlotNumber {
						slots[idxSlot].Block = block
					}
				}
				epochs[idxEpoch].Slots = append(epochs[idxEpoch].Slots, slots[idxSlot])
			}
		}
	}

	return epochs, nil
}

func (s *Store) KeepOnlyTop5(ctx context.Context, epochNumber uint64) error {
	qry, args, err := s.builder.Delete("epochs").
		Where(squirrel.LtOrEq{"epoch_number": epochNumber - 5}).
		ToSql()
	if err != nil {
		return fmt.Errorf("epochs delete query prep failed, err: %v", err.Error())
	}
	_, err = s.pool.Exec(ctx, qry, args...)
	if err != nil {
		return fmt.Errorf("epochs delete query failed, err: %v", err.Error())
	}
	return err
}
