package blockchain

import (
	"crypt-coin-payment/models"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	_ "github.com/jinzhu/gorm"
	"log"
	"os"
	"strings"
)

var BtcRpcClient *rpcclient.Client

func init() {
	connCfg := &rpcclient.ConnConfig{
		Host:		os.Getenv("btc_node_address") + ":" + os.Getenv("btc_rpc_port") ,
		User:		os.Getenv("btc_rpc_user"),
		Pass:       os.Getenv("btc_rpc_pass"),
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	var err error
	BtcRpcClient, err = rpcclient.New(connCfg, nil)
	if err != nil {
		log.Println(err)
	}
}

func GetBtcRpcClient(paymentMethodId uint) *rpcclient.Client {
	return BtcRpcClient
}


type BtcFacade struct {}

func (btcFacade *BtcFacade) GetBalance(address string) (float64, error)  {
	return 0, nil
}

func (btcFacade *BtcFacade) GetBlockHash(blockHeight int64) (string, error) {
	blockHash, err := BtcRpcClient.GetBlockHash(blockHeight)
	if blockHash == nil {
		return "", err
	}
	return blockHash.String(), err
}

func (btcFacade *BtcFacade) GetBestBlock() (string, int64, error)  {
	blockNumber, err := BtcRpcClient.GetBlockCount()
	blockHash, err := BtcRpcClient.GetBlockHash(blockNumber)
	if err != nil {
		return "", 0, err
	}
	return blockHash.String(), blockNumber, err
}

func (btcFacade *BtcFacade) NextBlock(blockHeight int64, blockHash string) (string, bool, error)  {
	nextBlockHash, err := BtcRpcClient.GetBlockHash(blockHeight + 1)
	nextBlockHeader, err := BtcRpcClient.GetBlockHeader(nextBlockHash)
	if err != nil {
		return "", true, err
	}
	if strings.Compare(blockHash, nextBlockHeader.PrevBlock.String()) != 0 {
		return nextBlockHeader.BlockHash().String(), false, err
	} else {
		return nextBlockHeader.BlockHash().String(), true, err
	}
}

func (btcFacade *BtcFacade) ApplyNextBlock(blockHash string, blockHeight int64) error {
	hash, err := chainhash.NewHashFromStr(blockHash)
	block, err := BtcRpcClient.GetBlock(hash)
	if err != nil {
		return err
	}
	db := models.GetDB()
	blockModel := &models.Block{
		BlockHash:       blockHash,
		BlockNumber:     blockHeight,
		PaymentMethodId: 1,
	}
	err = db.Create(blockModel).Error
	if err != nil {
		log.Println(err)
		return err
	}
	for _, tx := range block.Transactions {
		txInDb := models.GetTransaction(tx.TxHash().String())
		if txInDb == nil {
			for _, vout := range tx.TxOut {
				_, addresses, _, err := txscript.ExtractPkScriptAddrs(vout.PkScript, &chaincfg.TestNet3Params)
				if err != nil {
					log.Println(err)
					continue
				}
				if len(addresses) == 1 {
					addressInDb := models.GetAddress(addresses[0].String())
					if addressInDb != nil {
						newTx := models.Transaction{
							OrderId:         addressInDb.OrderId,
							TransactionHash: tx.TxHash().String(),
							To:              addressInDb.Address,
							Value:           float64(vout.Value),
							BlockHash:       blockHash,
							BlockNumber:     uint(blockHeight),
							Type:            models.TYPE_PAYMENT,
							PaymentMethodId: 1,
						}
						db.Create(newTx)
					}
				}
			}
		} else {
			txInDb.BlockHash = blockHash
			txInDb.BlockNumber = uint(blockHeight)
			db.Save(txInDb)
		}
	}
	return nil
}

func (btcFacade *BtcFacade) RevertBlock(blockNumber int64) error  {
	currentBlockNumber := blockNumber
	blockHash, err := BtcRpcClient.GetBlockHash(blockNumber)
	if err != nil {
		return err
	}
	blockInDb := models.GetBlockByBlockNumber(uint(blockNumber), 1)
	for strings.Compare(blockInDb.BlockHash, blockHash.String()) != 0 {
		err := models.RevertTransactionInBlock(uint(currentBlockNumber), 1)
		if err != nil {
			return err
		}
		models.GetDB().Delete(blockInDb)
		currentBlockNumber -= 1
		blockInDb = models.GetBlockByBlockNumber(uint(currentBlockNumber), 1)
		blockHash, err = BtcRpcClient.GetBlockHash(currentBlockNumber)
		if err != nil {
			return err
		}
	}
	return nil
}