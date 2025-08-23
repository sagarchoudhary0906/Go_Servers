package mainController

import (
	"database/sql"
	"errors"
	"expense_tracker/googledbcontroller"
	idgenerator "expense_tracker/idGenrateContoller"
	"log"
	"strconv"
)

const GOOGLE_LOGIN_TOKEN = 1
const PHONE_LOGIN_TOKEN = 2

var DB *sql.DB

func IdController(data map[string]string) (string, error) {
	log.Println("IdController: data: ", data, "data keys", data["lt"], "token", data["token"])
	login_token := data["lt"]
	login_token_int, _ := strconv.Atoi(login_token)
	token := data["token"]

	if login_token == "" || token == "" {
		log.Println("IdController: login_token or token is empty")
		return "", errors.New("login_token or token is empty")
	}

	var user_id string
	var err error

	switch login_token_int {
	case GOOGLE_LOGIN_TOKEN:
		log.Println("IdController: google login token")
		user_id, err = googleIdController(token)
		if err != nil {
			log.Println("IdController: error: ", err)
			return "", err
		}
		log.Println("IdController: user_id: ", user_id)
	case PHONE_LOGIN_TOKEN:
		//TODO: SAGAR  MAKE CHANGES FOR PHONE LOGIN TOKEN
		log.Println("IdController: phone login token")
	default:
		// This is the case for guest login
		user_id, err = idgenerator.GenerateUserID()
		if err != nil {
			log.Println("IdController: error: ", err)
			return "", err
		}
		log.Println("IdController: user_id: ", user_id)
	}

	return user_id, nil
}

func googleIdController(token string) (string, error) {
	log.Println("googleIdController: token: ", token)
	var exists bool
	exists, err := googledbcontroller.CheckIfUserIdExists(token, DB)
	if err != nil {
		log.Println("googleIdController: error: ", err)
		return "", err
	}
	if exists {
		log.Println("googleIdController: token present")
		refId, err := googledbcontroller.GetGoogleUserID(token, DB)
		if err != nil {
			log.Println("googleIdController: error: ", err)
			return "", err
		}
		log.Println("googleIdController: refId: ", refId)
		return refId, nil
	} else {
		// token not present
		log.Println("googleIdController: token not present hence generating a token now")
		refId, err := idgenerator.GenerateUserID()
		if err != nil {
			log.Println("googleIdController: error: ", err)
			return "", err
		}
		log.Println("googleIdController: refId: ", refId)
		err = googledbcontroller.InsertGoogleUserID(refId, token, DB)
		if err != nil {
			log.Println("googleIdController: error: ", err)
			return "", err
		}

	}
	return "", nil
}
