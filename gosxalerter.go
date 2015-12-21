package gosxalerter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Sound string
type ActivationType string

const (
	SoundDefault Sound = "'default'"
	SoundBasso   Sound = "Basso"
	SoundBlow    Sound = "Blow"
	SoundBottle  Sound = "Bottle"
	SoundFrog    Sound = "Frog"
	SoundFunk    Sound = "Funk"
	SoundGlass   Sound = "Glass"
	SoundHero    Sound = "Hero"
	SoundMorse   Sound = "Morse"
	SoundPing    Sound = "Ping"
	SoundPop     Sound = "Pop"
	SoundPurr    Sound = "Purr"
	SoundSosumi  Sound = "Sosumi"
	SoundTink    Sound = "Tink"
)

const (
	ActivationTypeClosed          ActivationType = "Closed"
	ActivationTypeTimeOut         ActivationType = "timeout"
	ActivationTypeContentsClicked ActivationType = "contentsClicked"
	ActivationTypeActionClicked   ActivationType = "actionClicked"
	ActivationTypeReplied         ActivationType = "replied"
)

type Alert struct {
	Options    *Options
	Activation *Activation
}
type Options struct {
	Message       string   //required
	Title         string   //optional
	Subtitle      string   //optional
	Sound         Sound    //optional
	Link          string   //optional
	Sender        string   //optional
	Group         string   //optional
	AppIcon       string   //optional
	ContentImage  string   //optional
	Actions       []string //optional
	Reply         bool
	CloseLabel    string
	DropdownLabel string
	Timeout       int
}

type Activation struct {
	At          string         `json:"activationAt"`
	Type        ActivationType `json:"activationType"`
	Value       string         `json:"activationValue"`
	DeliveredAt string         `json:"deliveredAt"`
	ValueIndex  string         `json:"activationValueIndex"`
}

func New(message string) *Alert {
	opts := &Options{
		Message: message,
		Reply:   false,
		Timeout: 0,
	}

	a := &Alert{
		Options: opts,
	}
	return a
}

func (a *Alert) DeliverAndWait() (*Activation, error) {
	activationChan, err := a.Deliver()
	if err != nil {
		return nil, fmt.Errorf("can not deliver - %s", err.Error())
	}
	activation := <-activationChan
	return activation, nil
}

func (a *Alert) Deliver() (chan *Activation, error) {
	name, args, err := buildCommand(a)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	cmd := exec.Command(name, args...)
	cmdOut, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	activation := make(chan *Activation, 1)
	cmdBytes, _ := ioutil.ReadAll(cmdOut)
	cmd.Wait()
	act := &Activation{}

	json.Unmarshal(cmdBytes, &act)
	a.Activation = act
	activation <- act

	return activation, nil
}

func buildCommand(a *Alert) (name string, arg []string, err error) {
	commandTuples := make([]string, 0)

	//check required commands
	if a.Options.Message == "" {
		return "", nil, errors.New("Please specifiy a proper message argument.")
	} else {
		commandTuples = append(commandTuples, []string{"-message", a.Options.Message}...)
	}

	//add closeLabel if found
	if a.Options.CloseLabel != "" {
		commandTuples = append(commandTuples, []string{"-closeLabel", a.Options.CloseLabel}...)
	}

	//add dropdownLabel if found
	if a.Options.DropdownLabel != "" {
		commandTuples = append(commandTuples, []string{"-dropdownLabel", a.Options.DropdownLabel}...)
	}

	//add title if found
	if len(a.Options.Actions) > 0 {
		commandTuples = append(commandTuples, []string{"-actions"}...)
		commandTuples = append(commandTuples, strings.Join(a.Options.Actions, ","))
	}

	//add Reply if found
	if a.Options.Reply == true {
		commandTuples = append(commandTuples, []string{"-reply"}...)
	}

	//add Reply if found
	if a.Options.Timeout > 0 {
		commandTuples = append(commandTuples, []string{"-timeout", strconv.Itoa(a.Options.Timeout)}...)
	}

	//add title if found
	if a.Options.Title != "" {
		commandTuples = append(commandTuples, []string{"-title", a.Options.Title}...)
	}

	//add subtitle if found
	if a.Options.Subtitle != "" {
		commandTuples = append(commandTuples, []string{"-subtitle", a.Options.Subtitle}...)
	}

	//add sound if specified
	if a.Options.Sound != "" {
		commandTuples = append(commandTuples, []string{"-sound", string(a.Options.Sound)}...)
	}

	//add group if specified
	if a.Options.Group != "" {
		commandTuples = append(commandTuples, []string{"-group", a.Options.Group}...)
	}

	//add appIcon if specified
	if a.Options.AppIcon != "" {
		commandTuples = append(commandTuples, []string{"-appIcon", a.Options.AppIcon}...)
	}

	//add contentImage if specified
	if a.Options.ContentImage != "" {
		commandTuples = append(commandTuples, []string{"-contentImage", a.Options.ContentImage}...)
	}

	//add sender if specified
	if strings.HasPrefix(strings.ToLower(a.Options.Sender), "com.") {
		commandTuples = append(commandTuples, []string{"-sender", a.Options.Sender}...)
	}

	if len(commandTuples) == 0 {
		return "", nil, errors.New("Please provide a Message and Type at a minimum.")
	}

	return finalPath, commandTuples, nil
}

const (
	executableFilename = "alerter"
	tempDirSuffix      = "gosxalterter"
)

var (
	finalPath string
)

func init() {
	if runtime.GOOS == "darwin" {
		if err := installAlerter(); err != nil {
			log.Fatal(err.Error())
		} else {
			finalPath = filepath.Join(os.TempDir(), executableFilename)
		}
	}
}

func installAlerter() error {
	finalPath := filepath.Join(os.TempDir(), executableFilename)

	//if alerter already installed no-need to re-install
	if _, err := os.Stat(finalPath); false == os.IsNotExist(err) {
		return nil
	}

	alerterB, _ := alerterBytes()
	err := ioutil.WriteFile(finalPath, alerterB, 0700)
	if err != nil {
		return errors.New("could not write alerter file")
	}

	err = os.Chmod(finalPath, 0755)
	if err != nil {
		return errors.New("could not make alerter executable")
	}

	return nil
}
