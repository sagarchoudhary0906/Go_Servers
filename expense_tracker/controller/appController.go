package mainController

import (
	"encoding/json"
	"log"
)

func AppController(req string) {
	var m map[string]json.RawMessage
	log.Println("AppController: start")
	if err := json.Unmarshal([]byte(req), &m); err != nil {
		log.Printf("invalid JSON: %v", err)
		return
	}

	c := string(m["c"])
	a := string(m["a"])
	d := string(m["d"])
	log.Printf("AppController: c=%q a=%q d=%v", c, a, d)

	log.Printf("AppController: c=%q a=%q d=%v", c, a, d)
}
