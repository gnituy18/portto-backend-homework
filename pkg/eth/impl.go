package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"prottohw/pkg/context"
)

func New(url string, db *gorm.DB) Eth {
	goEthClient, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}

	return &impl{
		goEthClient: goEthClient,
		db:          db,
		url:         url,
	}
}

type impl struct {
	goEthClient *ethclient.Client
	url         string
	db          *gorm.DB
}

func (im *impl) GetBlockNum(ctx context.Context) (uint64, error) {
	currNum, err := im.goEthClient.BlockNumber(ctx)
	if err != nil {
		ctx.With(zap.Error(err)).Error("goEthClient.BlockNumber failed in eth.GetBlockNum")
		return 0, err
	}

	return currNum, nil
}

func (im *impl) GetBlockByNumber(ctx context.Context, blockNum uint64) (*Block, error) {
	blockNumBig := big.NewInt(int64(blockNum))
	block, err := im.goEthClient.BlockByNumber(ctx, blockNumBig)
	if err != nil {
		ctx.With(zap.Error(err)).Error("goEthClient.BlockByNumber failed in eth.GetBlocks")
		return nil, err
	}

	return &Block{
		BlockNum:   block.NumberU64(),
		BlockHash:  block.Hash().String(),
		BlockTime:  block.Time(),
		ParentHash: block.ParentHash().String(),
	}, nil
}

func (im *impl) GetBlocks(ctx context.Context, n uint64) ([]*Block, error) {
	currNum, err := im.GetBlockNum(ctx)
	if err != nil {
		ctx.With(zap.Error(err)).Error("GetBlockNum failed in eth.GetBlocks")
		return nil, err
	}

	blocks := []*Block{}
	for i := 0; i < int(n); i++ {
		blockNum := currNum - n + 1 + uint64(i)
		block, err := im.GetBlockByNumber(ctx, blockNum)
		if err != nil {
			ctx.With(zap.Error(err)).Error("goEthClient.BlockByNumber failed in eth.GetBlocks")
			return nil, err
		}

		blocks = append(blocks, block)
	}

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

	var val int64
	if tx.Value() != nil {
		val = tx.Value().Int64()
	}
	var to string
	if tx.To() != nil {
		to = tx.To().String()
	}
	logs := []*Log{}
	for _, log := range recp.Logs {
		logs = append(logs, &Log{
			Index: log.Index,
			Data:  common.Bytes2Hex(log.Data),
		})
	}

	return &Transation{
		TxHash: tx.Hash().String(),
		From:   msg.From().Hex(),
		To:     to,
		Nonce:  tx.Nonce(),
		Data:   common.Bytes2Hex(tx.Data()),
		Value:  val,
		Logs:   logs,
	}, nil
}
