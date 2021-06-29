package eth

import (
	"context"
)

type Eth interface {
	GetBlocks(ctx context.Context, n int) ([]*Block, error)
	GetBlock(ctx context.Context, id int) (*Block, error)

	GetTransation(ctx context.Context, txHash Hash) (*Transation, error)
}
