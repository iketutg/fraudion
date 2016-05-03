package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/fraudion"
	"github.com/andmar/fraudion/logger"
	"github.com/andmar/fraudion/softswitch"

	"github.com/andmar/marlog"

	_ "github.com/go-sql-driver/mysql"
)

// Defines Constants
const (
	constDefaultConfigDir = "/etc/fraudion"
	constDefaultLogFile   = "/var/log/fraudion.log" // TODO: Should we keep the system defaulting to STDOUT or use this value?
)

// Defines expected CLI flags
var (
	argCLILogTo     = flag.String("logto", "", "<help message for 'logto'>") // NOTE: The default is "" because we use this to detect if the user has specified any file, if not, the system defaults to using STDOUT automatically.
	argCLIConfigDir = flag.String("configin", constDefaultConfigDir, "<help message for 'configin'>")
	argCLIDBPass    = flag.String("dbpass", "", "<help message for 'dbpass'>")
)

// Starts here!
func main() {

	fraudion := fraudion.Global // NOTE: fraudion.Global (and it's pointers) is (are) initialized on fraudion's package init() function

	argCLILogToS := "Buh"
	argCLILogTo := &argCLILogToS

	// Logger Setup
	log := marlog.MarLog
	log.Prefix = "FRAUDION"
	log.Flags = marlog.FlagLdate | marlog.FlagLtime | marlog.FlagLlongfile

	// TODO: Error handling here?
	log.SetStamp("ERROR", "*STDOUT")
	log.SetStamp("DEBUG", "*STDOUT")
	log.SetStamp("INFO", "*STDOUT")

	fraudion.StartUpTime = time.Now()
	log.LogS("INFO", fmt.Sprintf("Starting Fraudion at %s", fraudion.StartUpTime))
	log.LogS("INFO", "Parsing CLI flags...")
	//logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf(fmt.Sprintf("Starting Fraudion at %s", fraudion.StartUpTime)), false)
	//logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("Parsing CLI flags..."), false)
	flag.Parse()

	// TODO: This should default to constDefaultLogFile, maybe even handle a flag to disable logging
	if strings.ToLower(*argCLILogTo) != "" {

		logFile, err := os.OpenFile(*argCLILogTo, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {

			log.LogS("ERROR", fmt.Sprintf("Can't start, there was a problem (%s) opening the Log file. :(", err))
			//logger.Log.Write(logger.ConstLoggerLevelError, fmt.Sprintf("Can't start, there was a problem (%s) opening the Log file. :(", err), true)

			//os.Exit(1) // TODO: Add a configuration item to tell the system to quit when logging is not possible

		} else {

			// TODO: Error handling here
			log.SetOutputHandle("MAINFILE", logFile)

			log.AddOuputHandles("INFO", "MAINFILE")
			log.AddOuputHandles("DEBUG", "MAINFILE")
			log.AddOuputHandles("ERROR", "MAINFILE")

			log.LogS("INFO", fmt.Sprintf("Started logging to \"%s\"", *argCLILogTo))
			//logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("Outputting Log to \"%s\"", *argCLILogTo), false)

			logger.Log.SetHandles(logFile, logFile, logFile, logFile) // NOTE: Overwrite the default handles on the Logger object.

			//log.LogS("INFO", "\n")
			//logger.Log.Write(logger.ConstLoggerLevelInfo, "\n", false)

		}

	}

	// TODO: The parsing/validation/loading will be moved to a special Init function at package level, as the Softswitch package does
	logger.Log.Write(logger.ConstLoggerLevelInfo, fmt.Sprintf("Starting Fraudion Log at %s", fraudion.StartUpTime), false)

	configsJSON, err := config.Parse(*argCLIConfigDir)
	if err != nil {
		logger.Log.Write(logger.ConstLoggerLevelError, fmt.Sprintf("There was an error (%s) parsing the Fraudion JSON configuration file", err), true)
	}

	configs, err := config.Load(configsJSON)
	if err != nil {
		logger.Log.Write(logger.ConstLoggerLevelError, fmt.Sprintf("There was an error (%s) validating/loading the Fraudion configuration", err), true)
	}

	if configs.Triggers.DangerousDestinations.MaxActionChainRunCount != 0 {
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
	}

	softswitch.Init()

	fmt.Println(softswitch.Monitored, reflect.TypeOf(softswitch.Monitored))
	fmt.Println(softswitch.Monitored.GetCDRsSource(), reflect.TypeOf(softswitch.Monitored.GetCDRsSource()))

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
