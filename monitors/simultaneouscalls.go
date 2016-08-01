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

	for tickTime := range time.NewTicker(monitor.Config.ExecuteInterval).C {

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

			}

		}

	}

}
