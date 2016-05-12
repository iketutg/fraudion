package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"path/filepath"

	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/monitors"
	"github.com/andmar/fraudion/softswitches"
	"github.com/andmar/fraudion/system"

	"github.com/andmar/marlog"

	_ "github.com/go-sql-driver/mysql"
)

const (
	constDefaultConfigDir      = "." // NOTE: The system defaults to search the config file "fraudion.json" on the "run" directory
	constDefaultConfigFilename = "fraudion.json"
	constDefaultLogDir         = "."
	constDefaultLogFile        = "fraudion.log" // TODO: Should we keep the system defaulting to STDOUT or use this value?
)

var (
	argCLILogTo          = flag.String("logto", constDefaultLogDir, "<help message for 'logto'>")
	argCLILogFilename    = flag.String("logfile", constDefaultLogFile, "<help message for 'logfile'>")
	argCLIConfigIn       = flag.String("cfgin", constDefaultConfigDir, "<help message for 'cfgin'>")
	argCLIConfigFilename = flag.String("cfgfile", constDefaultConfigFilename, "<help message for 'cfgfile'>")
)

func main() {

	// * Logger Setup
	log := marlog.MarLog
	log.Prefix = "FRAUDION"
	log.Flags = marlog.FlagLdate | marlog.FlagLtime

	// TODO: Error handling here?
	log.SetStamp("ERROR", "*STDOUT")
	log.SetStamp("DEBUG", "*STDOUT")
	log.SetStamp("INFO", "*STDOUT")

	system.State.StartUpTime = time.Now()

	// TODO: Error handling here?
	log.LogS("INFO", "Starting Fraudion at "+system.State.StartUpTime.String())
	log.LogS("INFO", "Parsing CLI flags...")
	flag.Parse()

	// TODO: This should default to constDefaultLogFile, maybe even handle a flag to disable logging
	if strings.ToLower(*argCLILogTo) != "" && strings.ToLower(*argCLILogFilename) != "" {

		logFileFullName := filepath.Join(*argCLILogTo, *argCLILogFilename)

		log.LogS("INFO", "Setting up the main log file \""+logFileFullName+"\"...")

		logFile, err := os.OpenFile(logFileFullName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {

			log.LogO("ERROR", "Can't proceed. :( There was an Error ("+err.Error()+") opening/creating the Log file \""+logFileFullName+"\".", marlog.OptionFatal)

		} else {

			// TODO: Error handling here?
			log.SetOutputHandle("MAINFILE", logFile)

			log.AddOuputHandles("INFO", "MAINFILE")
			log.LogS("INFO", "Started logging INFO messages to \""+logFileFullName+"\"")
			log.AddOuputHandles("DEBUG", "MAINFILE")
			log.LogS("INFO", "Started logging DEBUG messages to \""+logFileFullName+"\"")
			log.AddOuputHandles("ERROR", "MAINFILE")
			log.LogS("INFO", "Started logging ERROR messages to \""+logFileFullName+"\"")

		}

	}

	// TODO: Debug: Remove this?
	// fmt.Println("Log:", log)

	if err := config.Load(*argCLIConfigIn, *argCLIConfigFilename, false); err != nil {
		log.LogO("ERROR", "Can't proceed. :( There was an Error ("+err.Error()+")", marlog.OptionFatal) // TODO: This has to be changed becase config.Validate() returns an array/slice of errors
	}

	// * Monitored Softswitch Setup
	log.LogS("INFO", "Configuring the monitored Softswitch...")
	switch config.Loaded.Softswitch.System {
	case softswitches.TypeAsterisk:

		log.LogS("DEBUG", "Softswitch type Asterisk")

		newSoftswitch := new(softswitches.Asterisk)
		newSoftswitch.Version = config.Loaded.Softswitch.Version

		sourceInfo, found := config.Loaded.CDRsSources[config.Loaded.Softswitch.CDRsSource]
		if found == false {
			log.LogO("ERROR", "Can't proceed. :( There was an Error (could not find CDR Source with name \""+config.Loaded.Softswitch.CDRsSource+"\" in Loaded configurations)", marlog.OptionFatal) // TODO: This should not happen in the future because it's going to be validated in the configuration parsing/loading phase
		}

		switch sourceInfo["type"] {
		case softswitches.CDRSourceDatabase:

			log.LogS("DEBUG", "CDRs Source is Database, DBMS \""+sourceInfo["dbms"]+"\"")

			newSource := new(softswitches.CDRsSourceDatabase)
			newSource.DBMS = sourceInfo["dbms"]
			newSource.UserName = sourceInfo["user_name"]
			newSource.UserPassword = sourceInfo["user_password"]
			newSource.DatabaseName = sourceInfo["database_name"]
			newSource.TableName = sourceInfo["table_name"]

			if err := newSource.Connect(); err != nil {
				log.LogO("ERROR", "Can't proceed. :( There was an Error (could not setup the Database connections pool)", marlog.OptionFatal)
			}

			newSoftswitch.CDRsSource = newSource

		default:
			log.LogO("ERROR", "Can't proceed. :( There was an Error (unknown CDR Source type \""+sourceInfo["type"]+"\" configured)", marlog.OptionFatal) // TODO: This should not happen in the future because it's going to be validated in the configuration parsing/loading phase
		}

		// TODO: Config Actions Chains/Actions/DataGroups here?
		// ...for each action type, get configs and create an action object of that type, fill configs with stuff from config.Loaded, do the same for the action chains and datagroups so that we can associate the softswitch with the action chain, etc

		log.LogS("INFO", "Softswitch is set up...")

		softswitches.Monitored = newSoftswitch

	default:
		log.LogO("ERROR", "Can't proceed. :( There was an Error (unknown Softswitch type \""+config.Loaded.Softswitch.System+"\" configured)", marlog.OptionFatal) // TODO: This should not happen in the future because it's going to be validated in the configuration parsing/loading phase
	}

	// fmt.Println("\nLoaded Configurations:")
	// fmt.Println(config.Loaded.General)
	// fmt.Println(config.Loaded.Softswitch)
	// fmt.Println(config.Loaded.CDRsSources)
	// fmt.Println(config.Loaded.Monitors.DangerousDestinations, config.Loaded.Monitors.ExpectedDestinations, config.Loaded.Monitors.SimultaneousCalls, config.Loaded.Monitors.SmallDurationCalls)
	// fmt.Println(config.Loaded.Actions.Email, config.Loaded.Actions.Call, config.Loaded.Actions.HTTP, config.Loaded.Actions.LocalCommands)
	// fmt.Println(config.Loaded.ActionChains)
	// fmt.Println(config.Loaded.DataGroups)
	// fmt.Println()
	//
	// fmt.Println("Loaded CDRs Sources:")
	// fmt.Println(config.Loaded.CDRsSources)
	// fmt.Println()
	//
	// fmt.Println("Softswitch:")
	// fmt.Println(softswitches.Monitored)
	// fmt.Println()
	//
	// fmt.Println("CDRs Source:")
	// fmt.Println(softswitches.Monitored.GetCDRsSource())
	// fmt.Println()

	// * Config/Start Monitors
	log.LogS("INFO", "Configuring the monitors...")
	if config.Loaded.Monitors.DangerousDestinations.Enabled == true {

		log.LogS("INFO", "Monitor \"DangerousDestinations\" is Enabled")

		ddMonitor := new(monitors.DangerousDestinations)
		ddMonitor.Config = &config.Loaded.Monitors.DangerousDestinations
		ddMonitor.Softswitch = softswitches.Monitored

		log.LogS("INFO", "Starting execution of monitor \"DangerousDestinations\"...")
		go ddMonitor.Run()

	}

	log.LogS("INFO", "All set, main thread is going to sleep now...")

	// Sleep!
	for {

		// Main "thread" has to Sleep or else 100% CPU...
		time.Sleep(100000 * time.Hour)

	}

}
