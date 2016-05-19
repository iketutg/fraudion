package main

import (
	"flag"
	"fmt"
	"os"
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
	log.Flags = marlog.FlagLdate | marlog.FlagLtime | marlog.FlagLshortfile

	// TODO: Error handling here?
	log.SetStamp("ERROR", "*STDOUT")
	log.SetStamp("DEBUG", "*STDOUT")
	log.SetStamp("INFO", "*STDOUT")

	// TODO: Remove this? This was just used for testing!
	log.SetStamp("VERBOSE", "*STDOUT")
	log.DeactivateStamps("VERBOSE")
	log.LogS("VERBOSE", "Test Verbose!")

	system.State.StartUpTime = time.Now()

	// TODO: Error handling here?
	log.LogS("INFO", "Starting Fraudion at "+system.State.StartUpTime.String())
	log.LogS("INFO", "Parsing CLI flags...")
	flag.Parse()

	logFileFullName := filepath.Join(*argCLILogTo, *argCLILogFilename)

	log.LogS("INFO", "Setting up the main log file \""+logFileFullName+"\"...")

	logFile, err := os.OpenFile(logFileFullName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.LogO("ERROR", "Can't proceed. :( There was an Error opening/creating the Log file \""+logFileFullName+"\" ("+err.Error()+").", marlog.OptionFatal)
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

	configFileFullName := filepath.Join(*argCLIConfigIn, *argCLIConfigFilename)
	if err := config.Load(configFileFullName, false); err != nil { // NOTE: Puts Verified configurations in config.Loaded
		log.LogO("ERROR", "Can't proceed. :( There was an error loading configurations from \""+configFileFullName+"\" ("+err.Error()+").", marlog.OptionFatal)
	}

	// * Monitored Softswitch Setup
	log.LogS("INFO", "Configuring the monitored Softswitch...")
	switch config.Loaded.Softswitch.Type {
	case softswitches.TypeAsterisk:

		log.LogS("DEBUG", "Softswitch type Asterisk")

		newSoftswitch := new(softswitches.Asterisk)
		newSoftswitch.Version = config.Loaded.Softswitch.Version

		sourceInfo, found := config.Loaded.CDRsSources[config.Loaded.Softswitch.CDRsSourceName]
		if found == false {
			log.LogO("ERROR", "Can't proceed. :( There was an error (could not find CDR Source with name \""+config.Loaded.Softswitch.CDRsSourceName+"\" in Loaded configurations)", marlog.OptionFatal) // TODO: This should not happen in the future because it's going to be validated in the configuration parsing/loading phase
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
		log.LogO("ERROR", "Can't proceed. :( There was an Error (unknown Softswitch type \""+config.Loaded.Softswitch.Type+"\" configured)", marlog.OptionFatal) // TODO: This should not happen in the future because it's going to be validated in the configuration parsing/loading phase
	}

	// * Config/Start Monitors
	log.LogS("INFO", "Configuring the monitors...")
	fmt.Println(config.Loaded)
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
