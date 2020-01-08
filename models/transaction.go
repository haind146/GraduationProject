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
	ApplicationId uint `json:"application_id"`
	TransactionHash string `gorm:"type:varchar(100);unique_index"`
	From string `gorm:"type:varchar(100);index"`
	To string `gorm:"type:varchar(100);index"`
	Value float64 `json:"value"`
	Fee float64 `json:"fee"`
	BlockHash string `json:"block_hash"`
	BlockNumber uint `json:"block_number"`
	Type uint `json:"type"`
	PaymentMethodId uint `json:"payment_method_id"`
}

type Utxo struct {
	gorm.Model
	TxId uint `json:"tx_id"`
	OutputIndex uint `json:"output_index"`
	Value float64 `json:"value"`
	Spent bool `json:"spent"`
}

func (transaction *Transaction) Create() error {
	result := GetDB().Create(transaction)
	return result.Error
}
func GetTransaction(txHash string) *Transaction {
	tx := &Transaction{}
	err := GetDB().Table("transactions").Where("transaction_hash = ?", txHash).First(tx).Error
	if err != nil {
		//log.Println("GetTransaction", err)
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

func TransactionsByOrder(orderId uint) []*Transaction {
	transactions := make([]*Transaction, 0)
	err := GetDB().Table("transactions").Where("order_id = ?", orderId).Find(&transactions).Error
	if err != nil {
		log.Println("TransactionsByAddress", err)
		return nil
	}
	return transactions
}

func SpendUtxo(txHash string)  {
	transaction := &Transaction{}
	db.Where("transaction_hash = ?", txHash).First(transaction)
	if transaction != nil {
		utxo := &Utxo{}
		db.Where("tx_id = ?", transaction.ID).First(utxo)
		utxo.Spent = true
		err := db.Save(utxo).Error
		if err != nil {
			log.Println("UpdateUtxo", err)
		}
	}
}