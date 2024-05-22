package main

import (
	"fmt"
	"sync"
	"time"
)

// 1.err处理 2.注释 3.readme
func main() {
	ethParser := EthereumParser{
		url:          "https://cloudflare-eth.com",
		lock:         sync.Mutex{},
		subscribers:  make(map[string]struct{}),
		transactions: make(map[string][]Transaction),
	}

	// ETH produces a block per 12 seconds. We can also subscribe to the event interface of the consensus layer here.
	ticker := time.NewTicker(12 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			blockNumber, err := ethParser.GetCurrentBlockNumber()
			if err != nil {
				fmt.Printf("Error getting current block number: %v\n", err)
				continue
			}
			fmt.Printf("Current block number: %s\n", blockNumber)

			txs, err := ethParser.GetBlockTransactions(blockNumber)
			if err != nil {
				fmt.Printf("Error getting transactions for block %s: %v\n", blockNumber, err)
				continue
			}
			fmt.Printf("Transactions in block %s: %v\n", blockNumber, txs)
		}
	}
}
