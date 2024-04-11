package main

import (
	"encoding/hex"
	"math/big"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
)

const TRANSFER_EVENT_TOPIC = "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" // NO 0x Prefix for ANY topic!!!

type Transfer struct {
	BlockNumber uint64 `json:"block_number" csv:"block_number"`

	TransactionHash string `json:"transaction_hash" csv:"transaction_hash"`
	LogIndex        uint   `json:"log_index" csv:"log_index"`
	// TxHashIdx       string      `csv:"id"`
	TokenAddress string `json:"token_address" csv:"token_address"`
	FromAddress  string `json:"from_address" csv:"from_address"`
	ToAddress    string `json:"to_address" csv:"to_address"`
	Value        string `json:"value" csv:"value"`
}

// ExtractTransferFromLog is a helper to extract EVM smart contract token transfer from the event logs
func ExtractTransferFromLog(logTopics []string, logData string, logContractAddress string, logIndex uint, logTxHash string, logBlockNum uint64) *Transfer {
	// topics := log.Topics
	if logTopics == nil || len(logTopics) < 1 {
		return nil
	}

	if logTopics[0] != TRANSFER_EVENT_TOPIC {
		return nil
	}

	topics_with_data := append(logTopics, chunkDataToHashes(logData)...)
	// txHash := log.TxHash
	// logIndex := log.Index
	if len(topics_with_data) != 4 {
		return nil
	}

	valBytes, err := hex.DecodeString(topics_with_data[3])
	chk(err)
	value := new(big.Int).SetBytes(valBytes)

	return &Transfer{
		BlockNumber:     logBlockNum,
		TokenAddress:    tron.EnsureTAddr(logContractAddress),
		FromAddress:     hash2Addr(topics_with_data[1]),
		ToAddress:       hash2Addr(topics_with_data[2]),
		Value:           value.String(),
		LogIndex:        logIndex,
		TransactionHash: logTxHash,
		// TxHashIdx:       log.TxHash.String() + "_" + strconv.Itoa(int(log.Index)),
	}
}

func chunkDataToHashes(b string) []string {
	rtn := make([]string, 0, (len(b)+31*2)/64) // len hash str == 64
	for i := 0; i+64 <= len(b); i += 64 {
		rtn = append(rtn, b[i:i+64])
	}

	return rtn
}

func hash2Addr(hash string) string {
	if len(hash) != 64 {
		panic("not a hash")
	}
	// addr := common.Address{}
	// copy(addr[:], )
	return tron.EnsureTAddr(hash[12*2:])
}
