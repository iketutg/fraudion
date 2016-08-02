package monitors

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
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
	// ActionEmail ...
	ActionEmail = "*email"
	// ActionLocalCommands ...
	ActionLocalCommands = "*local_commands"
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

// ExpectedDestinations ...
type ExpectedDestinations struct {
	monitorBase
	Config *config.MonitorExpectedDestinations
	State  StateExpectedDestinations
}

// SmallDurationCalls ...
type SmallDurationCalls struct {
	monitorBase
	Config *config.MonitorSmallDurationCalls
	State  StateSmallDurationCalls
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

// StateExpectedDestinations ...
type StateExpectedDestinations struct {
	stateBase
}

// StateSmallDurationCalls ...
type StateSmallDurationCalls struct {
	stateBase
}

var runActionChainmutex = &sync.Mutex{}

func runActionChain(monitor Monitor, skipNonRecurrentActions bool, data interface{}) error {

	runActionChainmutex.Lock()

	log := marlog.MarLog

	// TODO: This solution is a little bit weird, seems wrong...
	// NOTE: Only one will work...
	monitorDangerousDestinations, okDD := monitor.(*DangerousDestinations)
	monitorSimultaneousCalls, okSC := monitor.(*SimultaneousCalls)
	monitorExpectedDestinations, okED := monitor.(*ExpectedDestinations)
	if !okDD && !okSC && !okED {
		return fmt.Errorf("unable to detect monitor that tried to run the action chain")
	}

	var actionChainName string

	if okDD {
		actionChainName = monitorDangerousDestinations.Config.ActionChainName
	} else if okSC {
		actionChainName = monitorSimultaneousCalls.Config.ActionChainName
	} else {
		actionChainName = monitorExpectedDestinations.Config.ActionChainName
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
		case ActionEmail:

			if config.Loaded.Actions.Email.Enabled == true {

				if config.Loaded.Actions.Email.Recurrent == false && skipNonRecurrentActions == true {
					log.LogS("INFO", "Action is non recurrent, skipping...")
				} else {

					log.LogS("INFO", "Executing e-mail action...")

					subject := "ALERT @ " + config.Loaded.General.Hostname + ": "
					body := ""

					if okDD {

						dataAsserted, ok := data.(map[string]*softswitches.Hits)
						if !ok {
							log.LogS("ERROR", "could not convert data to e-mail action usable object")
						} else {

							prefixes := ""
							for key := range dataAsserted {
								prefixes = prefixes + key + ", "
							}
							prefixes = strings.TrimSuffix(prefixes, ", ")

							subject = subject + "Dangerous Destinations!"
							body = "Suspicious calls to:\n\n" + prefixes

						}

					} else if okSC {

						dataAsserted, ok := data.(uint32)
						if !ok {
							log.LogS("ERROR", "could not convert data to e-mail action usable object")
						} else {

							subject = subject + "Simultaneous Calls"
							body = "Currently active calls:\n\n" + strconv.Itoa(int(dataAsserted))

						}

					} else {

						dataAsserted, ok := data.(map[string]*softswitches.Hits)
						if !ok {
							log.LogS("ERROR", "could not convert data to e-mail action usable object")
						} else {

							prefixes := ""
							for key := range dataAsserted {
								prefixes = prefixes + key + ", "
							}
							prefixes = strings.TrimSuffix(prefixes, ", ")

							subject = subject + "Expected Destinations!"
							body = "Suspicious calls to:\n\n" + prefixes

						}

					}

					email := gmail.Compose(subject, "\n\n"+body)
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

				}

			}

		case ActionLocalCommands:

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

	runActionChainmutex.Unlock()

	return nil
}
