package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/hdkeychain"
	"net/http"
	"regexp"
)

func Message(status bool, message string) (map[string]interface{}) {
	return map[string]interface{} {"status" : status, "message" : message}
}

func Respond(w http.ResponseWriter, data map[string] interface{})  {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ValidateEmail(email string) bool {
	Re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return Re.MatchString(email)
}

func Generate64BytesRandom() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	return hex.EncodeToString(key)
}

func ValidateMasterPublicKey(pubKeyStr string) bool  {
	masterPublicKey, err := hdkeychain.NewKeyFromString(pubKeyStr)
	fmt.Println(err)
	return err == nil && !masterPublicKey.IsPrivate() && masterPublicKey.Depth() == 0
}