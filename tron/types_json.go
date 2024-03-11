package tron

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// JSONResponse is the base struct for JSONRPC2.0, working like ethereum's
type JSONResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
}

type JSONLatestBlock struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

// JSONBlockWithTxs represents a Block in which Transactions field is a list of Transaction structs
type JSONBlockWithTxs struct {
	BaseFeePerGas    interface{}       `json:"baseFeePerGas"`
	Difficulty       interface{}       `json:"difficulty"`
	ExtraData        interface{}       `json:"extraData"`
	GasLimit         *hexutil.Big      `json:"gasLimit"`
	GasUsed          *hexutil.Big      `json:"gasUsed"`
	Hash             string            `json:"hash"`
	LogsBloom        string            `json:"logsBloom"`
	Miner            string            `json:"miner"`
	MixHash          interface{}       `json:"mixHash"`
	Nonce            interface{}       `json:"nonce"` // null
	Number           *hexutil.Uint64   `json:"number"`
	ParentHash       string            `json:"parentHash"`
	ReceiptsRoot     interface{}       `json:"receiptsRoot"`
	Sha3Uncles       interface{}       `json:"sha3Uncles"`
	Size             *hexutil.Uint64   `json:"size"`
	StateRoot        string            `json:"stateRoot"`
	Timestamp        *hexutil.Uint64   `json:"timestamp"`
	TotalDifficulty  interface{}       `json:"totalDifficulty"`
	Transactions     []JSONTransaction `json:"transactions"`
	TransactionsRoot string            `json:"transactionsRoot"`
	Uncles           []string          `json:"uncles"`
}

// JSONBlockWithTxIDs represents a Block in which Transactions field is a list of Transaction IDs
type JSONBlockWithTxIDs struct {
	BaseFeePerGas    interface{}     `json:"baseFeePerGas"`
	Difficulty       interface{}     `json:"difficulty"`
	ExtraData        interface{}     `json:"extraData"`
	GasLimit         *hexutil.Big    `json:"gasLimit"`
	GasUsed          *hexutil.Big    `json:"gasUsed"`
	Hash             string          `json:"hash"`
	LogsBloom        string          `json:"logsBloom"`
	Miner            string          `json:"miner"`
	MixHash          interface{}     `json:"mixHash"`
	Nonce            interface{}     `json:"nonce"` // null
	Number           *hexutil.Uint64 `json:"number"`
	ParentHash       string          `json:"parentHash"`
	ReceiptsRoot     interface{}     `json:"receiptsRoot"`
	Sha3Uncles       interface{}     `json:"sha3Uncles"`
	Size             *hexutil.Uint64 `json:"size"`
	StateRoot        string          `json:"stateRoot"`
	Timestamp        *hexutil.Uint64 `json:"timestamp"`
	TotalDifficulty  interface{}     `json:"totalDifficulty"`
	Transactions     []string        `json:"transactions"`
	TransactionsRoot string          `json:"transactionsRoot"`
	Uncles           []string        `json:"uncles"`
}

// JSONTransaction represents a Transaction
type JSONTransaction struct {
	BlockHash        string          `json:"blockHash"`
	BlockNumber      *hexutil.Uint64 `json:"blockNumber"`
	From             string          `json:"from"`
	Gas              *hexutil.Big    `json:"gas"`
	GasPrice         *hexutil.Big    `json:"gasPrice"`
	Hash             string          `json:"hash"`
	Input            string          `json:"input"`
	Nonce            interface{}     `json:"nonce"` // always null
	R                string          `json:"r"`
	S                string          `json:"s"`
	To               string          `json:"to"`
	TransactionIndex *hexutil.Uint64 `json:"transactionIndex"`
	Type             string          `json:"type"`
	V                string          `json:"v"`
	Value            *hexutil.Big    `json:"value"`
}

type LatestBlock struct {
	BlockNumber *hexutil.Uint64 `json:"blockNumber"`
}
