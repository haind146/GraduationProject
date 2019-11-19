package service

import (
	"errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

type KeyInterface interface {
	GenerateAccountFromMasterPubKey(string, uint32) (string, error)
	GenerateAddressFromAccount(string, uint32) (string, error)
}

type BtcKey struct {}
type EthKey struct {}

func KeyFactory(paymentMethodId int) KeyInterface {
	switch paymentMethodId {
	case 1:
		return &BtcKey{}
	//case 2:
	//	return &EthKey{}
	}
	return nil
}

func (btcKey *BtcKey) GenerateAddressFromAccount(account string, index uint32) (string, error) {
	accountKey, err := hdkeychain.NewKeyFromString(account)
	if err != nil {
		return "", err
	}
	childKey, _ := accountKey.Child(index)
	address, _ := childKey.Address(&chaincfg.TestNet3Params)
	return address.String(), nil
}

func (btcKey *BtcKey) GenerateAccountFromMasterPubKey(masterPubKeyStr string, index uint32) (string, error)  {
	masterPubKey, err := hdkeychain.NewKeyFromString(masterPubKeyStr)
	if err != nil {
		return "", err
	}
	if masterPubKey.Depth() > 0 {
		return "", errors.New("Not a master public key")
	}
	account, err := masterPubKey.Child(index)
	return account.String(), nil
}

//func (ethKey EthKey)  GenerateAddressFromAccount(account string, index uint32) string {
//	return "def"
//}
