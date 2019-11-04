package controllers

import (
	"crypt-coin-payment/models"
	"encoding/json"
	"net/http"
	//"u"
	u "crypt-coin-payment/utils"
)



var CreateMasterPublicKey = func(w http.ResponseWriter, r *http.Request) {
	masterPublicKey := &models.MasterPublicKey{}
	userId := r.Context().Value("user") . (uint)
	err := json.NewDecoder(r.Body).Decode(masterPublicKey)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	masterPublicKey.UserId = userId
	resp := masterPublicKey.Create()
	u.Respond(w, resp)
}
