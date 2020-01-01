package controllers

import (
	"crypt-coin-payment/blockchain"
	"crypt-coin-payment/models"
	"encoding/json"
	"net/http"
	"strconv"

	//"u"
	u "crypt-coin-payment/utils"
)



var CreateApplication = func(w http.ResponseWriter, r *http.Request) {
	application := &models.Application{}
	userId := r.Context().Value("user") . (uint)
	err := json.NewDecoder(r.Body).Decode(application)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	application.UserId = userId
	resp := application.Create()
	u.Respond(w, resp)
}

var GetApplicationsList = func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user") . (uint)
	data := models.ApplicationsList(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetSweepMoneyInfo = func(w http.ResponseWriter, r *http.Request) {
	appId, ok := r.URL.Query()["application_id"]
	if !ok || len(appId[0]) < 1 {
		resp := u.Message(false, "Url Param 'app' is missing")
		u.Respond(w, resp)
		return
	}
	applicationId, err := strconv.ParseUint(appId[0], 10, 64)
	if err != nil {
		resp := u.Message(false, "Error parse appId")
		u.Respond(w, resp)
		return
	}
	data := blockchain.SweepInfo(uint(applicationId))
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}