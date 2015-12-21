package main

import (
	"log"

	"github.com/vjeantet/gosx-alerter"
)

func main() {
	alert := gosxalerter.New("Deploy now on UAT ?")
	alert.Options.Actions = []string{"Now", "Later today", "Tomorrow"}
	alert.Options.AppIcon = "http://vjeantet.fr/images/logo.png"
	alert.Options.CloseLabel = "NO"
	alert.Options.DropdownLabel = "When ?"
	alert.Options.Title = "Alerter"
	alert.Options.Sound = gosxalerter.SoundHero

	alertActivation, err := alert.DeliverAndWait()
	if err != nil {
		log.Fatalln("error:", err)
	}

	log.Printf("Type : %s", alertActivation.Type)
	log.Printf("Value : %s", alertActivation.Value)
	log.Printf("Activated at : %s", alertActivation.At)
}
