package subscriber

import (
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
	rpcClient := GetRpcClient(1)
	tx, err := rpcClient.DecodeRawTransaction(rawTx)
	if err != nil {
		log.Println(err)
	}

	for _, vout := range tx.Vout  {
		for _, address := range vout.ScriptPubKey.Addresses {
			order := models.FindOrderByAddress(address)
			if order != nil {

			}
		}
	}

	//log.Println(tx)
}