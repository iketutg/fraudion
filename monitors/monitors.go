package monitors

import (
	"fmt"
	"reflect"
	"time"

	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/softswitches"
	"github.com/andmar/marlog"
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

func RunActionChain(monitor Monitor, skipNonRecurrentActions bool, data map[string]string) error {

	log := marlog.MarLog

	log.LogS("INFO", "One... Two...")

	a := reflect.TypeOf(monitor)
	b := reflect.TypeOf(a)

	fmt.Println(a, b)

	switch monitor.(type) {
	case *DangerousDestinations:
		fmt.Println("Dangerous Destinations!")

	}

	// actionChainName := monitor.Config.ActionChainName
	// if actionChainName == "" {
	// 	actionChainName = "default"
	// }
	//
	// log.LogS("DEBUG", "ActionChain to execute has name \""+actionChainName+"\"")
	//
	// actionChain, found := config.Loaded.ActionChains[actionChainName]
	// if found == false {
	// 	log.LogS("ERROR", "ActionChain not found by name")
	// } else {
	//
	// 	log.LogS("INFO", "ActionChain found, looping Actions...")
	//
	// 	for _, v := range actionChain {
	//
	// 		if v.ActionName == "*email" && config.Loaded.Actions.Email.Enabled == true {
	//
	// 			if config.Loaded.Actions.Email.Recurrent == false && skipNonRecurrentActions == true {
	// 				log.LogS("INFO", "Action is non recurrent, skipping...")
	// 			} else {
	//
	// 				log.LogS("INFO", "Executing e-mail action")
	//
	// 				body := fmt.Sprintf("Found:\n\n%v", hits)
	// 				//body := fmt.Sprintf("Test!")
	//
	// 				email := gmail.Compose("Fraudion ALERT: Dangerous Destinations!", fmt.Sprintf("\n\n%s", body))
	// 				email.From = config.Loaded.Actions.Email.Username
	// 				email.Password = config.Loaded.Actions.Email.Password
	// 				email.ContentType = "text/html; charset=utf-8"
	//
	// 				for _, dataGroupName := range v.DataGroupNames {
	//
	// 					email.AddRecipient(dataGroups[dataGroupName].EmailAddress)
	//
	// 				}
	//
	// 				err := email.Send()
	// 				if err != nil {
	// 					log.LogS("ERROR", "Could not send the e-mail, an error ("+err.Error()+") ocurred")
	// 				}
	//
	// 			}
	//
	// 		} else if v.ActionName == "*local_commands" && config.Loaded.Actions.LocalCommands.Enabled == true {
	//
	// 			if config.Loaded.Actions.Email.Recurrent == false && skipNonRecurrentActions == true {
	// 				log.LogS("ERROR", "Action is non recurrent, skipping...")
	// 			} else {
	//
	// 				log.LogS("INFO", "Executing local command action")
	//
	// 				for _, dataGroupName := range v.DataGroupNames {
	//
	// 					command := exec.Command(dataGroups[dataGroupName].CommandName, dataGroups[dataGroupName].CommandArguments)
	//
	// 					err := command.Run()
	// 					if err != nil {
	// 						log.LogS("ERROR", "Could not execute the command, an error ("+err.Error()+") ocurred")
	// 					}
	//
	// 				}
	//
	// 			}
	//
	// 		} else {
	//
	// 			log.LogS("ERROR", "Unsupported Action")
	//
	// 		}
	//
	// 	}
	//
	// }

	return nil
}
