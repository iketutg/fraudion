package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"path/filepath"

	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/softswitch"
	"github.com/andmar/fraudion/state"

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

	// Logger Setup
	log := marlog.MarLog
	log.Prefix = "FRAUDION"
	log.Flags = marlog.FlagLdate | marlog.FlagLtime | marlog.FlagLlongfile

	// TODO: Error handling here?
	log.SetStamp("ERROR", "*STDOUT")
	log.SetStamp("DEBUG", "*STDOUT")
	log.SetStamp("INFO", "*STDOUT")

	state.System.StartUpTime = time.Now()

	// TODO: Error handling here?
	log.LogS("INFO", fmt.Sprintf("Starting Fraudion at %s", state.System.StartUpTime))
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

		}

	}

	if err := config.Load(*argCLIConfigIn, *argCLIConfigFilename, false); err != nil {
		log.LogO("ERROR", err.Error(), marlog.OptionFatal) // TODO: This has to be changed becase config.Validate() returns an array/slice of errors
	}

	fmt.Println("\nLoaded Configurations:")
	fmt.Println(config.Loaded.General)
	fmt.Println(config.Loaded.Softswitch)
	fmt.Println(config.Loaded.CDRsSources)
	fmt.Println(config.Loaded.Monitors.DangerousDestinations, config.Loaded.Monitors.ExpectedDestinations, config.Loaded.Monitors.SimultaneousCalls, config.Loaded.Monitors.SmallDurationCalls)
	fmt.Println(config.Loaded.Actions.Email, config.Loaded.Actions.Call, config.Loaded.Actions.HTTP, config.Loaded.Actions.LocalCommands)
	fmt.Println(config.Loaded.ActionChains)
	fmt.Println(config.Loaded.DataGroups)
	fmt.Println()

	fmt.Println("Softswitch:")
	fmt.Println(softswitch.Monitored)
	fmt.Println()

	fmt.Println("CDRs Sources:")
	fmt.Println(config.Loaded.CDRsSources)
	fmt.Println()

	log.LogS("INFO", "All Ok. Main thread is going to sleep now...")

	// Sleep!
	for {

		// Main "thread" has to Sleep or else 100% CPU...
		time.Sleep(100000 * time.Hour)

	}

}
