package models

import (
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	OrderId string `json:"order_id"`
	TransactionHash string `gorm:"type:varchar(100);unique_index"`
	From string `gorm:"type:varchar(100);index"`
	To string `gorm:"type:varchar(100);index"`
	Value uint `json:"value"`
	Fee uint `json:"fee"`
	BlockHash uint `json:"block_hash"`
	BlockNumber uint `json:"block_number"`
	Type uint `json:"type"`
	PaymentMethodIs uint `json:"payment_method_id"`
}


func (transaction *Transaction) Create() error {
	result := GetDB().Create(transaction)
	return result.Error
}
//func GetTransaction(id uint) (*Transaction) {
//	//
//	//address := &Address{}
//	//err := GetDB().Table("addresses").Where("id = ?", id).First(application).Error
//	//if err != nil {
//	//	return nil
//	//}
//	//return address
//}
