package mainController

import (
	"encoding/json"
	"log"
	"net/http"
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
	msg := req.FormValue("msg")
	log.Println("MainController: msg: ", msg)

	status := "not ok"
	if msg != "" {
		AppController(msg)
		status = "ok"
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"response": request,
		"status":   status,
	})
	log.Println("MainController: end")
}
