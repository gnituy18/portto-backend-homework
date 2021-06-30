package main

import (
	"sync"
	"time"

	"prottohw/pkg/context"
	"prottohw/pkg/eth"
	"prottohw/pkg/log"
)

var (
	updateBlockNumInterval = time.Minute
)

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
