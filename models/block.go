package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Block struct {
	gorm.Model
	BlockHash string `gorm:"type:varchar(100);unique_index"`
	BlockNumber int64 `json:"block_number"`
	PaymentMethodId uint `json:"payment_method_id"`
}

func GetLatestBlock(paymentMethodId uint) *Block  {
	block := &Block{}
	err := GetDB().Table("blocks").Where("payment_method_id = ?", paymentMethodId).Order("block_number DESC").First(block).Error
	if err != nil {
		fmt.Println(err.Error)
		return nil
	}
	return block
}

func GetBlockByBlockNumber(blockNumber uint, paymentMethodId uint) *Block {
	block := &Block{}
	err := GetDB().Table("blocks").Where("block_number = ? AND payment_method_id = ?", blockNumber, paymentMethodId).First(block)
	if err != nil {
		fmt.Println(err)
		return nil
	}


	return block
}

//func (block *Block) Create() (map[string] interface{}) {
//	//GetDB().Create(address)
//}

//func GetBlock(id uint) (*Address) {
//	//
//	//address := &Address{}
//	//err := GetDB().Table("addresses").Where("id = ?", id).First(application).Error
//	//if err != nil {
//	//	return nil
//	//}
//	//return address
//}
