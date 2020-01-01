package models

import (
	"crypt-coin-payment/keychain"
	u "crypt-coin-payment/utils"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

const (
	ORDER_CREATED = 1
	ORDER_MEMPOOL = 2
	ORDER_INBLOCK = 3
	ORDER_SWEEPED = 4
)

type Order struct {
	gorm.Model
	PartnerOrderId string `json:"partner_order_id"`
	Amount float64 `json:"amount"`
	ReceivedValue float64 `json:"receive"`
	ReceiveAddress string `gorm:"type:varchar(100);unique_index"`
	PaymentMethodId uint `json:"payment_method_id"`
	ApplicationId uint `json:"application_id"`
	Status uint `json:"status"`
	ExpiredTime uint `json:"expired_time"`
}

/*
 This struct function validate the required parameters sent through the http request body

returns message and true if the requirement is met
*/
func (order *Order) Validate() (map[string] interface{}, bool) {

	if order.Amount <= 0 {
		return u.Message(false, "Amount should greater than 0"), false
	}
	if order.ApplicationId <= 0 {
		return u.Message(false, "Application Id not found"), false
	}
	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (order *Order) Create() map[string] interface{} {

	if resp, ok := order.Validate(); !ok {
		return resp
	}

	appPubKey := GetApplicationPubicKey(order.ApplicationId)
	if appPubKey == nil {
		return u.Message(false, "App public key not found")
	}

	keyService := keychain.KeyFactory(1)
	receiveAddress, err := keyService.GenerateAddressFromAccount(appPubKey.PublicKey, uint32(appPubKey.NumOfAddressGenerated))
	if err != nil {
		return u.Message(false, "Error when generate address")
	}
	//err = blockchain.ImportAddress(receiveAddress)
	//if err != nil {
	//	return u.Message(false, "Error when import address to node")
	//}

	order.ReceiveAddress = receiveAddress
	order.Status = ORDER_CREATED
	appPubKey.NumOfAddressGenerated = appPubKey.NumOfAddressGenerated + 1

	GetDB().Create(order)

	address := &Address{
		Address: receiveAddress,
		OrderId: order.ID,
		Balance: 0,
		PendingReceive: 0,
		PendingSent: 0,
		MnemonicPath: strconv.Itoa(int(appPubKey.Index)) + "/" + strconv.Itoa(int(appPubKey.NumOfAddressGenerated -1)),
	}
	err = GetDB().Create(address).Error
	fmt.Println("Create Address", err)
	GetDB().Save(appPubKey)
	resp := u.Message(true, "success")
	resp["order"] = order
	return resp
}

func FindOrderByAddress (address string) *Order {
	order := &Order{}
	err := GetDB().Table("orders").Where("receive_address = ?", address).First(order).Error
	if err != nil {
		return nil
	}
	return order
}

func FindOrerById(id uint) *Order {
	order := &Order{}
	err := GetDB().First(order, id).Error
	if err != nil {
		return nil
	}
	return order
}

func OrdersList(appId uint) []*Order {
	orders := make([]*Order, 0)
	err := GetDB().Table("orders").Where("application_id = ?", appId).Find(&orders).Error
	if err != nil {
		log.Println("OrdersList", err)
		return nil
	}
	return orders
}