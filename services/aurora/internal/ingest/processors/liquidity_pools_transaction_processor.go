package processors

import (
	"context"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/toid"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

type liquidityPool struct {
	internalID     int64 // Bigint auto-generated by postgres
	transactionSet map[int64]struct{}
	operationSet   map[int64]struct{}
}

func (b *liquidityPool) addTransactionID(id int64) {
	if b.transactionSet == nil {
		b.transactionSet = map[int64]struct{}{}
	}
	b.transactionSet[id] = struct{}{}
}

func (b *liquidityPool) addOperationID(id int64) {
	if b.operationSet == nil {
		b.operationSet = map[int64]struct{}{}
	}
	b.operationSet[id] = struct{}{}
}

type LiquidityPoolsTransactionProcessor struct {
	sequence         uint32
	liquidityPoolSet map[string]liquidityPool
	qLiquidityPools  history.QHistoryLiquidityPools
}

func NewLiquidityPoolsTransactionProcessor(Q history.QHistoryLiquidityPools, sequence uint32) *LiquidityPoolsTransactionProcessor {
	return &LiquidityPoolsTransactionProcessor{
		qLiquidityPools:  Q,
		sequence:         sequence,
		liquidityPoolSet: map[string]liquidityPool{},
	}
}

func (p *LiquidityPoolsTransactionProcessor) ProcessTransaction(ctx context.Context, transaction ingest.LedgerTransaction) error {
	err := p.addTransactionLiquidityPools(p.liquidityPoolSet, p.sequence, transaction)
	if err != nil {
		return err
	}

	err = p.addOperationLiquidityPools(p.liquidityPoolSet, p.sequence, transaction)
	if err != nil {
		return err
	}

	return nil
}

func (p *LiquidityPoolsTransactionProcessor) addTransactionLiquidityPools(lpSet map[string]liquidityPool, sequence uint32, transaction ingest.LedgerTransaction) error {
	transactionID := toid.New(int32(sequence), int32(transaction.Index), 0).ToInt64()
	transactionLiquidityPools, err := liquidityPoolsForTransaction(
		sequence,
		transaction,
	)
	if err != nil {
		return errors.Wrap(err, "Could not determine liquidity pools for transaction")
	}

	for _, lp := range transactionLiquidityPools {
		entry := lpSet[lp]
		entry.addTransactionID(transactionID)
		lpSet[lp] = entry
	}

	return nil
}

func liquidityPoolsForTransaction(
	sequence uint32,
	transaction ingest.LedgerTransaction,
) ([]string, error) {
	changes, err := transaction.GetChanges()
	if err != nil {
		return nil, err
	}
	lps, err := liquidityPoolsForChanges(changes)
	if err != nil {
		return nil, errors.Wrapf(err, "reading transaction %v liquidity pools", transaction.Index)
	}
	return dedupeLiquidityPools(lps)
}

func dedupeLiquidityPools(in []string) (out []string, err error) {
	set := map[string]struct{}{}
	for _, id := range in {
		set[id] = struct{}{}
	}

	for id := range set {
		out = append(out, id)
	}
	return
}

func liquidityPoolsForChanges(
	changes []ingest.Change,
) ([]string, error) {
	var lps []string

	for _, c := range changes {
		if c.Type != xdr.LedgerEntryTypeLiquidityPool {
			continue
		}

		if c.Pre == nil && c.Post == nil {
			return nil, errors.New("Invalid io.Change: change.Pre == nil && change.Post == nil")
		}

		if c.Pre != nil {
			poolID := c.Pre.Data.MustLiquidityPool().LiquidityPoolId
			lps = append(lps, PoolIDToString(poolID))
		}
		if c.Post != nil {
			poolID := c.Post.Data.MustLiquidityPool().LiquidityPoolId
			lps = append(lps, PoolIDToString(poolID))
		}
	}

	return lps, nil
}

func (p *LiquidityPoolsTransactionProcessor) addOperationLiquidityPools(lpSet map[string]liquidityPool, sequence uint32, transaction ingest.LedgerTransaction) error {
	liquidityPools, err := liquidityPoolsForOperations(transaction, sequence)
	if err != nil {
		return errors.Wrap(err, "could not determine operation liquidity pools")
	}

	for operationID, lps := range liquidityPools {
		for _, lp := range lps {
			entry := lpSet[lp]
			entry.addOperationID(operationID)
			lpSet[lp] = entry
		}
	}

	return nil
}

func liquidityPoolsForOperations(transaction ingest.LedgerTransaction, sequence uint32) (map[int64][]string, error) {
	lps := map[int64][]string{}

	for opi, op := range transaction.Envelope.Operations() {
		operation := transactionOperationWrapper{
			index:          uint32(opi),
			transaction:    transaction,
			operation:      op,
			ledgerSequence: sequence,
		}

		changes, err := transaction.GetOperationChanges(uint32(opi))
		if err != nil {
			return lps, err
		}
		c, err := liquidityPoolsForChanges(changes)
		if err != nil {
			return lps, errors.Wrapf(err, "reading operation %v liquidity pools", operation.ID())
		}
		lps[operation.ID()] = c
	}

	return lps, nil
}

func (p *LiquidityPoolsTransactionProcessor) Commit(ctx context.Context) error {
	if len(p.liquidityPoolSet) > 0 {
		if err := p.loadLiquidityPoolIDs(ctx, p.liquidityPoolSet); err != nil {
			return err
		}

		if err := p.insertDBTransactionLiquidityPools(ctx, p.liquidityPoolSet); err != nil {
			return err
		}

		if err := p.insertDBOperationsLiquidityPools(ctx, p.liquidityPoolSet); err != nil {
			return err
		}
	}

	return nil
}

func (p *LiquidityPoolsTransactionProcessor) loadLiquidityPoolIDs(ctx context.Context, liquidityPoolSet map[string]liquidityPool) error {
	ids := make([]string, 0, len(liquidityPoolSet))
	for id := range liquidityPoolSet {
		ids = append(ids, id)
	}

	toInternalID, err := p.qLiquidityPools.CreateHistoryLiquidityPools(ctx, ids, maxBatchSize)
	if err != nil {
		return errors.Wrap(err, "Could not create liquidity pool ids")
	}

	for _, id := range ids {
		internalID, ok := toInternalID[id]
		if !ok {
			return errors.Errorf("no internal id found for liquidity pool %s", id)
		}

		lp := liquidityPoolSet[id]
		lp.internalID = internalID
		liquidityPoolSet[id] = lp
	}

	return nil
}

func (p LiquidityPoolsTransactionProcessor) insertDBTransactionLiquidityPools(ctx context.Context, liquidityPoolSet map[string]liquidityPool) error {
	batch := p.qLiquidityPools.NewTransactionLiquidityPoolBatchInsertBuilder(maxBatchSize)

	for _, entry := range liquidityPoolSet {
		for transactionID := range entry.transactionSet {
			if err := batch.Add(ctx, transactionID, entry.internalID); err != nil {
				return errors.Wrap(err, "could not insert transaction liquidity pool in db")
			}
		}
	}

	if err := batch.Exec(ctx); err != nil {
		return errors.Wrap(err, "could not flush transaction liquidity pools to db")
	}
	return nil
}

func (p LiquidityPoolsTransactionProcessor) insertDBOperationsLiquidityPools(ctx context.Context, liquidityPoolSet map[string]liquidityPool) error {
	batch := p.qLiquidityPools.NewOperationLiquidityPoolBatchInsertBuilder(maxBatchSize)

	for _, entry := range liquidityPoolSet {
		for operationID := range entry.operationSet {
			if err := batch.Add(ctx, operationID, entry.internalID); err != nil {
				return errors.Wrap(err, "could not insert operation liquidity pool in db")
			}
		}
	}

	if err := batch.Exec(ctx); err != nil {
		return errors.Wrap(err, "could not flush operation liquidity pools to db")
	}
	return nil
}
