package subscriber

import (
	"github.com/btcsuite/btcd/rpcclient"
	"log"
	"os"
)

type SubInterface interface {
	Subscribe() error
}

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

func SubscriberFactory(paymentMethodId uint) SubInterface {
	switch paymentMethodId {
	case 1:
		return &BtcSubscriber{
			ZmqPubEndpoint: os.Getenv("btc_zmq_public_endpoint"),
		}
	default:
		return nil
	}
}

func GetRpcClient(patmentMethodId uint) *rpcclient.Client  {
	return BtcRpcClient
}
