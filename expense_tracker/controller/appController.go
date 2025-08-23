package mainController

import (
	"encoding/json"
	"log"
)

const GET_USER_ID = "guid"

func AppController(req string) {
	var m map[string]json.RawMessage
	log.Println("AppController: start")
	if err := json.Unmarshal([]byte(req), &m); err != nil {
		log.Printf("invalid JSON: %v", err)
		return
	}

	var a string
	if err := json.Unmarshal(m["a"], &a); err != nil {
		// handle error
	}
	var c string
	var d map[string]string
	_ = json.Unmarshal(m["c"], &c)
	_ = json.Unmarshal(m["d"], &d)
	log.Printf("AppController: c=%q a=%q d=%v  res = %v", c, a, d, a == GET_USER_ID)

	if a == GET_USER_ID {
		IdController(d)
	}
}
