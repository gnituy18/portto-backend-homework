package eth

import (
	"context"
	"errors"
)

func New(url string) Eth {
	return &impl{}
}

type impl struct {
	url string
}

func (im *impl) GetBlocks(ctx context.Context, n int) ([]*Block, error) {
	return nil, errors.New("TODO")
}

func (im *impl) GetBlock(ctx context.Context, id int) (*Block, error) {
	return nil, errors.New("TODO")
}

func (im *impl) GetTransation(ctx context.Context, txHash Hash) (*Transation, error) {
	return nil, errors.New("TODO")
}
