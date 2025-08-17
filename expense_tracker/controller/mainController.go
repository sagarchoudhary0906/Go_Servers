package mainController

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// GET API Example :- curl -G --data-urlencode 'msg={"c":"expense","a":"create","d":{"amount":"10"}, "s":123}' http://localhost:8080/expense
// POST API Example :- curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' --data-urlencode 'msg={"c":"expense","a":"create","d":{"amount":"10"}, "s":123}' http://localhost:8080/expense
type JSONRequest struct {
	Controller string            `json:"c"`
	Action     string            `json:"a"`
	Data       map[string]string `json:"d"`
}

func MainController(w http.ResponseWriter, req *http.Request) {
	log.Println("MainController: start")
	var request JSONRequest

	ct := req.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			log.Printf("Error decoding JSON body: %v", err)
			return
		}
	} else {
		msg := req.FormValue("msg")
		if msg != "" {
			if err := json.Unmarshal([]byte(msg), &request); err != nil {
				http.Error(w, "invalid msg JSON", http.StatusBadRequest)
				log.Printf("Error unmarshalling msg: %v", err)
				return
			}
		}
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"response": request,
		"status":   "ok",
	})
	log.Println("MainController: end")
}
