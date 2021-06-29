package eth

import (
	"prottohw/pkg/context"

	"github.com/ethereum/go-ethereum/common"
)

type Eth interface {
	GetBlocks(ctx context.Context, n uint64) ([]*Block, error)
	GetBlock(ctx context.Context, id uint64) (*Block, error)

	GetTransation(ctx context.Context, txHash common.Hash) (*Transation, error)
}
