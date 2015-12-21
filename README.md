# gosx-alerter
gosx-alerter is a Go package for delivering desktop alert notifications to OSX 10.8 or higher

[![GoDoc](http://godoc.org/github.com/vjeantet/gosx-alerter?status.png)](http://godoc.org/github.com/vjeantet/gosx-alerter)

```go
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
```

![](../master/alerter-actions.png?raw=true)

```go
    alert := gosxalerter.New("Name this release please")
    alert.Options.Title = "Alerter"
    alert.Options.Reply = true

    alertActivation, err := alert.DeliverAndWait()
    if err != nil {
        log.Fatalln("error:", err)
    }

    log.Printf("Type : %s", alertActivation.Type)
    log.Printf("Value : %s", alertActivation.Value)
    log.Printf("Activated at : %s", alertActivation.At)
```
![](../master/alerter-reply.png?raw=true)
![](../master/alerter-replytext.png?raw=true)