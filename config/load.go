package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"net/http"

	"github.com/andmar/marlog"
)

// Loaded ...
var Loaded *loadedValues

// Load ...
func Load(configOriginData string, configOrigin string) error {

	log := marlog.MarLog

	if configOrigin == ConstOriginURL {

		log.LogS("INFO", "Fetching configuration from URL for reading...")
		r, err := http.Get(configOriginData)
		if err != nil {
			log.LogS("ERROR", "Could not fetch config json from URL for validation: "+err.Error())
			return err
		}

		defer r.Body.Close()

		rb, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.LogS("ERROR", "Could not read the fetched config json from URL into buffer: "+err.Error())
			return err
		}

		reader := bytes.NewReader(rb)

		log.LogS("INFO", "Validating configuration format in file...")
		if err := ValidateFromURL(reader); err != nil {
			return err
		}

		log.LogS("INFO", "Fetched contents passed Validation.")

		log.LogS("INFO", "Parsing configuration...")
		// NOTE: This seek has to be done here so that we can receive io.Reader on the parse function instead of a bytes.Reader which is not acceptable by the JsonConfigReader lib
		reader.Seek(0, 0)
		if err := parseFromURL(reader); err != nil {
			return err
		}

		log.LogS("INFO", "Data Parsed.")

	} else {

		log.LogS("INFO", "Opening configuration file for reading...")
		configFile, err := os.Open(configOriginData)
		if err != nil {
			log.LogS("ERROR", "Could not open config file for validation: "+err.Error())
			return err
		}

		defer configFile.Close()

		log.LogS("INFO", "Validating configuration format in file...")
		if err := ValidateFromFile(configFile); err != nil {
			return err
		}

		log.LogS("INFO", "File contents passed Validation.")

		log.LogS("INFO", "Parsing configuration...")
		if err := parseFromFile(configFile); err != nil {
			return err
		}

		log.LogS("INFO", "File Parsed.")

	}

	log.LogS("INFO", "Loading configuration...")
	if err := loadFromParsed(); err != nil {
		// NOTE: Remove anything that ended up in this variable inspite of the failure in the loading
		Loaded = nil
		return err
	}

	log.LogS("INFO", "Configurations Loaded.")

	return nil

}

func loadFromParsed() error {

	Loaded = new(loadedValues)

	// * General Section
	Loaded.General.Hostname = parsed.General.Hostname

	// * Softswitch Section
	Loaded.Softswitch.Type = parsed.Softswitch.Type
	Loaded.Softswitch.Version = parsed.Softswitch.Version
	Loaded.Softswitch.CDRsSource = *parsed.Softswitch.CDRsSource

	// * Monitors
	if parsed.Monitors.SimultaneousCalls == nil {
		Loaded.Monitors.SimultaneousCalls.Enabled = false
	} else {
		Loaded.Monitors.SimultaneousCalls.Enabled = parsed.Monitors.SimultaneousCalls.Enabled
		executeInterval, err := time.ParseDuration(parsed.Monitors.SimultaneousCalls.ExecuteInterval)
		if err != nil {
			return fmt.Errorf("error converting string to time.Duration on Load, this should not happen... ever")
		}
		Loaded.Monitors.SimultaneousCalls.ExecuteInterval = executeInterval
		Loaded.Monitors.SimultaneousCalls.HitThreshold = parsed.Monitors.SimultaneousCalls.HitThreshold
		Loaded.Monitors.SimultaneousCalls.MinimumNumberLength = parsed.Monitors.SimultaneousCalls.MinimumNumberLength
		Loaded.Monitors.SimultaneousCalls.ActionChainName = parsed.Monitors.SimultaneousCalls.ActionChainName
	}

	if parsed.Monitors.DangerousDestinations == nil {
		Loaded.Monitors.DangerousDestinations.Enabled = false
	} else {
		Loaded.Monitors.DangerousDestinations.Enabled = parsed.Monitors.DangerousDestinations.Enabled
		executeInterval, err := time.ParseDuration(parsed.Monitors.DangerousDestinations.ExecuteInterval)
		if err != nil {
			return fmt.Errorf("error converting string to time.Duration on Load, this should not happen... ever")
		}
		Loaded.Monitors.DangerousDestinations.ExecuteInterval = executeInterval
		Loaded.Monitors.DangerousDestinations.HitThreshold = parsed.Monitors.DangerousDestinations.HitThreshold
		Loaded.Monitors.DangerousDestinations.MinimumNumberLength = parsed.Monitors.DangerousDestinations.MinimumNumberLength
		Loaded.Monitors.DangerousDestinations.ActionChainName = parsed.Monitors.DangerousDestinations.ActionChainName
		if considerFromLast, err := time.ParseDuration(parsed.Monitors.DangerousDestinations.ConsiderCDRsFromLast); err != nil {
			considerFromLastUInt, err := strconv.Atoi(parsed.Monitors.DangerousDestinations.ConsiderCDRsFromLast)
			if err != nil {
				return fmt.Errorf("error converting string to int on Load, this should not happen... ever")
			}
			considerFromLastDurationFromInt, err := time.ParseDuration(strconv.FormatUint(uint64(considerFromLastUInt*24), 10) + "h")
			if err != nil {
				return fmt.Errorf("error converting string to time.Duration from int on Load, this should not happen... ever")
			}
			Loaded.Monitors.DangerousDestinations.ConsiderCDRsFromLast = considerFromLastDurationFromInt
		} else {
			Loaded.Monitors.DangerousDestinations.ConsiderCDRsFromLast = considerFromLast
		}
		Loaded.Monitors.DangerousDestinations.PrefixList = parsed.Monitors.DangerousDestinations.PrefixList
		Loaded.Monitors.DangerousDestinations.MatchRegex = parsed.Monitors.DangerousDestinations.MatchRegex
		Loaded.Monitors.DangerousDestinations.IgnoreRegex = parsed.Monitors.DangerousDestinations.IgnoreRegex
	}

	if parsed.Monitors.ExpectedDestinations == nil {
		Loaded.Monitors.ExpectedDestinations.Enabled = false
	} else {
		Loaded.Monitors.ExpectedDestinations.Enabled = parsed.Monitors.ExpectedDestinations.Enabled
		executeInterval, err := time.ParseDuration(parsed.Monitors.ExpectedDestinations.ExecuteInterval)
		if err != nil {
			return fmt.Errorf("error converting string to time.Duration on Load, this should not happen... ever")
		}
		Loaded.Monitors.ExpectedDestinations.ExecuteInterval = executeInterval
		Loaded.Monitors.ExpectedDestinations.HitThreshold = parsed.Monitors.ExpectedDestinations.HitThreshold
		Loaded.Monitors.ExpectedDestinations.MinimumNumberLength = parsed.Monitors.ExpectedDestinations.MinimumNumberLength
		Loaded.Monitors.ExpectedDestinations.ActionChainName = parsed.Monitors.ExpectedDestinations.ActionChainName
		if considerFromLast, err := time.ParseDuration(parsed.Monitors.ExpectedDestinations.ConsiderCDRsFromLast); err != nil {
			considerFromLastUInt, err := strconv.Atoi(parsed.Monitors.ExpectedDestinations.ConsiderCDRsFromLast)
			if err != nil {
				return fmt.Errorf("error converting string to int on Load, this should not happen... ever")
			}
			considerFromLastDurationFromInt, err := time.ParseDuration(strconv.FormatUint(uint64(considerFromLastUInt*24), 10) + "h")
			if err != nil {
				return fmt.Errorf("error converting string to time.Duration from int on Load, this should not happen... ever")
			}
			Loaded.Monitors.ExpectedDestinations.ConsiderCDRsFromLast = considerFromLastDurationFromInt
		} else {
			Loaded.Monitors.ExpectedDestinations.ConsiderCDRsFromLast = considerFromLast
		}
		Loaded.Monitors.ExpectedDestinations.PrefixList = parsed.Monitors.ExpectedDestinations.PrefixList
		Loaded.Monitors.ExpectedDestinations.MatchRegex = parsed.Monitors.ExpectedDestinations.MatchRegex
		Loaded.Monitors.ExpectedDestinations.IgnoreRegex = parsed.Monitors.ExpectedDestinations.IgnoreRegex
	}

	if parsed.Monitors.ExpectedDestinations == nil {
		Loaded.Monitors.ExpectedDestinations.Enabled = false
	} else {
		Loaded.Monitors.ExpectedDestinations.Enabled = parsed.Monitors.ExpectedDestinations.Enabled
		executeInterval, err := time.ParseDuration(parsed.Monitors.ExpectedDestinations.ExecuteInterval)
		if err != nil {
			return fmt.Errorf("error converting string to time.Duration on Load, this should not happen... ever")
		}
		Loaded.Monitors.ExpectedDestinations.ExecuteInterval = executeInterval
		Loaded.Monitors.ExpectedDestinations.HitThreshold = parsed.Monitors.ExpectedDestinations.HitThreshold
		Loaded.Monitors.ExpectedDestinations.MinimumNumberLength = parsed.Monitors.ExpectedDestinations.MinimumNumberLength
		Loaded.Monitors.ExpectedDestinations.ActionChainName = parsed.Monitors.ExpectedDestinations.ActionChainName
		if considerFromLast, err := time.ParseDuration(parsed.Monitors.ExpectedDestinations.ConsiderCDRsFromLast); err != nil {
			considerFromLastUInt, err := strconv.Atoi(parsed.Monitors.ExpectedDestinations.ConsiderCDRsFromLast)
			if err != nil {
				return fmt.Errorf("error converting string to int on Load, this should not happen... ever")
			}
			considerFromLastDurationFromInt, err := time.ParseDuration(strconv.FormatUint(uint64(considerFromLastUInt*24), 10) + "h")
			if err != nil {
				return fmt.Errorf("error converting string to time.Duration from int on Load, this should not happen... ever")
			}
			Loaded.Monitors.ExpectedDestinations.ConsiderCDRsFromLast = considerFromLastDurationFromInt
		} else {
			Loaded.Monitors.ExpectedDestinations.ConsiderCDRsFromLast = considerFromLast
		}
		Loaded.Monitors.ExpectedDestinations.PrefixList = parsed.Monitors.ExpectedDestinations.PrefixList
		Loaded.Monitors.ExpectedDestinations.MatchRegex = parsed.Monitors.ExpectedDestinations.MatchRegex
		Loaded.Monitors.ExpectedDestinations.IgnoreRegex = parsed.Monitors.ExpectedDestinations.IgnoreRegex
	}

	if parsed.Monitors.SmallDurationCalls == nil {
		Loaded.Monitors.SmallDurationCalls.Enabled = false
	} else {
		Loaded.Monitors.SmallDurationCalls.Enabled = parsed.Monitors.SmallDurationCalls.Enabled
		executeInterval, err := time.ParseDuration(parsed.Monitors.SmallDurationCalls.ExecuteInterval)
		if err != nil {
			return fmt.Errorf("error converting string to time.Duration on Load, this should not happen... ever")
		}
		Loaded.Monitors.SmallDurationCalls.ExecuteInterval = executeInterval
		Loaded.Monitors.SmallDurationCalls.HitThreshold = parsed.Monitors.SmallDurationCalls.HitThreshold
		Loaded.Monitors.SmallDurationCalls.MinimumNumberLength = parsed.Monitors.SmallDurationCalls.MinimumNumberLength
		Loaded.Monitors.SmallDurationCalls.ActionChainName = parsed.Monitors.SmallDurationCalls.ActionChainName
		if considerFromLast, err := time.ParseDuration(parsed.Monitors.SmallDurationCalls.ConsiderCDRsFromLast); err != nil {
			considerFromLastUInt, err := strconv.Atoi(parsed.Monitors.SmallDurationCalls.ConsiderCDRsFromLast)
			if err != nil {
				return fmt.Errorf("error converting string to int on Load, this should not happen... ever")
			}
			considerFromLastDurationFromInt, err := time.ParseDuration(strconv.FormatUint(uint64(considerFromLastUInt*24), 10) + "h")
			if err != nil {
				return fmt.Errorf("error converting string to time.Duration from int on Load, this should not happen... ever")
			}
			Loaded.Monitors.SmallDurationCalls.ConsiderCDRsFromLast = considerFromLastDurationFromInt
		} else {
			Loaded.Monitors.SmallDurationCalls.ConsiderCDRsFromLast = considerFromLast
		}
		durationThreshold, err := time.ParseDuration(parsed.Monitors.SmallDurationCalls.DurationThreshold)
		if err != nil {
			return fmt.Errorf("error converting string to time.Duration on Load, this should not happen... ever")
		}
		Loaded.Monitors.SmallDurationCalls.DurationThreshold = durationThreshold
	}

	// * Actions
	if parsed.Actions.Email == nil {
		Loaded.Actions.Email.Enabled = false
	} else {
		Loaded.Actions.Email.Enabled = parsed.Actions.Email.Enabled
		Loaded.Actions.Email.Recurrent = parsed.Actions.Email.Recurrent
		Loaded.Actions.Email.Username = parsed.Actions.Email.Username
		Loaded.Actions.Email.Password = parsed.Actions.Email.Password
		Loaded.Actions.Email.Title = parsed.Actions.Email.Title
		Loaded.Actions.Email.Body = parsed.Actions.Email.Body
	}

	if parsed.Actions.LocalCommands == nil {
		Loaded.Actions.LocalCommands.Enabled = false
	} else {
		Loaded.Actions.LocalCommands.Enabled = parsed.Actions.LocalCommands.Enabled
		Loaded.Actions.LocalCommands.Recurrent = parsed.Actions.LocalCommands.Recurrent
	}

	// * Action Chains
	// NOTE: Action Chains configured for the Monitors exist?
	if Loaded.Monitors.DangerousDestinations.Enabled == true {
		existsForDangerousDestinations := false
		for chainName := range *parsed.ActionChains {
			if Loaded.Monitors.DangerousDestinations.ActionChainName == chainName {
				existsForDangerousDestinations = true
			}
		}
		if existsForDangerousDestinations == false {
			return fmt.Errorf("action chain for Dangerous Destinations not enabled")
		}
	}

	if Loaded.Monitors.ExpectedDestinations.Enabled == true {
		existsForExpectedDestinations := false
		for chainName := range *parsed.ActionChains {
			if Loaded.Monitors.ExpectedDestinations.ActionChainName == chainName {
				existsForExpectedDestinations = true
			}
		}
		if existsForExpectedDestinations == false {
			return fmt.Errorf("action chain for Expected Destinations not enabled")
		}
	}

	if Loaded.Monitors.SmallDurationCalls.Enabled == true {
		existsForSmallDurationCalls := false
		for chainName := range *parsed.ActionChains {
			if Loaded.Monitors.SmallDurationCalls.ActionChainName == chainName {
				existsForSmallDurationCalls = true
			}
		}
		if existsForSmallDurationCalls == false {
			return fmt.Errorf("action chain for Small Duration Calls not enabled")
		}
	}

	if Loaded.Monitors.SimultaneousCalls.Enabled == true {
		existsForSimultaneousCalls := false
		for chainName := range *parsed.ActionChains {
			if Loaded.Monitors.SimultaneousCalls.ActionChainName == chainName {
				existsForSimultaneousCalls = true
			}
		}
		if existsForSimultaneousCalls == false {
			return fmt.Errorf("action chain for Simultaneous Calls not enabled")
		}
	}

	// * Action Chains
	Loaded.ActionChains = *parsed.ActionChains
	// NOTE: All Actions in Chains are enabled?
	allEnabled := true
	for _, chain := range *parsed.ActionChains {

		for _, v := range chain {

			switch v.ActionName {
			case "*email":
				if Loaded.Actions.Email.Enabled == false {
					allEnabled = false
				}
			case "*local_commands":
				if Loaded.Actions.LocalCommands.Enabled == false {
					allEnabled = false
				}
			default:
				return fmt.Errorf("unknown action in action chain config, this should not happen... ever")
			}

		}

	}
	if allEnabled == false {
		return fmt.Errorf("some configured action in chain is not enabled")
	}

	// * Data Groups
	Loaded.DataGroups = *parsed.DataGroups
	// NOTE: All DataGroups used in Chains have the information for the specified Action?
	for _, chain := range *parsed.ActionChains {

		for _, action := range chain {

			for _, dataGroupName := range action.DataGroupNames {

				dataGroup, exists := Loaded.DataGroups[dataGroupName]
				if exists == false {
					return fmt.Errorf("some configured data group in chain is not defined")
				}

				switch action.ActionName {
				case "*email":
					if dataGroup.EmailAddress == "" {
						return fmt.Errorf("some E-mail action requires email_address value defined in Data Group")
					}
				case "*local_commands":
					if dataGroup.CommandName == "" {
						return fmt.Errorf("some Local Commands action requires command_name value defined in Data Group")
					}
				default:
					return fmt.Errorf("unknown action, this should never happen")
				}

			}

		}

	}

	return nil

}

type loadedValues struct {
	General      general
	Softswitch   softswitch
	Monitors     monitors
	Actions      actions
	ActionChains actionChains
	DataGroups   dataGroups
}

type general struct {
	Hostname string
}

type softswitch struct {
	Type       string
	Version    string
	CDRsSource cdrsSource
}

type cdrsSource map[string]string

type monitors struct {
	SimultaneousCalls     MonitorSimultaneousCalls
	DangerousDestinations MonitorDangerousDestinations
	ExpectedDestinations  MonitorExpectedDestinations
	SmallDurationCalls    MonitorSmallCallDurations
}

type monitorBase struct {
	Enabled             bool
	ExecuteInterval     time.Duration
	HitThreshold        uint32
	MinimumNumberLength uint32
	ActionChainName     string
}

// MonitorSimultaneousCalls ...
type MonitorSimultaneousCalls struct {
	monitorBase
}

// MonitorDangerousDestinations ...
type MonitorDangerousDestinations struct {
	monitorBase
	ConsiderCDRsFromLast time.Duration
	PrefixList           []string
	MatchRegex           string
	IgnoreRegex          string
}

// MonitorExpectedDestinations ...
type MonitorExpectedDestinations struct {
	monitorBase
	ConsiderCDRsFromLast time.Duration
	PrefixList           []string
	MatchRegex           string
	IgnoreRegex          string
}

// MonitorSmallCallDurations ...
type MonitorSmallCallDurations struct {
	monitorBase
	ConsiderCDRsFromLast time.Duration
	DurationThreshold    time.Duration
}

type actions struct {
	Email         actionEmail
	LocalCommands actionLocalCommands
}

type actionBase struct {
	Enabled   bool
	Recurrent bool
}

type actionEmail struct {
	actionBase
	Type     string
	Username string
	Password string
	Title    string
	Body     string
}

type actionLocalCommands struct {
	actionBase
}

type actionChains map[string][]actionChainAction

type actionChainAction struct {
	ActionName     string   `json:"action_name"`
	DataGroupNames []string `json:"data_groups"`
}

type dataGroups map[string]dataGroup

type dataGroup struct {
	PhoneNumber      string            `json:"phone_number"`
	EmailAddress     string            `json:"email_address"`
	HTTPURL          string            `json:"http_url"`
	HTTPMethod       string            `json:"http_method"`
	HTTPParameters   map[string]string `json:"http_parameters"`
	CommandName      string            `json:"command_name"`
	CommandArguments string            `json:"command_arguments"`
}
