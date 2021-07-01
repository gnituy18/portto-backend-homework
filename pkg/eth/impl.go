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

func (im *impl) GetBlocks(ctx context.Context, n uint64) ([]*Block, error) {
	currNum, err := im.GetCurrNum(ctx)
	if err != nil {
		ctx.With(zap.Error(err)).Error("GetCurrNum failed in eth.GetBlocks")
		return nil, err
	}

	// search in db first
	dbBlocks, err := im.getBlocksDB(ctx, currNum, n)
	if err != nil {
		ctx.With(zap.Error(err)).Error("getBlocksDb failed in eth.GetBlocks")
		return nil, err
	}

	numMap := map[uint64]*Block{}
	for _, b := range dbBlocks {
		numMap[b.BlockNum] = b
	}

	blocks := []*Block{}
	for i := 0; i < int(n); i++ {
		blockNum := currNum - n + 1 + uint64(i)
		if b, ok := numMap[blockNum]; ok {
			blocks = append(blocks, b)
			continue
		}

		// inject missing block from rpc
		block, err := im.getBlockByNumberRPC(ctx, blockNum)
		if err != nil {
			ctx.With(zap.Error(err)).Error("goEthClient.BlockByNumber failed in eth.GetBlocks")
			return nil, err
		}

		im.saveBlock(block)
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

func (im *impl) GetTransaction(ctx context.Context, txHash common.Hash) (*Transaction, error) {
	recp, err := im.goEthClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("txHash", txHash.String()),
		).Error("goEthClient.TransactionReceipt failed in eth.GetTransaction")
		return nil, err
	}

	tx, _, err := im.goEthClient.TransactionByHash(ctx, txHash)
	if err != nil {
		ctx.With(
			zap.Error(err),
			zap.String("txHash", txHash.String()),
		).Error("goEthClient.TransactionReceipt failed in eth.GetTransaction")
		return nil, err
	}

	chainID, err := im.goEthClient.NetworkID(ctx)
	if err != nil {
		ctx.With(zap.Error(err)).Error("goEthClient.NetWorkID failed in eth.GetTransaction")
		return nil, err
	}
	msg, err := tx.AsMessage(types.NewEIP155Signer(chainID), nil)
	if err != nil {
		ctx.With(zap.Error(err)).Error("tx.AsMessage failed in eth.GetTransaction")
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

	return &Transaction{
		TxHash: tx.Hash().String(),
		From:   msg.From().Hex(),
		To:     to,
		Nonce:  tx.Nonce(),
		Data:   common.Bytes2Hex(tx.Data()),
		Value:  val,
		Logs:   logs,
	}, nil
}

func (im *impl) SaveTransaction(b *Transaction) error {
	res := im.db.Create(b)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (im *impl) GetCurrNum(ctx context.Context) (uint64, error) {
	currNum, err := im.goEthClient.BlockNumber(ctx)
	if err != nil {
		ctx.With(zap.Error(err)).Error("goEthClient.BlockNumber failed in eth.GetCurrNum")
		return 0, err
	}

	return currNum, nil
}

func (im *impl) getBlocksDB(ctx context.Context, currNum, n uint64) ([]*Block, error) {
	blocks := []*Block{}
	res := im.db.Where("block_num > ?", currNum-n).Order("block_num").Find(&blocks)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	} else if res.Error != nil {
		return nil, res.Error
	}

	return blocks, nil
}

func (im *impl) FetchBlockAndSave(ctx context.Context, blockNum uint64) (bool, error) {
	if _, err := im.getBlockByNumberDB(ctx, blockNum); err == nil {
		return false, nil
	} else if err != nil && err != ErrNotFound {
		ctx.With(zap.Error(err)).Error("getBlockByNumberDB failed")
		return false, err
	}

	block, err := im.getBlockByNumberRPC(ctx, blockNum)
	if err != nil {
		ctx.With(zap.Error(err)).Error("getBlockByNumberRPC failed")
		return false, err
	}

	if err := im.saveBlock(block); err != nil {
		ctx.With(zap.Error(err)).Error("saveBlock failed")
		return false, err
	}

	return true, nil
}

func (im *impl) getBlockByNumberDB(ctx context.Context, blockNum uint64) (*Block, error) {
	block := &Block{}
	res := im.db.First(block, "block_num", blockNum)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	} else if res.Error != nil {
		return nil, res.Error
	}

	return block, nil
}

func (im *impl) getBlockByNumberRPC(ctx context.Context, blockNum uint64) (*Block, error) {
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

func (im *impl) saveBlock(b *Block) error {
	res := im.db.Create(b)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
