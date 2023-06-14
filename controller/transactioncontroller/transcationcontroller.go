package transactioncontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	if err := models.DB.Where("Number = ?", userInput.Number).First(&user).Error; err != nil {
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

	// get status
	status := models.Status{}
	if err := models.DB.Where("acc_number = ?", userInput.Number).First(&status).Error; err != nil {
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
	// check isBlocked
	if status.IsBlocked == 0 {
		// process the status db for pin related
		// verify is pin valid
		if err := bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(userInput.Pin)); err != nil {
			// status.PinAttempt -1, if it reaches 0 (from 3), status.IsBlocked = 1
			newStatus := pinValidation(status)
			models.DB.Save(newStatus)
			response := map[string]string{"message": "Pin salah"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}
		// if the user succesfully to log in, do the reset Validation
		newStatus := resetValidation(status)
		models.DB.Save(newStatus)
		bal := strconv.Itoa(int(user.Balance))
		response := map[string]string{"balance": bal}
		helper.ResponseJSON(w, http.StatusOK, response)
	} else {
		response := map[string]string{"message": "Account Blocked"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

}

func Transfer(w http.ResponseWriter, r *http.Request) {
	// get the user and recipient data
	transferInput := models.Transfer{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&transferInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()
	// pinValidation, resetValidation

	// get source account data
	sourceAccount := models.Account{}
	// process the pin (hash and match it with db)
	if err := models.DB.Where("Number = ?", transferInput.Source).First(&sourceAccount).Error; err != nil {
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
	// source status

	// get status
	sourceStatus := models.Status{}
	if err := models.DB.Where("acc_number = ?", transferInput.Source).First(&sourceStatus).Error; err != nil {
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
	// check isBlocked
	if sourceStatus.IsBlocked == 0 {
		// process the status db for pin related
		// verify is pin valid
		if err := bcrypt.CompareHashAndPassword([]byte(sourceAccount.Pin), []byte(transferInput.Pin)); err != nil {
			// status.PinAttempt -1, if it reaches 0 (from 3), status.IsBlocked = 1
			newStatus := pinValidation(sourceStatus)
			models.DB.Save(newStatus)
			response := map[string]string{"message": "Pin salah"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}
		// if the user succesfully to log in, do the reset Validation
		newStatus := resetValidation(sourceStatus)
		models.DB.Save(newStatus)
		// recipient Number isValid
		recipient := models.Account{}
		if err := models.DB.Where("Number = ?", transferInput.Recipient).First(&recipient).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				response := map[string]string{"message": "Destination Invalid"}
				helper.ResponseJSON(w, http.StatusBadRequest, response)
				return
			default:
				response := map[string]string{"message": err.Error()}
				helper.ResponseJSON(w, http.StatusInternalServerError, response)
				return
			}
		}
		// account Balance isSufficient
		if sourceAccount.Balance >= transferInput.Amount {
			// transfer
			sourceAccount.Balance = sourceAccount.Balance - transferInput.Amount
			recipient.Balance = recipient.Balance + transferInput.Amount

			// store to activity db
			storeActivity(sourceAccount.UserId, sourceStatus.AccNumber, "Transfer", recipient.Number, w)

			// save to account db
			models.DB.Save(sourceAccount)
			models.DB.Save(recipient)
			response := map[string]string{"Message": "Transfer Success"}
			helper.ResponseJSON(w, http.StatusOK, response)
		} else {
			response := map[string]string{"message": "Insufficient Balance"}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}
	} else {
		response := map[string]string{"message": "Account Blocked"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

}

func Getrecent(w http.ResponseWriter, r *http.Request) {
	// create new Database about account activity
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
	if err := models.DB.Where("Number = ?", userInput.Number).First(&user).Error; err != nil {
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

	// get status
	status := models.Status{}
	if err := models.DB.Where("acc_number = ?", userInput.Number).First(&status).Error; err != nil {
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
	// access the database using user input

	var activities []models.Activity

	// check isBlocked
	if status.IsBlocked == 0 {
		// process the status db for pin related
		// verify is pin valid
		if err := bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(userInput.Pin)); err != nil {
			// status.PinAttempt -1, if it reaches 0 (from 3), status.IsBlocked = 1
			newStatus := pinValidation(status)
			models.DB.Save(newStatus)
			response := map[string]string{"message": "Pin salah"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}
		// if the user succesfully to log in, do the reset Validation
		newStatus := resetValidation(status)

		models.DB.Save(newStatus)

		// get the data activity
		if err := models.DB.Where("acc_number = ?", userInput.Number).Find(&activities).Error; err != nil {
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

		response := activities
		helper.ResponseJSON(w, http.StatusOK, response)
	} else {
		response := map[string]string{"message": "Account Blocked"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}
}

func pinValidation(status models.Status) models.Status {
	// get the data from number
	if status.IsBlocked == 0 {
		if status.PinAttempt > 0 {
			status.PinAttempt = status.PinAttempt - 1
			// push to db
			// fmt.Println(status)
			return status
		}
		// fmt.Println("account blocked")
		if status.PinAttempt == 0 {
			status.IsBlocked = 1
			// push to db
			fmt.Println("account blocked")
			// fmt.Println(status)
			return status
		}
		return status
	} else {
		// fmt.Println("account blocked")
		return status
	}

	// check the isBlocked
	// if it's not blocked reduce the attempt
	// if the attempt is reaching zero, change the isBlocked
}

func resetValidation(status models.Status) models.Status {
	// set the attempt to max (3)
	status.PinAttempt = 3
	return status
}

func storeActivity(sourceUserId int64, sourceAccNumber int64, description string, recipientAccNumber int64, w http.ResponseWriter) {
	activity := models.Activity{
		UserId:      sourceUserId,
		AccNumber:   sourceAccNumber,
		Description: description,
		Recipient:   recipientAccNumber,
		DateTime:    time.Now(),
	}
	if err := models.DB.Create(&activity).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}
}
