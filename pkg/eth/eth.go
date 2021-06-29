package eth

import (
	"context"
)

type Eth interface {
	GetBlocks(ctx context.Context, n int64) ([]*Block, error)
	GetBlock(ctx context.Context, id int64) (*Block, error)

	GetTransation(ctx context.Context, txHash Hash) (*Transation, error)
}
