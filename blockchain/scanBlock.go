package blockchain

import (
	"crypt-coin-payment/models"
	"log"
)

func ScanBlock(paymentMethodId uint) {
	facadeInstance := GetFacadeInterface(paymentMethodId)
	var currentBlockNumber int64
	var currentBlockHash string
	lastBlock := models.GetLatestBlock(paymentMethodId)
	if lastBlock == nil {
		blockHash, blockNumber, err := facadeInstance.GetBestBlock()
		if err != nil {
			log.Println("getBestBlock", err)
			return
		}
		err = facadeInstance.ApplyNextBlock(blockHash, blockNumber)
		if err != nil {
			log.Println("ApplyNextBlock", err)
			return
		}
		currentBlockNumber = blockNumber
		currentBlockHash = blockHash

	} else {
		currentBlockNumber = lastBlock.BlockNumber
	}
	for {
		nextBlockHash, isValid, err := facadeInstance.NextBlock(currentBlockNumber, currentBlockHash)
		if err != nil || nextBlockHash == "" {
			return
		}
		if nextBlockHash != "" && isValid {
			err := facadeInstance.ApplyNextBlock(nextBlockHash, currentBlockNumber+1)
			if err != nil {
				log.Println("ApplyNextBlock", err)
				return
			}
			currentBlockNumber++
			currentBlockHash = nextBlockHash
		}
		if !isValid {
			err := facadeInstance.RevertBlock(currentBlockNumber)
			if err != nil {
				log.Println("Revert Block", err)
			}
			return
		}
	}


}
