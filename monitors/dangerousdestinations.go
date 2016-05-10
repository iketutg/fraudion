package monitors

import (
	"fmt"
	"regexp"
	//"strconv"
	"strings"
	"time"

	//"database/sql"
	//"os/exec"

	//"github.com/SlyMarbo/gmail"

	"github.com/andmar/fraudion/config"
	//"github.com/andmar/fraudion/logger"
	"github.com/andmar/fraudion/softswitch"
	"github.com/andmar/fraudion/state"
	//"github.com/andmar/fraudion/utils"

	"github.com/andmar/marlog"
)

// DangerousDestinationsRun ...
func DangerousDestinationsRun() {

	log := marlog.MarLog

	monitorConfig := config.Loaded.Monitors.DangerousDestinations
	cdrsSource := softswitch.Monitored.GetCDRsSource()

	log.LogS("INFO", "Started Monitor: DangerousDestinations")

	switch cdrsSource.(type) {
	case softswitch.CDRsSourceDatabase:

		cdrsSourceConverted := cdrsSource.(softswitch.CDRsSourceDatabase)

		log.LogS("INFO", fmt.Sprintf("CDRs Source is Database, DBMS:%s", cdrsSourceConverted.DBMS))

		for tickTime := range time.NewTicker(monitorConfig.ExecuteInterval).C { // NOTE: Replace "_" with "currentTime" and Log execution start time

			log.LogS("INFO", fmt.Sprintf("Monitor ticked at %s.", tickTime))

			cdrsSourceConverted, ok := cdrsSource.(softswitch.CDRsSourceDatabase)
			if ok == false {
				log.LogS("ERROR", "Something strange happened")
			}

			if err := cdrsSourceConverted.GetConnections().Ping(); err != nil {
				log.LogS("ERROR", "Something strange happened while checking if DB is alive")
			} else {

				log.LogS("INFO", "Database connection is A-Ok!")

				// Variables to hold "hit" information
				destinationHits := make(map[string]uint32)
				for _, prefix := range monitorConfig.PrefixList {
					destinationHits[prefix] = 0
				}
				hitValues := []string{}

				// NOTE: "guardDuration" and "guardTime" makes it so that when the service is restarted (maybe after an attack, to reset the monitors values), the CDRs, from the new start up time forward will only be considered from "startUpTime" - "guardTime" onwards, to try to prevent the system from redetecting the attack and reexecuting the associated Action Chain
				stringGuardDuration := "1h" // TODO: This value should also come from the configuration file
				guardDuration, err := time.ParseDuration(stringGuardDuration)
				if err != nil {
					log.LogS("ERROR", fmt.Sprintf("Something (%s) happened while trying to parse \"stringGuardTime\"", err.Error()))
				}
				guardTime := state.System.StartUpTime.Add(-guardDuration)
				durationSinceGuardTime := time.Now().Sub(guardTime)

				// Get CDRs from the Source
				// TODO: From here on what is done is Elastix2.3 specific, where the tests were made, so later we'll have to add some conditions to check what is the configured softswitch
				fmt.Println(fmt.Sprintf("SELECT * FROM cdr WHERE calldate >= DATE_SUB(CURDATE(), INTERVAL %v HOUR) AND calldate >= DATE_SUB(CURDATE(), INTERVAL %v HOUR) ORDER BY calldate DESC;", uint32(monitorConfig.ConsiderCDRsFromLast.Hours()), uint32(durationSinceGuardTime.Hours())))
				rows, err := cdrsSourceConverted.GetConnections().Query(fmt.Sprintf("SELECT * FROM cdr WHERE calldate >= DATE_SUB(CURDATE(), INTERVAL %v HOUR) ORDER BY calldate DESC;", uint32(monitorConfig.ConsiderCDRsFromLast.Hours())))
				if err != nil {
					log.LogS("ERROR", fmt.Sprintf("Something (%s) happened while trying to GET the CDRs", err.Error()))
				} else {

					for rows.Next() {

						var calldate, clid, src, dst, dcontext, channel, dstchannel, lastapp, lastdata, disposition, accountcode, uniqueid, userfield string
						var duration, billsec, amaflags uint32

						err := rows.Scan(&calldate, &clid, &src, &dst, &dcontext, &channel, &dstchannel, &lastapp, &lastdata, &duration, &billsec, &disposition, &amaflags, &accountcode, &uniqueid, &userfield)

						if err != nil {
							log.LogS("ERROR", fmt.Sprintf("Something (%s) happened while trying to GET the CDRs Values", err.Error()))
						} else {

							fmt.Println(calldate,
								clid,
								src,
								dst,
								//dcontext,
								//channel,
								//dstchannel,
								lastapp,
								lastdata,
								duration,
								billsec,
								disposition,
								//amaflags,
								//accountcode,
								//uniqueid,
								//userfield
							)

							// TODO: Should we match dials to more than one destination SIP/test/<number>&SIP/test/<number2>
							// TODO: Maybe the dial string match code should be from the interfaces because it's a softswitch specific thing
							// TODO: This is also Elastix2.3 specific, where the tests were made, so later we'll have to add some conditions to check what is the configured softswitch
							matchesDialString := regexp.MustCompile("(?:SIP|DAHDI)/[^@&]+/([0-9]+)") // NOTE: Supported dial string format
							matchedString := matchesDialString.FindString(lastdata)
							if lastapp != "Dial" /*|| strings.Contains(lastapp, "Local") || !test */ || matchedString == "" { // NOTE: Ignore if "lastapp" is not Dial and "lastdata" does not contain an expected dial string
								// TODO: Add some special Stamp to log this ignores?
								continue
							}

							dialedNumber := matchesDialString.FindStringSubmatch(lastdata)[1]

							if uint32(len(dialedNumber)) >= monitorConfig.MinimumNumberLength {

								for _, prefix := range monitorConfig.PrefixList {

									matchStringWithTag := monitorConfig.MatchRegex
									matchString := strings.Replace(matchStringWithTag, "__prefix__", prefix, 1)
									foundMatch, err := regexp.MatchString(matchString, lastdata)

									if err != nil {
										log.LogS("ERROR", fmt.Sprintf("Something (%s) happened while trying to match (found) a Prefix with regexp", err.Error()))
									}

									matchStringWithTag = monitorConfig.IgnoreRegex
									matchString = strings.Replace(matchStringWithTag, "__prefix__", prefix, 1)
									foundIgnore, err := regexp.MatchString(matchString, lastdata)

									if err != nil {
										log.LogS("ERROR", fmt.Sprintf("Something (%s) happened while trying to match (ignore) a Prefix with regexp", err.Error()))
									}

									if foundMatch == true && foundIgnore == false {
										destinationHits[prefix] = destinationHits[prefix] + 1
										hitValues = append(hitValues, dst)
									} else {
										// TODO: Add some special Stamp to log this ignores?
									}

								}

							}

						}

					}

					runActionChain := false
					for _, hits := range destinationHits {
						if hits >= monitorConfig.HitThreshold {
							runActionChain = true
						}
					}

					fmt.Println(runActionChain)

					/*
						actionChainGuardTime := stateTrigger.LastActionChainRunTime.Add(configs.General.DefaultActionChainHoldoffPeriod)

						if runActionChain && actionChainGuardTime.Before(time.Now()) && stateTrigger.ActionChainRunCount > 0 {

							stateTrigger.ActionChainRunCount = stateTrigger.ActionChainRunCount - 1

							actionChainName := monitorConfig.ActionChainName
							if actionChainName == "" {
								actionChainName = "*default"
							}

							logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("Running action chain: %s", actionChainName), false)
							stateTrigger.LastActionChainRunTime = time.Now()

							actionChain := configs.ActionChains.List[actionChainName]
							dataGroups := configs.DataGroups.List

							for _, v := range actionChain {

								if v.ActionName == "*email" {

									// TODO: Should we assert here that Email Action is enabled here or on config validation?

									body := fmt.Sprintf("Found:\n\n%v", hits)

									email := gmail.Compose("Fraudion ALERT: Dangerous Destinations!", fmt.Sprintf("\n\n%s", body))
									email.From = configs.Actions.Email.Username
									email.Password = configs.Actions.Email.Password
									fmt.Println(configs.Actions.Email.Username, configs.Actions.Email.Password)
									email.ContentType = "text/html; charset=utf-8"
									for _, dataGroupName := range v.DataGroupNames {
										fmt.Println(dataGroups[dataGroupName].EmailAddress)
										email.AddRecipient(dataGroups[dataGroupName].EmailAddress)
									}

									err := email.Send()
									if err != nil {
										fmt.Println(err.Error())
									}

								} else if v.ActionName == "*localcommand" {

									// TODO: Should we assert here that the run user of the process has "root" permissions?

									for _, dataGroupName := range v.DataGroupNames {

										command := exec.Command(dataGroups[dataGroupName].CommandName, dataGroups[dataGroupName].CommandArguments)

										err := command.Run()
										if err != nil {
											fmt.Println(err.Error())
										}

									}

								} else {

									fmt.Println("Unsupported Action!")

								}

							}

						}*/

				}

			}

		}

	default:
		log.LogS("ERROR", "Unknown CDRs Source")

	}

}
