package main

import (
	"flag"
	"os"
	"time"

	"net/http"
	"path/filepath"

	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/monitors"
	"github.com/andmar/fraudion/softswitches"
	"github.com/andmar/fraudion/system"

	"github.com/andmar/marlog"

	_ "github.com/go-sql-driver/mysql"
)

const (
	constDefaultLogDir = "."
	// TODO: Should we keep the system defaulting to STDOUT or use this value as we do now?
	constDefaultLogFile = "fraudion.log"
)

var (
	argCLILogTo              = flag.String("logto", constDefaultLogDir, "Directory where to save the log file.")
	argCLILogFilename        = flag.String("logfile", constDefaultLogFile, "Log file's name.")
	argCLIConfigOrigin       = flag.String("cfgorigin", config.ConstDefaultOrigin, "Config data origin.")
	argCLIConfigIn           = flag.String("cfgin", config.ConstDefaultConfigDir, "Directory/URL where to get the config JSON data.")
	argCLIConfigFilename     = flag.String("cfgfile", config.ConstDefaultConfigFilename, "Config file's name (if -cfgorigin is \"file\").")
	argCLIConfigURL          = flag.String("cfgurl", config.ConstDefaultConfigURL, "Config URL (if -cfgorigin is \"url\").")
	argCLIValidateConfigOnly = flag.Bool("cfgvalidate", false, "Validate config file only.")
)

func main() {

	// * Logger Setup
	log := marlog.MarLog
	log.Prefix = "FRAUDION"
	log.Flags = marlog.FlagLdate | marlog.FlagLtime | marlog.FlagLshortfile

	// NOTE: What kind of "logs" will be available (that log to STDOUT)
	log.SetStamp("ERROR", "*STDOUT")
	log.SetStamp("DEBUG", "*STDOUT")
	log.SetStamp("INFO", "*STDOUT")
	log.SetStamp("VERBOSE", "*STDOUT")

	system.State.StartUpTime = time.Now()

	log.LogS("INFO", "Fraudion started at "+system.State.StartUpTime.String())
	log.LogS("INFO", "Parsing CLI flags...")

	flag.Parse()

	logFileFullName := filepath.Join(*argCLILogTo, *argCLILogFilename)

	log.LogS("INFO", "Setting up the Log file \""+logFileFullName+"\"...")

	// NOTE: Fraudion has to have permission to create the file in the folder it's being executed
	if logFile, err := os.OpenFile(logFileFullName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600); err != nil {
		log.LogO("ERROR", "Can't proceed. :( There was an Error opening/creating the Log file \""+logFileFullName+"\" ("+err.Error()+").", marlog.OptionFatal)
	} else {

		// NOTE: Do this to each file created to separate from previous log entries
		logFile.WriteString("\n")

		log.SetOutputHandle("MAINFILE", logFile)

		// NOTE: What kind of "logs" will be available (that log to "MAINFILE")
		log.AddOuputHandles("INFO", "MAINFILE")
		log.LogS("INFO", "Started logging INFO messages to \""+logFileFullName+"\" at "+system.State.StartUpTime.String())

		log.AddOuputHandles("DEBUG", "MAINFILE")
		log.LogS("INFO", "Started logging DEBUG messages to \""+logFileFullName+"\" at "+system.State.StartUpTime.String())

		log.AddOuputHandles("ERROR", "MAINFILE")
		log.LogS("INFO", "Started logging ERROR messages to \""+logFileFullName+"\" at "+system.State.StartUpTime.String())

		log.AddOuputHandles("VERBOSE", "MAINFILE")
		log.LogS("INFO", "Started logging VERBOSE messages to \""+logFileFullName+"\" at "+system.State.StartUpTime.String())

	}

	configOriginData := ""

	// * Config Validation
	if *argCLIValidateConfigOnly {

		if *argCLIConfigOrigin == config.ConstOriginFile {

			log.LogS("INFO", "Config Origin will be a file.")

			configOriginData = filepath.Join(*argCLIConfigIn, *argCLIConfigFilename)

			configFile, err := os.Open(configOriginData)
			if err != nil {
				log.LogS("ERROR", "Could not open config file for validation: "+err.Error())
			} else {

				defer configFile.Close()

				log.LogS("INFO", "Validating configuration JSON...")
				if err := config.ValidateFromFile(configFile); err != nil {
					log.LogS("INFO", "Config file FAILED validation: "+err.Error())
				} else {
					log.LogS("INFO", "Config file PASSED validation. :)")
				}

			}

		} else {

			log.LogS("INFO", "Config Origin will be a URL.")

			configOriginData = config.ConstDefaultConfigURL

			if r, err := http.Get(*argCLIConfigURL); err != nil {
				log.LogS("ERROR", "Could not fetch config json from URL for validation: "+err.Error())
			} else {

				defer r.Body.Close()

				log.LogS("INFO", "Validating configuration JSON...")
				if err := config.ValidateFromURL(r.Body); err != nil {
					log.LogS("INFO", "Config file FAILED validation: "+err.Error())
				} else {
					log.LogS("INFO", "Config file PASSED validation. :)")
				}

			}

		}

		os.Exit(0)

	}

	// * Config Loading
	originLabel := "File"
	configOriginData = *argCLIConfigFilename
	if *argCLIConfigOrigin == config.ConstOriginURL {
		originLabel = "URL"
		configOriginData = *argCLIConfigURL
	}
	log.LogS("INFO", "Setting up Config loading from "+originLabel+" (\""+configOriginData+"\")")

	// NOTE: If all goes well this puts Parsed/Verified configurations in config.Loaded
	if err := config.Load(configOriginData, *argCLIConfigOrigin); err != nil {
		log.LogO("ERROR", "Can't proceed. :( There was an error loading configurations ("+err.Error()+").", marlog.OptionFatal)
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

		log.LogS("INFO", "Softswitch is set up...")

		softswitches.Monitored = newSoftswitch

	default:
		// NOTE: This should not happen in the future because it's going to be validated in the configuration parsing/loading phase
		log.LogO("ERROR", "Can't proceed. :( There was an Error (unknown Softswitch type \""+config.Loaded.Softswitch.Type+"\" configured)", marlog.OptionFatal)
	}

	// * Config/Start Monitors
	log.LogS("INFO", "Configuring the monitors...")

	if config.Loaded.Monitors.DangerousDestinations.Enabled == true {

		log.LogS("INFO", "Monitor \"DangerousDestinations\" is Enabled")

		monitor := new(monitors.DangerousDestinations)
		monitor.Config = &config.Loaded.Monitors.DangerousDestinations
		monitor.Softswitch = softswitches.Monitored

		log.LogS("INFO", "Starting execution of monitor \"DangerousDestinations\"...")
		go monitor.Run()

	}

	if config.Loaded.Monitors.SimultaneousCalls.Enabled == true {

		log.LogS("INFO", "Monitor \"SimultaneousCalls\" is Enabled")

		monitor := new(monitors.SimultaneousCalls)
		monitor.Config = &config.Loaded.Monitors.SimultaneousCalls
		monitor.Softswitch = softswitches.Monitored

		log.LogS("INFO", "Starting execution of monitor \"SimultaneousCalls\"...")
		go monitor.Run()

	}

	log.LogS("INFO", "All set, main thread is going to sleep now...")
	for {

		// NOTE: Main "thread" has to Sleep or else 100% CPU...
		time.Sleep(100000 * time.Hour)

	}

}
