package models

import (
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	OrderId string `json:"order_id"`
	TransactionHash uint `json:"transaction_hash"`
	From uint `json:"from"`
	To uint `json:"to"`
	Value uint `json:"value"`
	Fee uint `json:"fee"`
	BlockHash uint `json:"block_hash"`
	BlockNumber uint `json:"block_number"`
	Type uint `json:"type"`
	PaymentMethodIs uint `json:"payment_method_id"`
}

//func (transaction *Transaction) Create() (map[string] interface{}) {
//	//GetDB().Create(address)
//}
//
//func GetTransaction(id uint) (*Transaction) {
//	//
//	//address := &Address{}
//	//err := GetDB().Table("addresses").Where("id = ?", id).First(application).Error
//	//if err != nil {
//	//	return nil
//	//}
//	//return address
//}
