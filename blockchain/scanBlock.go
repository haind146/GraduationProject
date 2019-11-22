package blockchain

import (
	"crypt-coin-payment/models"
	"log"
)

func ScanBlock(paymentMethodId uint)  {
	facadeInstance := GetFacadeInterface(paymentMethodId)
	var currentBlockNumber int64
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
	} else {
		currentBlockNumber = lastBlock.BlockNumber
	}
	if currentBlockNumber > 0 {
		nextBlock, err := facadeInstance.GetBlockHash(currentBlockNumber + 1)
		if nextBlock == "" || err != nil {
			return
		}
		nextBlockHash, isValid, err := facadeInstance.NextBlock(currentBlockNumber+1, nextBlock)
		if nextBlockHash != "" && isValid {
			err := facadeInstance.ApplyNextBlock(nextBlockHash, currentBlockNumber+1)
			if err != nil {
				log.Println("ApplyNextBlock", err)
				return
			}
		}
		if !isValid {
			err := facadeInstance.RevertBlock(currentBlockNumber)
			if err != nil {
				log.Println("Revert Block", err)
				return
			}
		}
	}
}