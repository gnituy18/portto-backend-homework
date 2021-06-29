package eth

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	"prottohw/pkg/context"
)

func New(url string) Eth {
	goEthClient, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}

	return &impl{
		goEthClient: goEthClient,
		url:         url,
	}
}

type impl struct {
	goEthClient *ethclient.Client
	url         string
}

func (im *impl) GetBlocks(ctx context.Context, n uint64) ([]*Block, error) {
	currNum, err := im.goEthClient.BlockNumber(ctx)
	if err != nil {
		ctx.With(zap.Error(err)).Error("goEthClient.BlockNumber failed in eth.GetBlocks")
		return nil, err
	}

	blocks := []*Block{}
	for i := 0; i < int(n); i++ {
		blockNum := big.NewInt(int64(currNum-n+1) + int64(i))
		block, err := im.goEthClient.BlockByNumber(ctx, blockNum)
		if err != nil {
			ctx.With(zap.Error(err)).Error("goEthClient.BlockByNumber failed in eth.GetBlocks")
			return nil, err
		}

		blocks = append(blocks, &Block{
			BlockNum:   block.NumberU64(),
			BlockHash:  block.Hash().String(),
			BlockTime:  block.Time(),
			ParentHash: block.ParentHash().String(),
		})
	}

	ctx.Infof("%v", blocks)
	return blocks, nil
}

func (im *impl) GetBlock(ctx context.Context, hash common.Hash) (*Block, error) {
	block, err := im.goEthClient.BlockByHash(ctx, hash)
	if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("hash", hash.String()),
		).Error("goEthClient.BlockByNumber failed in eth.GetBlocks")
		return nil, err
	}

	txHashs := []string{}
	for _, tx := range block.Transactions() {
		txHashs = append(txHashs, tx.Hash().String())
	}

	return &Block{
		BlockNum:     block.NumberU64(),
		BlockHash:    block.Hash().String(),
		BlockTime:    block.Time(),
		ParentHash:   block.ParentHash().String(),
		Transactions: txHashs,
	}, nil
}

func (im *impl) GetTransation(ctx context.Context, txHash common.Hash) (*Transation, error) {
	recp, err := im.goEthClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("txHash", txHash.String()),
		).Error("goEthClient.TransactionReceipt failed in eth.GetTransation")
		return nil, err
	}

	logs := []*Log{}
	for _, log := range recp.Logs {
		logs = append(logs, &Log{
			Index: log.Index,
			Data:  string(log.Data),
		})
	}

	tx, _, err := im.goEthClient.TransactionByHash(ctx, txHash)
	if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("txHash", txHash.String()),
		).Error("goEthClient.TransactionReceipt failed in eth.GetTransation")
		return nil, err
	}

	chainID, err := im.goEthClient.NetworkID(ctx)
	if err != nil {
		ctx.With(zap.Error(err)).Error("goEthClient.NetWorkID failed in eth.GetTransation")
		return nil, err
	}

	msg, err := tx.AsMessage(types.NewEIP155Signer(chainID), nil)
	if err != nil {
		ctx.With(zap.Error(err)).Error("tx.AsMessage failed in eth.GetTransation")
		return nil, err
	}

	return &Transation{
		TxHash: tx.Hash().String(),
		From:   msg.From().Hex(),
		To:     tx.To().String(),
		Nonce:  tx.Nonce(),
		Data:   string(tx.Data()),
		Value:  tx.Value().String(),
		Logs:   logs,
	}, nil
}
