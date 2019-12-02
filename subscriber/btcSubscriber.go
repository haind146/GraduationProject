package subscriber

import (
	"crypt-coin-payment/blockchain"
	"crypt-coin-payment/models"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	zmq "github.com/pebbe/zmq4"
	"log"
)

type BtcSubscriber struct {
	ZmqPubEndpoint string
}

func (subscriber *BtcSubscriber) Subscribe() error {
	socket, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		return err
	}
	defer socket.Close()
	err = socket.Connect(subscriber.ZmqPubEndpoint)
	if err != nil {
		log.Println("SocketConnect", err)
		return err
	}
	err = socket.SetSubscribe("hashtx")
	if err != nil {
		log.Println("SetSubscribe", err)
		return err
	}

	for {
		msg, e := socket.RecvMessageBytes(0)
		hash := hex.EncodeToString(msg[1])
		if e != nil {
			log.Println("RecvMessageBytes", e)
			break
		}
		hashtx, _ := chainhash.NewHashFromStr(hash)
		go HandleNewTransaction(hashtx)
	}
	return nil
}

func HandleNewTransaction(hash *chainhash.Hash)  {
	rpcClient := blockchain.GetBtcRpcClient(1)
	tx, err := rpcClient.GetRawTransaction(hash)
	if err != nil {
		log.Println("GetRawTransaction", err)
		return
	}
	existTxInDb := models.GetTransaction(tx.Hash().String())
	if existTxInDb != nil {
		return
	}

	for _, vout := range tx.MsgTx().TxOut  {
		_, addresses, _, err := txscript.ExtractPkScriptAddrs(vout.PkScript, &chaincfg.TestNet3Params)
		for _, address := range addresses {
			addressModel := models.GetAddress(address.String())

			if addressModel != nil {
				fmt.Println(addressModel.Address)
				dbTx := models.GetDB().Begin()
				transaction := &models.Transaction{
					OrderId:         addressModel.OrderId,
					TransactionHash: tx.Hash().String(),
					To:              address.String(),
					Value:           float64(vout.Value)/100000000,
					Type:            models.TYPE_PAYMENT,
					PaymentMethodId: 1,
				}
				err = dbTx.Create(transaction).Error
				log.Println(err)
				order := models.FindOrderByAddress(address.String())
				order.ReceivedValue += transaction.Value
				err = dbTx.Save(order).Error
				if err != nil {
					log.Println(err)
					dbTx.Rollback()
				} else {
					dbTx.Commit()
				}
				break
			}
		}
	}
	//log.Println(tx)
}