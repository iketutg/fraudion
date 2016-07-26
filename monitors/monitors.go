package monitors

import (
	"fmt"
	"time"

	"os/exec"

	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/softswitches"
	"github.com/andmar/marlog"

	"github.com/SlyMarbo/gmail"
)

const (
	// RunModeNormal ...
	RunModeNormal = iota
	// RunModeInWarning ...
	RunModeInWarning
	// RunModeInAlarm ...
	RunModeInAlarm
)

// Monitor ...
type Monitor interface {
	Run()
}

// monitorBase ...
type monitorBase struct {
	Softswitch softswitches.Softswitch
}

// DangerousDestinations ...
type DangerousDestinations struct {
	monitorBase
	Config *config.MonitorDangerousDestinations
	State  StateDangerousDestinations
}

// SimultaneousCalls ...
type SimultaneousCalls struct {
	monitorBase
	Config *config.MonitorSimultaneousCalls
	State  StateSimultaneousCalls
}

type stateBase struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
	RunMode                int
}

// StateDangerousDestinations ...
type StateDangerousDestinations struct {
	stateBase
}

// StateSimultaneousCalls ...
type StateSimultaneousCalls struct {
	stateBase
}

func runActionChain(monitor Monitor, skipNonRecurrentActions bool, data map[string]string) error {

	log := marlog.MarLog

	var actionChainName string

	switch monitor.(type) {
	case *DangerousDestinations:
		monitor, ok := monitor.(*DangerousDestinations)
		if !ok {
			return fmt.Errorf("could not convert monitor value to a DangerousDestinations object")
		}
		actionChainName = monitor.Config.ActionChainName
		fmt.Println("Dangerous Destinations!")
	case *SimultaneousCalls:
		fmt.Println("Simultaneous Calls!")
	default:
		return fmt.Errorf("Unknown Monitor, this is probably a bug.")
	}

	log.LogS("DEBUG", "ActionChain to execute has name \""+actionChainName+"\"")

	actionChain, found := config.Loaded.ActionChains[actionChainName]
	if found == false {
		log.LogS("ERROR", "ActionChain not found by name")
		return fmt.Errorf("action chain not found by name")
	}

	log.LogS("INFO", "ActionChain found, looping Actions...")

	dataGroups := config.Loaded.DataGroups

	for _, action := range actionChain {

		switch action.ActionName {
		case "email":

			if config.Loaded.Actions.Email.Enabled == true && skipNonRecurrentActions == true {

				if config.Loaded.Actions.Email.Recurrent == false && skipNonRecurrentActions == true {
					log.LogS("INFO", "Action is non recurrent, skipping...")
				} else {

					log.LogS("INFO", "Executing e-mail action...")

					if data != nil && data["hits"] != "" {

						body := "Found:\n\n" + data["hits"]

						email := gmail.Compose("Fraudion ALERT @ "+config.Loaded.General.Hostname+": Monitor Dangerous Destinations!", "\n\n"+body)
						email.From = config.Loaded.Actions.Email.Username
						email.Password = config.Loaded.Actions.Email.Password
						email.ContentType = "text/html; charset=utf-8"

						for _, dataGroupName := range action.DataGroupNames {

							log.LogS("DEBUG", "Adding "+dataGroups[dataGroupName].EmailAddress+" as a recipient")

							email.AddRecipient(dataGroups[dataGroupName].EmailAddress)

						}

						err := email.Send()
						if err != nil {
							log.LogS("ERROR", "could not send the e-mail, an error ("+err.Error()+") ocurred")
						}

					} else {

						log.LogS("ERROR", "could not execute e-mail action because data was empty")

					}

				}

			}

		case "localcommands":

			if config.Loaded.Actions.LocalCommands.Enabled == true {

				if config.Loaded.Actions.LocalCommands.Recurrent == false && skipNonRecurrentActions == true {
					log.LogS("INFO", "Action is non recurrent, skipping...")
				} else {

					log.LogS("INFO", "Executing local command action...")

					for _, dataGroupName := range action.DataGroupNames {

						log.LogS("DEBUG", "Executing: "+dataGroups[dataGroupName].CommandName+" with arguments: "+dataGroups[dataGroupName].CommandArguments)

						command := exec.Command(dataGroups[dataGroupName].CommandName, dataGroups[dataGroupName].CommandArguments)

						err := command.Run()
						if err != nil {
							log.LogS("ERROR", "could not execute the command, an error ("+err.Error()+") ocurred")
						}

					}

				}

			}

		default:
			return fmt.Errorf("unsupported action in action chain")

		}

	}

	return nil
}
