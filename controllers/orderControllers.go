package controllers

import (
	"crypt-coin-payment/models"
	"encoding/json"
	"net/http"
	"strconv"

	//"u"
	u "crypt-coin-payment/utils"
)


var CreateOrder = func(w http.ResponseWriter, r *http.Request) {
	order := &models.Order{}
	userId := r.Context().Value("user") . (uint)
	err := json.NewDecoder(r.Body).Decode(order)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	application := models.GetApplication(order.ApplicationId)
	if application == nil || application.UserId != userId {
		u.Respond(w, u.Message(false, "Application Id not found"))
		return
	}

	resp := order.Create()
	u.Respond(w, resp)
}

var GetOrdersList = func(w http.ResponseWriter, r *http.Request) {
	appId, ok := r.URL.Query()["app"]
	if !ok || len(appId[0]) < 1 {
		resp := u.Message(false, "Url Param 'app' is missing")
		u.Respond(w, resp)
		return
	}
	appIdUint, err := strconv.ParseUint(appId[0], 10, 64)
	if err != nil {
		resp := u.Message(false, "Error parse appId")
		u.Respond(w, resp)
		return
	}

	data := models.OrdersList(uint(appIdUint))
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetTransactionsByOrder = func(w http.ResponseWriter, r *http.Request) {
	appId, ok := r.URL.Query()["order_id"]
	if !ok || len(appId[0]) < 1 {
		resp := u.Message(false, "Url Param 'app' is missing")
		u.Respond(w, resp)
		return
	}
	orderId, err := strconv.ParseUint(appId[0], 10, 64)
	if err != nil {
		resp := u.Message(false, "Error parse appId")
		u.Respond(w, resp)
		return
	}

	data := models.TransactionsByOrder(uint(orderId))
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}