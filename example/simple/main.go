package main

import (
	"log"

	"github.com/vjeantet/gosx-alerter"
)

func main() {
	alert := gosxalerter.New("kmlj")
	alert.Options.Actions = []string{"YES", "MAYBE"}
	alert.Options.CloseLabel = "NO"
	alert.Options.DropdownLabel = "Actions"
	alert.Options.Title = "Alerter"
	alert.Options.Sound = gosxalerter.SoundHero
	alert.Options.Timeout = 10

	alertActivation, err := alert.DeliverAndWait()
	if err != nil {
		log.Fatalln("error:", err)
	}

	log.Printf("Type : %s", alertActivation.Type)
	log.Printf("Value : %s", alertActivation.Value)
	log.Printf("Activated at : %s", alertActivation.At)
}
