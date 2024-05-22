package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// EthereumParser implements Parser
type EthereumParser struct {
	lock         sync.Mutex
	url          string
	currentBlock int
	subscribers  map[string]struct{}
	transactions map[string][]Transaction
}

type Block struct {
	Result Transactions `json:"result"`
}

type Transactions struct {
	Transaction []Transaction `json:"transactions"`
}

type Transaction struct {
	Hash  string `json:"hash"` // transaction hash
	From  string `json:"from"` // sender address
	To    string `json:"to"`   // receiver address
	Value string `json:"value"`
}

func (ep *EthereumParser) GetCurrentBlock() int {
	return ep.currentBlock
}

func (ep *EthereumParser) Subscribe(address string) bool {
	ep.lock.Lock()
	defer ep.lock.Unlock()
	// Check if the address has been subscribed
	_, exists := ep.subscribers[address]
	if !exists {
		ep.subscribers[address] = struct{}{}
		ep.transactions[address] = []Transaction{}
	}
	return true
}

// GetTransactions Get transactions of subscribed users
func (ep *EthereumParser) GetTransactions(address string) []Transaction {
	return ep.transactions[address]
}

// GetBlockTransactions Get all transactions of the Block
func (ep *EthereumParser) GetBlockTransactions(blockNumber string) ([]string, error) {
	payload := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "eth_getBlockByNumber",
		"params": ["%s", true],
		"id": 1
	}`, blockNumber)
	resp, err := http.Post(ep.url, "application/json", strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Block
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var txHashes []string
	ep.lock.Lock()
	defer ep.lock.Unlock()
	for _, tx := range result.Result.Transaction {
		fmt.Println(tx)
		if _, ok := ep.subscribers[tx.From]; ok {
			ep.transactions[tx.From] = append(ep.transactions[tx.From], tx)
		}
		if _, ok := ep.subscribers[tx.To]; ok {
			ep.transactions[tx.To] = append(ep.transactions[tx.To], tx)
		}
	}

	return txHashes, nil
}

func (ep *EthereumParser) GetCurrentBlockNumber() (string, error) {
	payload := `{
		"jsonrpc": "2.0",
		"method": "eth_blockNumber",
		"params": [],
		"id": 1
	}`
	resp, err := http.Post(ep.url, "application/json", strings.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	blockNumber, ok := result["result"].(string)
	if !ok {
		return "", fmt.Errorf("invalid block number format")
	}
	if blockInt, err := strconv.ParseInt(blockNumber, 16, 64); err == nil {
		ep.currentBlock = int(blockInt)
	}
	return blockNumber, nil
}
