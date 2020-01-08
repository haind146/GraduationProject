package subscriber

import (
	"crypt-coin-payment/blockchain"
	"crypt-coin-payment/models"
	"encoding/hex"
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
		log.Println("Subcribe", err)
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

	for outputIndex, vout := range tx.MsgTx().TxOut  {
		_, addresses, _, err := txscript.ExtractPkScriptAddrs(vout.PkScript, &chaincfg.TestNet3Params)
		for _, address := range addresses {
			addressModel := models.GetAddress(address.String())

			if addressModel != nil {
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
				utxo := &models.Utxo{
					TxId:        transaction.ID,
					OutputIndex: uint(outputIndex),
					Value:       float64(vout.Value)/100000000,
					Spent:       false,
				}
				err = dbTx.Create(utxo).Error

				log.Println(err)
				order := models.FindOrderByAddress(address.String())
				order.ReceivedValue += transaction.Value
				order.Status = models.ORDER_MEMPOOL
				err = dbTx.Save(order).Error
				if err != nil {
					log.Println("Save Order", err)
					dbTx.Rollback()
				} else {
					dbTx.Commit()
				}
				return
			}
		}
	}

	txInDb := models.GetTransaction(tx.MsgTx().TxIn[0].PreviousOutPoint.Hash.String())
	if txInDb != nil && txInDb.Type == models.TYPE_PAYMENT {

		transaction := &models.Transaction{
			TransactionHash: tx.Hash().String(),
			Type:			 models.TYPE_SWEEP,
			PaymentMethodId: 1,
		}
		transaction.Value = float64(tx.MsgTx().TxOut[0].Value)/100000000
		_, addresses, _, err := txscript.ExtractPkScriptAddrs(tx.MsgTx().TxOut[0].PkScript, &chaincfg.TestNet3Params)
		if err != nil {
			log.Println("ExtractPkScriptAddrs", err)
		}
		transaction.To = addresses[0].String()

		totalValue := 0.0
		for _, vin := range tx.MsgTx().TxIn {
			prevTxHash := vin.PreviousOutPoint.Hash.String()
			models.SpendUtxo(prevTxHash)
			txDb := models.GetTransaction(prevTxHash)
			totalValue += txDb.Value
			models.UpdateOrderStatus(txDb.OrderId, models.ORDER_SWEEPED)
		}

		transaction.Fee = totalValue - transaction.Value
		order := models.FindOrerById(txInDb.OrderId)
		transaction.ApplicationId = order.ApplicationId
		err = models.GetDB().Create(transaction).Error
		if err != nil {
			log.Println("SaveSweepTransaction", err)
		}
	}
}