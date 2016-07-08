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
	constDefaultConfigDir      = "."
	constDefaultConfigFilename = "fraudion.json"
	constDefaultLogDir         = "."
	// TODO: Should we keep the system defaulting to STDOUT or use this value?
	constDefaultLogFile = "fraudion.log"
)

var (
	argCLILogTo              = flag.String("logto", constDefaultLogDir, "Directory where to save the log file.")
	argCLILogFilename        = flag.String("logfile", constDefaultLogFile, "Log file's name.")
	argCLIConfigIn           = flag.String("cfgin", constDefaultConfigDir, "Directory where to search the config file.")
	argCLIConfigFilename     = flag.String("cfgfile", constDefaultConfigFilename, "Config file's name.")
	argCLIValidateConfigOnly = flag.Bool("cfgvalidate", false, "Validate config file only.")
)

func main() {

	// * Logger Setup
	log := marlog.MarLog
	log.Prefix = "FRAUDION"
	log.Flags = marlog.FlagLdate | marlog.FlagLtime | marlog.FlagLshortfile

	log.SetStamp("ERROR", "*STDOUT")
	log.SetStamp("DEBUG", "*STDOUT")
	log.SetStamp("INFO", "*STDOUT")

	system.State.StartUpTime = time.Now()

	log.LogS("INFO", "Fraudion started at "+system.State.StartUpTime.String())
	log.LogS("INFO", "Parsing CLI flags...")
	flag.Parse()

	logFileFullName := filepath.Join(*argCLILogTo, *argCLILogFilename)

	log.LogS("INFO", "Setting up the Log file \""+logFileFullName+"\"...")

	if logFile, err := os.OpenFile(logFileFullName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600); err != nil {
		log.LogO("ERROR", "Can't proceed. :( There was an Error opening/creating the Log file \""+logFileFullName+"\" ("+err.Error()+").", marlog.OptionFatal)
	} else {

		log.SetOutputHandle("MAINFILE", logFile)

		log.AddOuputHandles("INFO", "MAINFILE")
		log.LogS("INFO", "Started logging INFO messages to \""+logFileFullName+"\"")
		log.AddOuputHandles("DEBUG", "MAINFILE")
		log.LogS("INFO", "Started logging DEBUG messages to \""+logFileFullName+"\"")
		log.AddOuputHandles("ERROR", "MAINFILE")
		log.LogS("INFO", "Started logging ERROR messages to \""+logFileFullName+"\"")

		log.AddOuputHandles("VERBOSE", "MAINFILE")
		log.LogS("VERBOSE", "Started logging DEBUG messages to \""+logFileFullName+"\"")

	}

	// * Config Loading
	configFileFullName := filepath.Join(*argCLIConfigIn, *argCLIConfigFilename)

	if *argCLIValidateConfigOnly {

		log.LogS("INFO", "Validating configuration file...")

		configFile, err := os.Open(configFileFullName)

		if err != nil {

			log.LogS("INFO", "Could not open config file for validation: "+err.Error())

		} else {

			defer configFile.Close()

			if err := config.ValidateFromFile(configFile); err != nil {
				log.LogS("INFO", "Config file FAILED validation: "+err.Error())
			} else {
				log.LogS("INFO", "Config file PASSED validation. :)")
			}

		}

		os.Exit(0)
	}

	log.LogS("INFO", "Setting up Config loading from config file \""+configFileFullName+"\"...")

	// NOTE: If all goes well this puts Parsed/Verified configurations in config.Loaded
	if err := config.Load(configFileFullName); err != nil {
		log.LogO("ERROR", "Can't proceed. :( There was an error loading configurations from \""+configFileFullName+"\" ("+err.Error()+").", marlog.OptionFatal)
	}

	// * Monitored Softswitch Setup
	log.LogS("INFO", "Configuring the monitored Softswitch...")
	switch config.Loaded.Softswitch.Type {
	case softswitches.TypeAsterisk:

		log.LogS("DEBUG", "Softswitch type Asterisk")

		newSoftswitch := new(softswitches.Asterisk)
		newSoftswitch.Version = config.Loaded.Softswitch.Version
		newSoftswitch.CDRsSource = config.Loaded.Softswitch.CDRsSource

		switch config.Loaded.Softswitch.CDRsSource["type"] {
		case softswitches.CDRSourceDatabase:

			log.LogS("DEBUG", "CDRs Source is Database, DBMS \""+config.Loaded.Softswitch.CDRsSource["dbms"]+"\"")

			newSource := new(softswitches.CDRsSourceDatabase)
			newSource.DBMS = config.Loaded.Softswitch.CDRsSource["dbms"]
			newSource.UserName = config.Loaded.Softswitch.CDRsSource["user_name"]
			newSource.UserPassword = config.Loaded.Softswitch.CDRsSource["user_password"]
			newSource.DatabaseName = config.Loaded.Softswitch.CDRsSource["database_name"]
			newSource.TableName = config.Loaded.Softswitch.CDRsSource["table_name"]

			if err := newSource.Connect(); err != nil {
				log.LogO("ERROR", "Can't proceed. :( There was an Error (could not setup the Database connections pool)", marlog.OptionFatal)
			}

			newSoftswitch.CDRsSource = newSource

		default:
			// NOTE: This should not happen in the future because it's going to be validated in the configuration parsing/loading phase
			log.LogO("ERROR", "Can't proceed. :( There was an Error (unknown CDR Source type \""+config.Loaded.Softswitch.CDRsSource["type"]+"\" configured)", marlog.OptionFatal)
		}

		// TODO: Config Actions Chains/Actions/DataGroups here?
		// ...for each action type, get configs and create an action object of that type, fill configs with stuff from config.Loaded, do the same for the action chains and datagroups so that we can associate the softswitch with the action chain, etc

		log.LogS("INFO", "Softswitch is set up...")

		softswitches.Monitored = newSoftswitch

	default:
		// NOTE: This should not happen in the future because it's going to be validated in the configuration parsing/loading phase
		log.LogO("ERROR", "Can't proceed. :( There was an Error (unknown Softswitch type \""+config.Loaded.Softswitch.Type+"\" configured)", marlog.OptionFatal)
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

	if config.Loaded.Monitors.SimultaneousCalls.Enabled == true {

		log.LogS("INFO", "Monitor \"SimultaneousCalls\" is Enabled")

		ddMonitor := new(monitors.SimultaneousCalls)
		ddMonitor.Config = &config.Loaded.Monitors.SimultaneousCalls
		ddMonitor.Softswitch = softswitches.Monitored

		log.LogS("INFO", "Starting execution of monitor \"SimultaneousCalls\"...")
		go ddMonitor.Run()

	}

	// "Sleep!""
	log.LogS("INFO", "All set, main thread is going to sleep now...")
	for {

		// Main "thread" has to Sleep or else 100% CPU...
		time.Sleep(100000 * time.Hour)

	}

}
