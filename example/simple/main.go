package main

import (
	"fmt"
	"log"
	"time"

	"github.com/vjeantet/gosx-alerter"
)

func main() {
	alert, _ := gosxalerter.New("Name this release please")
	alert.Options.Reply = true

	activationChan, err := alert.Deliver()
	if err != nil {
		log.Fatalln("error:", err)
	}

	// This is for example purpose, you can set a timeout options on an alert
	// when needed.
	select {
	case activation := <-activationChan:
		log.Printf("Type : %s", activation.Type)
		log.Printf("Value : %s", activation.Value)
		log.Printf("Activated at : %s", activation.At)

	case <-time.After(5 * time.Second):
		fmt.Println("BOOM!")
		alert.Close()
	}

}
