package monitors

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/andmar/fraudion/utils"

	"github.com/andmar/marlog"
)

// Run ...
func (monitor *SmallDurationCalls) Run() {

	log := marlog.MarLog

	log.LogS("INFO", "Started Monitor SmallDurationCalls!")

	matches := func(destination string, args ...uint32) (string, bool, error) {

		if uint32(len(destination)) >= monitor.Config.MinimumNumberLength {

			hasPrefix, prefix := utils.FindIntlPrefix(destination)

			if !hasPrefix {
				return "", false, nil
			}

			matchStringWithTag := monitor.Config.MatchRegex
			matchString := strings.Replace(matchStringWithTag, "__prefix__", prefix, 1)

			foundMatch, err := regexp.MatchString(matchString, destination)
			if err != nil {
				log.LogS("ERROR", "an error  ("+err.Error()+") ocurred while trying to match a Prefix with regexp")
				return "", false, err
			}

			matchStringWithTag = monitor.Config.IgnoreRegex
			matchString = strings.Replace(matchStringWithTag, "__prefix__", prefix, 1)

			foundIgnore, err := regexp.MatchString(matchString, destination)
			if err != nil {
				log.LogS("ERROR", "an error ("+err.Error()+") ocurrerd while trying to match (to ignore) a Prefix with regexp")
				return "", false, err
			}

			callDuration, err := time.ParseDuration(strconv.FormatInt(int64(args[0]), 10) + "s")
			if err != nil {
				log.LogS("ERROR", "an error ("+err.Error()+") ocurrerd while trying to parse a duration")
				return "", false, err
			}

			if foundMatch == true && callDuration.Seconds() < monitor.Config.DurationThreshold.Seconds() && foundIgnore == false {
				return prefix, true, nil
			}

			return "", false, nil

		}

		return "", false, nil

	}

	log.LogS("DEBUG", "Setting up time Ticker with interval \""+monitor.Config.ExecuteInterval.String()+"\"")

	for tickTime := range time.NewTicker(monitor.Config.ExecuteInterval).C {

		log.LogS("INFO", "Monitor SmallDurationCalls ticked at "+tickTime.String())

		log.LogS("DEBUG", "Querying Softswitch for Hits (matches in CDRs) from the past \""+monitor.Config.ConsiderCDRsFromLast.String()+"\"...")

		hits, err := monitor.Softswitch.GetHits(matches, monitor.Config.ConsiderCDRsFromLast, true)
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

			log.LogS("INFO", "Checking if some Hits are above threshold \""+strconv.Itoa(int(monitor.Config.HitThreshold))+"\"")

			for _, v := range hits {

				if v.NumberOfHits > monitor.Config.HitThreshold {
					log.LogS("DEBUG", "Hits above threshold (first detected) \""+strconv.Itoa(int(monitor.Config.HitThreshold))+"\" on prefix "+v.Prefix+" found: "+fmt.Sprintf("%v", v.Destinations)+"!!")
					monitor.State.RunMode = RunModeInAlarm
					break
				}

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

			log.LogS("INFO", "RunMode after Hits check is "+runModeString)

			if monitor.State.RunMode != RunModeNormal {

				log.LogS("INFO", "Will execute action chain...")

				runActionChain(monitor, skipNonRecurrentActions, hits)

			}

		}

	}

}
