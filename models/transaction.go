package models

import (
	"github.com/jinzhu/gorm"
	"log"
)

const TYPE_PAYMENT  = 1
const TYPE_SWEEP  = 2

type Transaction struct {
	gorm.Model
	OrderId uint `json:"order_id"`
	TransactionHash string `gorm:"type:varchar(100);unique_index"`
	From string `gorm:"type:varchar(100);index"`
	To string `gorm:"type:varchar(100);index"`
	Value float64 `json:"value"`
	Fee uint `json:"fee"`
	BlockHash string `json:"block_hash"`
	BlockNumber uint `json:"block_number"`
	Type uint `json:"type"`
	PaymentMethodId uint `json:"payment_method_id"`
}


func (transaction *Transaction) Create() error {
	result := GetDB().Create(transaction)
	return result.Error
}
func GetTransaction(txHash string) *Transaction {
	tx := &Transaction{}
	err := GetDB().Table("transactions").Where("transaction_hash = ?", txHash).First(tx).Error
	if err != nil {
		return nil
	}
	return tx
}

func GetFirstTransaction(paymentMethodId uint) *Transaction {
	tx := &Transaction{}
	err := GetDB().Table("transactions").Where("payment_method_id = ?", paymentMethodId).Order("id ASC").First(tx).Error
	if err != nil {
		log.Println(err)
		return nil
	}
	return tx
}

func RevertTransactionInBlock(blockNumber uint, paymentMethodId uint) error  {
	dbTx := db.Begin()
	err := dbTx.Table("orders").Where("id IN (SELECT order_id FROM transactions WHERE block_number = ? AND payment_method_id = ?)", blockNumber, paymentMethodId).Updates(map[string]interface{}{"status": 1}).Error
	err = dbTx.Table("transactions").Where("block_number = ? AND payment_method_id = ?", blockNumber, paymentMethodId).Updates(map[string]interface{}{"block_number": 0}).Error
	if err != nil {
		dbTx.Rollback()
		return err
	}
	dbTx.Commit()
	return nil
}