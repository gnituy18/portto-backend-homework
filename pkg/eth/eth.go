package eth

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"prottohw/pkg/context"
)

var (
	ErrNotFound = ethereum.NotFound
)

type Eth interface {
	GetBlockNum(ctx context.Context) (uint64, error)
	GetBlock(ctx context.Context, hash common.Hash) (*Block, error)
	GetBlockByNumber(ctx context.Context, n uint64) (*Block, error)
	GetBlocks(ctx context.Context, n uint64) ([]*Block, error)

	GetTransation(ctx context.Context, txHash common.Hash) (*Transation, error)
}
