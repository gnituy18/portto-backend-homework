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
	GetCurrNum(ctx context.Context) (uint64, error)
	GetBlocks(ctx context.Context, n uint64) ([]*Block, error)
	GetBlock(ctx context.Context, hash common.Hash) (*Block, error)

	GetBlockByNumberRPC(ctx context.Context, n uint64) (*Block, error)
	GetBlockByNumberDB(ctx context.Context, n uint64) (*Block, error)
	SaveBlock(b *Block) error

	GetTransaction(ctx context.Context, txHash common.Hash) (*Transaction, error)
}
