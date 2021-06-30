package main

import (
	"fmt"
	"time"

	"prottohw/pkg/context"
	"prottohw/pkg/db"
	"prottohw/pkg/eth"
	"prottohw/pkg/log"
)

var (
	rpcEndpoint      = "https://data-seed-prebsc-2-s3.binance.org:8545"
	startingBlockNum = 100000
	// Modify this num to adjust the
	workerNum = 10
)

func main() {
	ctx := context.Background()

	// db conn
	_, err := db.NewPostgres()
	if err != nil {
		panic(err)
	}

	// eth client
	ethclient := eth.New(rpcEndpoint)

	// init worker
	numCh := make(chan int)
	for i := 0; i <= workerNum; i++ {
		go craw(numCh)
	}

	for blockNum := startingBlockNum; ; blockNum++ {
		n, err := ethclient.GetBlockNum(ctx)
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

func craw(numCh <-chan int) {
	for n := range numCh {
		time.Sleep(time.Second)
		fmt.Println(n)
	}
}
