package main

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"prottohw/pkg/context"
	"prottohw/pkg/db"
	"prottohw/pkg/eth"
	"prottohw/pkg/log"
)

var (
	rpcEndpoint      = "https://data-seed-prebsc-2-s3.binance.org:8545"
	startingBlockNum = uint64(100000)
	// Modify this num to adjust the worker number
	workerNum = 10
)

func main() {
	pg, err := db.NewPostgres()
	if err != nil {
		panic(err)
	}

	ethclient := eth.New(rpcEndpoint, pg)
	numReader := NewBlockNumReader(ethclient)

	// crawers
	numCh := make(chan uint64)
	for i := 0; i <= workerNum; i++ {
		go worker(numCh, ethclient)
	}

	for blockNum := startingBlockNum; ; blockNum++ {
		n := numReader.GetBlockNum()
		if err != nil {
			log.Global().Error("ethclient.GetBlockNum failed")
		}

		// wait 1 min if no new blocks to curl
		if n == uint64(blockNum) {
			time.Sleep(time.Minute)
			continue
		}

		numCh <- blockNum
	}
}

func worker(numCh <-chan uint64, ethclient eth.Eth) {
	for n := range numCh {
		block, err := ethclient.GetBlockByNumber(context.Background(), n)
		if err != nil {
			log.Global().Error("block hash:" + block.BlockHash)
			continue
		}

		log.Global().With(zap.String("block hash", block.BlockHash)).Info("GetBlockByNumber")
	}
}

var (
	updateBlockNumInterval = time.Minute
)

// BlockNumReader try to minimize RPC call by getting latest block number once per minute
func NewBlockNumReader(ethclient eth.Eth) *BlockNumReader {
	reader := &BlockNumReader{}

	go func() {
		c := time.Tick(updateBlockNumInterval)
		for range c {
			n, err := ethclient.GetBlockNum(context.Background())
			if err != nil {
				log.Global().Error("ethclient.GetBlockNum failed")
				continue
			}

			reader.mux.Lock()
			log.Global().Error("update curr block num")
			defer reader.mux.Unlock()
			reader.currBlockNum = n

		}
	}()

	return reader
}

type BlockNumReader struct {
	currBlockNum uint64
	mux          sync.RWMutex
}

func (r *BlockNumReader) GetBlockNum() uint64 {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.currBlockNum
}
