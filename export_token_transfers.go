package main

import (
	"encoding/csv"
	"io"
	"log"
	"sync"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"github.com/jszwec/csvutil"
	"golang.org/x/exp/slices"
)

type ExportTransferOptions struct {
	// outputType string // failed to tfOutput to a conn with
	tfOutput         io.Writer
	logOutput        io.Writer
	internalTxOutput io.Writer
	receiptOutput    io.Writer

	Worker int

	ProviderURI string `json:"provider_uri,omitempty"`
	StartBlock  uint64 `json:"start_block,omitempty"`
	EndBlock    uint64 `json:"end_block,omitempty"`

	// extension
	StartTimestamp string `json:"start_timestamp,omitempty"`
	EndTimestamp   string `json:"end_timestamp,omitempty"`

	Contracts []string `json:"contracts,omitempty"`
}

// ExportTransfers is the main func for handling export_transfers command
func ExportTransfers(options *ExportTransferOptions) {
	cli := tron.NewTronClient(options.ProviderURI)

	var tfEncoder, logEncoder, internalTxEncoder, receiptEncoder *csvutil.Encoder

	if options.tfOutput != nil {
		tfWriter := csv.NewWriter(options.tfOutput)
		defer tfWriter.Flush()
		tfEncoder = csvutil.NewEncoder(tfWriter)
	}

	if options.logOutput != nil {
		logWriter := csv.NewWriter(options.logOutput)
		defer logWriter.Flush()
		logEncoder = csvutil.NewEncoder(logWriter)
	}

	if options.internalTxOutput != nil {
		internalTxWriter := csv.NewWriter(options.internalTxOutput)
		defer internalTxWriter.Flush()
		internalTxEncoder = csvutil.NewEncoder(internalTxWriter)
	}

	if options.receiptOutput != nil {
		receiptWriter := csv.NewWriter(options.receiptOutput)
		defer receiptWriter.Flush()
		receiptEncoder = csvutil.NewEncoder(receiptWriter)
	}

	filterLogContracts := make([]string, len(options.Contracts))
	for i, addr := range options.Contracts {
		filterLogContracts[i] = tron.EnsureHexAddr(addr)[2:] // hex addr with 41 prefix
	}

	if options.StartTimestamp != "" {
		// fast locate estimate start height

		number, err := BlockNumberFromDateTime(cli, options.StartTimestamp, FirstAfterTimestamp)
		if err != nil {
			panic(err)
		}
		options.StartBlock = *number

	}

	if options.EndTimestamp != "" {
		number, err := BlockNumberFromDateTime(cli, options.EndTimestamp, LastBeforeTimestamp)
		if err != nil {
			panic(err)
		}
		options.EndBlock = *number
	}

	log.Printf("try parsing token transfers from block %d to %d", options.StartBlock, options.EndBlock)

	if options.StartBlock != 0 && options.EndBlock != 0 {
		for number := options.StartBlock; number <= options.EndBlock; number++ {
			txInfos := cli.GetTxInfosByNumber(number)
			for txIndex, txInfo := range txInfos {
				txHash := txInfo.ID

				if receiptEncoder != nil {
					receiptEncoder.Encode(NewCsvReceipt(number, txHash, uint(txIndex), txInfo.ContractAddress, txInfo.Fee, txInfo.Receipt))
				}

				for logIndex, log := range txInfo.Log {
					if len(filterLogContracts) != 0 && !slices.Contains(filterLogContracts, log.Address) {
						continue
					}

					if tfEncoder != nil {
						tf := ExtractTransferFromLog(log.Topics, log.Data, log.Address, uint(logIndex), txHash, number)
						if tf != nil {
							err := tfEncoder.Encode(tf)
							chk(err)
						}
					}

					if logEncoder != nil {
						err := logEncoder.Encode(NewCsvLog(number, txHash, uint(logIndex), log))
						chk(err)
					}

				}

				if internalTxEncoder != nil {
					for internalIndex, internalTx := range txInfo.InternalTransactions {
						for callInfoIndex, callInfo := range internalTx.CallValueInfo {
							err := internalTxEncoder.Encode(NewCsvInternalTx(number, txHash, uint(internalIndex), internalTx, uint(callInfoIndex), callInfo.TokenID, callInfo.CallValue))
							chk(err)
						}
					}
				}

			}

			log.Printf("parsed block %d", number)
		}

		return
	}
}

// ExportTransfers is the main func for handling export_transfers command
func ExportTransfersWithWorkers(options *ExportTransferOptions, workers uint) {
	cli := tron.NewTronClient(options.ProviderURI)

	var tfEncCh, logEncCh, internalTxEncCh, receiptEncCh chan any
	var receiverWG sync.WaitGroup

	if options.tfOutput != nil {
		tfWriter := csv.NewWriter(options.tfOutput)
		defer tfWriter.Flush()
		tfEncoder := csvutil.NewEncoder(tfWriter)
		tfEncCh = createCSVEncodeCh(&receiverWG, tfEncoder, workers)
	}

	if options.logOutput != nil {
		logWriter := csv.NewWriter(options.logOutput)
		defer logWriter.Flush()
		logEncoder := csvutil.NewEncoder(logWriter)
		logEncCh = createCSVEncodeCh(&receiverWG, logEncoder, workers)
	}

	if options.internalTxOutput != nil {
		internalTxWriter := csv.NewWriter(options.internalTxOutput)
		defer internalTxWriter.Flush()
		internalTxEncoder := csvutil.NewEncoder(internalTxWriter)
		internalTxEncCh = createCSVEncodeCh(&receiverWG, internalTxEncoder, workers)
	}

	if options.receiptOutput != nil {
		receiptWriter := csv.NewWriter(options.receiptOutput)
		defer receiptWriter.Flush()
		receiptEncoder := csvutil.NewEncoder(receiptWriter)
		receiptEncCh = createCSVEncodeCh(&receiverWG, receiptEncoder, workers)
	}

	filterLogContracts := make([]string, len(options.Contracts))
	for i, addr := range options.Contracts {
		filterLogContracts[i] = tron.EnsureHexAddr(addr)[2:] // hex addr with 41 prefix
	}

	log.Printf("try parsing token transfers from block %d to %d", options.StartBlock, options.EndBlock)

	exportWork := func(wg *sync.WaitGroup, workerID uint) {
		for number := options.StartBlock + uint64(workerID); number <= options.EndBlock; number += uint64(workers) {
			txInfos := cli.GetTxInfosByNumber(number)
			for txIndex, txInfo := range txInfos {
				txHash := txInfo.ID

				if options.receiptOutput != nil {
					receiptEncCh <- NewCsvReceipt(number, txHash, uint(txIndex), txInfo.ContractAddress, txInfo.Fee, txInfo.Receipt)
				}

				for logIndex, log := range txInfo.Log {
					if len(filterLogContracts) != 0 && !slices.Contains(filterLogContracts, log.Address) {
						continue
					}

					if options.tfOutput != nil {
						tf := ExtractTransferFromLog(log.Topics, log.Data, log.Address, uint(logIndex), txHash, number)
						if tf != nil {
							tfEncCh <- tf
						}
					}

					if options.logOutput != nil {
						logEncCh <- NewCsvLog(number, txHash, uint(logIndex), log)
					}

				}

				if options.internalTxOutput != nil {
					for internalIndex, internalTx := range txInfo.InternalTransactions {
						for callInfoIndex, callInfo := range internalTx.CallValueInfo {
							internalTxEncCh <- NewCsvInternalTx(number, txHash, uint(internalIndex), internalTx, uint(callInfoIndex), callInfo.TokenID, callInfo.CallValue)
						}
					}
				}

			}

			log.Printf("parsed block %d", number)
		}

		wg.Done()
	}

	var senderWG sync.WaitGroup
	for workerID := uint(0); workerID < workers; workerID++ {
		senderWG.Add(1)
		go exportWork(&senderWG, workerID)
	}

	senderWG.Wait()
	if options.tfOutput != nil {
		close(tfEncCh)
	}
	if options.logOutput != nil {
		close(logEncCh)
	}
	if options.internalTxOutput != nil {
		close(internalTxEncCh)
	}
	if options.receiptOutput != nil {
		close(receiptEncCh)
	}

	receiverWG.Wait()
}
