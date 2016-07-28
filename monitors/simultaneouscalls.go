package monitors

import (
	"strconv"
	"time"

	"github.com/andmar/marlog"
)

// Run ...
func (monitor *SimultaneousCalls) Run() {

	log := marlog.MarLog

	log.LogS("INFO", "Started Monitor SimultaneousCalls!")

	log.LogS("DEBUG", "Setting up time Ticker with interval \""+monitor.Config.ExecuteInterval.String()+"\"")

	for tickTime := range time.NewTicker(monitor.Config.ExecuteInterval).C { // NOTE: Replace "_" with "currentTime" and Log execution start time

		log.LogS("INFO", "Monitor SimultaneousCalls ticked at "+tickTime.String())

		log.LogS("DEBUG", "Querying Softswitch for Current Active Calls...")

		numberOfCalls, err := monitor.Softswitch.GetCurrentActiveCalls(monitor.Config.MinimumNumberLength)
		if err != nil {
			log.LogS("ERROR: ", err.Error())
		} else {

			// NOTE: This block has to be here because we reset the value of monitor.State.RunMode below, this catches state changes
			skipNonRecurrentActions := false
			if monitor.State.RunMode != RunModeNormal {
				skipNonRecurrentActions = true
			}

			// NOTE: Resets RunMode in each Tick so that the System can detect when it's out of an alarm situation
			monitor.State.RunMode = RunModeNormal

			log.LogS("INFO", "Current active Calls "+strconv.Itoa(int(numberOfCalls)))

			if numberOfCalls > monitor.Config.HitThreshold {
				log.LogS("INFO", "Number above threshold \""+strconv.Itoa(int(monitor.Config.HitThreshold))+"\"!!")
				monitor.State.RunMode = RunModeInAlarm
			}

			runModeString := ""
			switch monitor.State.RunMode {
			case RunModeInWarning:
			case RunModeInAlarm:
				runModeString = "Alarm/Warning"
				log.LogS("DEBUG", "System is in Alarm/Warning")
			default:
				runModeString = "Normal"
				log.LogS("DEBUG", "System detected nothing. :)")
			}

			log.LogS("INFO", "RunMode after Simultaneous Calls check is "+runModeString)

			if monitor.State.RunMode != RunModeNormal {

				log.LogS("INFO", "Will execute action chain...")

				runActionChain(monitor, skipNonRecurrentActions, numberOfCalls)

				// actionChainName := monitor.Config.ActionChainName
				// if actionChainName == "" {
				// 	actionChainName = "default"
				// }
				// dataGroups := config.Loaded.DataGroups
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
				// 				body := fmt.Sprintf("Found:\n\n%v", numberOfCalls)
				// 				//body := fmt.Sprintf("Test!")
				//
				// 				email := gmail.Compose("Fraudion ALERT: Simultaneous Calls!", fmt.Sprintf("\n\n%s", body))
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

			}

		}

	}

}
