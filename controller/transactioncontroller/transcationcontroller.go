package transactioncontroller

import (
	"encoding/json"
	"net/http"

	"github.com/meinantoyuriawan/gobank-v2/helper"
	"github.com/meinantoyuriawan/gobank-v2/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Getsaldo(w http.ResponseWriter, r *http.Request) {
	userInput := models.Account{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	// get account data
	user := models.Account{}
	// process the pin (hash and match it with db)
	if err := models.DB.Where("Username = ?", userInput.Number).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// verify is pin valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(userInput.Pin)); err != nil {
		response := map[string]string{"message": "Pin salah"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

}

func Transfer(w http.ResponseWriter, r *http.Request) {

}

func Getrecent(w http.ResponseWriter, r *http.Request) {

}
