package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"path/filepath"

	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/fraudion"
	"github.com/andmar/fraudion/logger"

	"github.com/andmar/marlog"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// NOTE: The system defaults to search the config file "fraudion.json" on the "run" directory
	constDefaultConfigDir      = "."
	constDefaultConfigFilename = "fraudion.json"
	constDefaultLogDir         = "."
	constDefaultLogFile        = "fraudion.log" // TODO: Should we keep the system defaulting to STDOUT or use this value?

)

var (
	argCLILogTo          = flag.String("logto", constDefaultLogDir, "<help message for 'logto'>")
	argCLILogFilename    = flag.String("logfile", constDefaultLogFile, "<help message for 'logfile'>")
	argCLIConfigIn       = flag.String("cfgin", constDefaultConfigDir, "<help message for 'cfgin'>")
	argCLIConfigFilename = flag.String("cfgfile", constDefaultConfigFilename, "<help message for 'cfgfile'>")

	// TODO: This is temporary only, this information will always come from the config file
	argCLIDBPass = flag.String("dbpass", "", "<help message for 'dbpass'>")
)

func main() {

	// TODO: This is to be removed
	fraudion := fraudion.Global // NOTE: fraudion.Global (and it's pointers) is (are) initialized on fraudion's package init() function

	// Logger Setup
	log := marlog.MarLog
	log.Prefix = "FRAUDION"
	log.Flags = marlog.FlagLdate | marlog.FlagLtime | marlog.FlagLlongfile

	// TODO: Error handling here?
	log.SetStamp("ERROR", "*STDOUT")
	log.SetStamp("DEBUG", "*STDOUT")
	log.SetStamp("INFO", "*STDOUT")

	fraudion.StartUpTime = time.Now()

	// TODO: Error handling here?
	log.LogS("INFO", fmt.Sprintf("Starting Fraudion at %s", fraudion.StartUpTime))
	log.LogS("INFO", "Parsing CLI flags...")
	flag.Parse()

	// TODO: This should default to constDefaultLogFile, maybe even handle a flag to disable logging
	if strings.ToLower(*argCLILogTo) != "" && strings.ToLower(*argCLILogFilename) != "" {

		logFile, err := os.OpenFile(filepath.Join(*argCLILogTo, *argCLILogFilename), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {

			log.LogO("ERROR", fmt.Sprintf("Can't start, there was a problem (%s) opening the Log file. :(", err), marlog.OptionFatal)

		} else {

			// TODO: Error handling here
			log.SetOutputHandle("MAINFILE", logFile)

			log.AddOuputHandles("INFO", "MAINFILE")
			log.AddOuputHandles("DEBUG", "MAINFILE")
			log.AddOuputHandles("ERROR", "MAINFILE")

			log.LogS("INFO", fmt.Sprintf("Started logging to \"%s\"", filepath.Join(*argCLILogTo, *argCLILogFilename)))

			logger.Log.SetHandles(logFile, logFile, logFile, logFile) // NOTE: Overwrite the default handles on the Logger object.

		}

	}

	if err := config.Parse(*argCLIConfigIn, *argCLIConfigFilename); err != nil {
		log.LogO("ERROR", err.Error(), marlog.OptionFatal)
	}

	// fmt.Println(config.Parsed.General)
	// fmt.Println(config.Parsed.Softswitch)
	// fmt.Println(config.Parsed.CDRsSources)
	// fmt.Println(config.Parsed.Monitors.DangerousDestinations, config.Parsed.Monitors.ExpectedDestinations, config.Parsed.Monitors.SimultaneousCalls, config.Parsed.Monitors.SmallDurationCalls)
	// fmt.Println(config.Parsed.Actions.Email, config.Parsed.Actions.Call, config.Parsed.Actions.HTTP, config.Parsed.Actions.LocalCommands)
	// fmt.Println(config.Parsed.ActionChains)
	// fmt.Println(config.Parsed.DataGroups)

	//if err := config.Validate(); err != nil {
	//log.LogO("ERRORS:", err.Error(), marlog.OptionFatal) // TODO: This has to be changed becase config.Validate() returns an array/slice of errors
	//}

	if err := config.Load(); err != nil {
		log.LogO("ERROR", err.Error(), marlog.OptionFatal) // TODO: This has to be changed becase config.Validate() returns an array/slice of errors
	}

	fmt.Println(config.Loaded.General)
	fmt.Println(config.Loaded.Softswitch)
	fmt.Println(config.Loaded.CDRsSources)
	fmt.Println(config.Loaded.Monitors.DangerousDestinations, config.Loaded.Monitors.ExpectedDestinations, config.Loaded.Monitors.SimultaneousCalls, config.Loaded.Monitors.SmallDurationCalls)
	fmt.Println(config.Loaded.Actions.Email, config.Loaded.Actions.Call, config.Loaded.Actions.HTTP, config.Loaded.Actions.LocalCommands)
	fmt.Println(config.Loaded.ActionChains)
	fmt.Println(config.Loaded.DataGroups)

	os.Exit(-1)

	/*
		configs, err := config.Load(configsJSON)
		if err != nil {
			logger.Log.Write(logger.ConstLoggerLevelError, fmt.Sprintf("There was an error (%s) validating/loading the Fraudion configuration", err), true)
		}*/

	/*if configs.Triggers.DangerousDestinations.MaxActionChainRunCount != 0 {
		fraudion.State.Triggers.StateDangerousDestinations.ActionChainRunCount = configs.Triggers.DangerousDestinations.MaxActionChainRunCount
	} else {
		fraudion.State.Triggers.StateDangerousDestinations.ActionChainRunCount = configs.General.DefaultActionChainRunCount
	}

	if configs.Triggers.ExpectedDestinations.MaxActionChainRunCount != 0 {
		fraudion.State.Triggers.StateExpectedDestinations.ActionChainRunCount = configs.Triggers.ExpectedDestinations.MaxActionChainRunCount
	} else {
		fraudion.State.Triggers.StateDangerousDestinations.ActionChainRunCount = configs.General.DefaultActionChainRunCount
	}

	if configs.Triggers.SimultaneousCalls.MaxActionChainRunCount != 0 {
		fraudion.State.Triggers.StateSimultaneousCalls.ActionChainRunCount = configs.Triggers.SimultaneousCalls.MaxActionChainRunCount
	} else {
		fraudion.State.Triggers.StateDangerousDestinations.ActionChainRunCount = configs.General.DefaultActionChainRunCount
	}

	if configs.Triggers.SmallDurationCalls.MaxActionChainRunCount != 0 {
		fraudion.State.Triggers.StateSmallDurationCalls.ActionChainRunCount = configs.Triggers.SmallDurationCalls.MaxActionChainRunCount
	} else {
		fraudion.State.Triggers.StateDangerousDestinations.ActionChainRunCount = configs.General.DefaultActionChainRunCount
	}*/

	//softswitch.Init()

	//fmt.Println(softswitch.Monitored, reflect.TypeOf(softswitch.Monitored))
	//fmt.Println(softswitch.Monitored.GetCDRsSource(), reflect.TypeOf(softswitch.Monitored.GetCDRsSource()))

	// TODO: We'll call config.Validate() here in the future

	//fmt.Println(fraudion.State)

	// TODO: This will maybe be done elsewhere!
	/*var db *sql.DB
	if configs.Triggers.DangerousDestinations.Enabled == true || configs.Triggers.ExpectedDestinations.Enabled == true || configs.Triggers.SmallDurationCalls.Enabled == true {
		fraudion.LogInfo.Println("Connecting to the CDRs Database...")
		var dbstring string
		if *argCLIDBPass == "" {
			dbstring = fmt.Sprintf("root:@tcp(localhost:3306)/asteriskcdrdb?allowOldPasswords=1")
		} else {
			dbstring = fmt.Sprintf("root:%s@tcp(localhost:3306)/asteriskcdrdb?allowOldPasswords=1", *argCLIDBPass)
		}
		db, err = sql.Open("mysql", dbstring)
		if err != nil {
			fraudion.LogError.Fatalf("There was an error (%s) trying to open a connection to the database\n", err)
		}
	}

	// Launch Triggers!
	fraudion.LogInfo.Println("Launching enabled triggers...")
	if configs.Triggers.SimultaneousCalls.Enabled == true {
		go monitors.SimultaneousCallsRun()
	}
	if configs.Triggers.DangerousDestinations.Enabled == true {
		go monitors.DangerousDestinationsRun(db)
	}

	/*if configs.Triggers.ExpectedDestinations.Enabled == true {
		go triggers.ExpectedDestinationsRun(configs, db)
	}

	if configs.Triggers.SmallDurationCalls.Enabled == true {
		go triggers.SmallDurationCallsRun(configs, db)
	}*/

	log.LogS("INFO", "Main thread is going to sleep...")

	// Sleep!
	for {

		// Main "thread" has to Sleep or else 100% CPU...
		time.Sleep(100000 * time.Hour)

	}

}
