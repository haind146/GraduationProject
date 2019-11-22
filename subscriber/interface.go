package subscriber

import (
	"os"
)

type SubInterface interface {
	Subscribe() error
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
