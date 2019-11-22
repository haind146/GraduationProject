package subscriber

import (
	"crypt-coin-payment/blockchain"
	"crypt-coin-payment/models"
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
	log.Println("websocket connected")
	defer socket.Close()
	err = socket.Connect(subscriber.ZmqPubEndpoint)
	if err != nil {
		log.Println(err)
		return err
	}
	err = socket.SetSubscribe("rawtx")
	if err != nil {
		log.Println(err)
		return err
	}

	for {
		msg, e := socket.RecvMessageBytes(0)
		if e != nil {
			log.Println(e)
			break
		}
		go HandleNewTransaction(msg[1])
	}
	return nil
}

func HandleNewTransaction(rawTx []byte)  {
	rpcClient := blockchain.GetBtcRpcClient(1)
	tx, err := rpcClient.DecodeRawTransaction(rawTx)
	if err != nil {
		log.Println(err)
		return
	}
	existTxInDb := models.GetTransaction(tx.Hash)
	if existTxInDb != nil {
		return
	}

	for _, vout := range tx.Vout  {
		for _, address := range vout.ScriptPubKey.Addresses {
			addressModel := models.GetAddress(address)
			if addressModel != nil {
				dbTx := models.GetDB().Begin()
				transaction := &models.Transaction{
					OrderId:         addressModel.OrderId,
					TransactionHash: tx.Hash,
					To:              address,
					Value:           vout.Value,
					BlockHash:       tx.BlockHash,
					Type:            models.TYPE_PAYMENT,
					PaymentMethodId: 1,
				}
				err = dbTx.Create(transaction).Error
				order := models.FindOrderByAddress(address)
				order.ReceivedValue += vout.Value
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