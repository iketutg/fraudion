package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/andmar/marlog"
)

const (
	defaultExecuteInterval     = "5m"
	defaultHitThreshold        = 5
	defaultMinimumNumberLength = 5
	defaultDataGroupName       = "*default"
	defaultActionChainName     = "*default"
)

var Loaded *loadedValues

func Load(configFileFullName string, validateOnly bool) error {

	log := marlog.MarLog

	log.LogS("INFO", "Preparing configuration loading...")

	configFile, err := os.Open(configFileFullName)
	if err != nil {
		return err
	}
	defer configFile.Close()

	if err := validateFromFile(configFile); err != nil {
		return err
	}

	log.LogS("INFO", "File contents passed Validation.")

	if validateOnly == true {
		return nil
	}

	if err := parseFromFile(configFile); err != nil {
		parsed = nil
		return err
	}

	log.LogS("INFO", "File Parsed.")

	if err := loadFromParsed(); err != nil {
		Loaded = nil
		return err
	}

	log.LogS("INFO", "Configurations Loaded.")

	return nil

}

func loadFromParsed() error {

	Loaded = new(loadedValues)

	// General Section

	// Softswitch Section
	Loaded.Softswitch.Type = parsed.Softswitch.Type
	Loaded.Softswitch.Version = parsed.Softswitch.Version
	Loaded.Softswitch.CDRsSourceName = parsed.Softswitch.CDRsSourceName

	// CDRs Sources
	// NOTE: Source configured for the Softswitch exists?
	exists := false
	for sourceName := range *parsed.CDRsSources {
		if Loaded.Softswitch.CDRsSourceName == sourceName {
			exists = true
			break
		}
	}

	if exists == false {
		return fmt.Errorf("CDR Source Information for Softswitch is not configured")
	}

	Loaded.CDRsSources = *parsed.CDRsSources

	// Monitors
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

	// Actions
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

	// Action Chains
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

	Loaded.ActionChains = *parsed.ActionChains

	// Data Groups
	// TODO: All DataGroups used in Chains have the information for the specified Action?
	Loaded.DataGroups = *parsed.DataGroups

	fmt.Println("\nParsed Configurations:")
	fmt.Println(parsed)
	fmt.Println(parsed.General)
	fmt.Println(parsed.Softswitch)
	fmt.Println(parsed.CDRsSources)
	fmt.Println(parsed.Monitors.DangerousDestinations, parsed.Monitors.ExpectedDestinations, parsed.Monitors.SimultaneousCalls, parsed.Monitors.SmallDurationCalls)
	fmt.Println(parsed.Actions.Email, parsed.Actions.LocalCommands)
	fmt.Println(parsed.ActionChains)
	fmt.Println(parsed.DataGroups)
	fmt.Println()

	fmt.Println("*\n\n\n*")

	fmt.Println("\nParsed Configurations:")
	fmt.Println(Loaded)
	fmt.Println()

	return nil

}

/*

 */

type loadedValues struct {
	General      general
	Softswitch   softswitch
	CDRsSources  cdrsSources
	Monitors     monitors
	Actions      actions
	ActionChains actionChains
	DataGroups   dataGroups
}

type general struct{}

type softswitch struct {
	Type           string
	Version        string
	CDRsSourceName string
}

type cdrsSources map[string]map[string]string

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

type MonitorSimultaneousCalls struct {
	monitorBase
}

type MonitorDangerousDestinations struct {
	monitorBase
	ConsiderCDRsFromLast time.Duration
	PrefixList           []string
	MatchRegex           string
	IgnoreRegex          string
}

type MonitorExpectedDestinations struct {
	monitorBase
	ConsiderCDRsFromLast time.Duration
	PrefixList           []string
	MatchRegex           string
	IgnoreRegex          string
}

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

type dataGroups map[string]dataGroup

type actionChainAction struct {
	ActionName     string   `json:"action_name"`
	DataGroupNames []string `json:"data_groups"`
}

type dataGroup struct {
	PhoneNumber      string            `json:"phone_number"`
	EmailAddress     string            `json:"data_groups"`
	HTTPURL          string            `json:"http_url"`
	HTTPMethod       string            `json:"http_method"`
	HTTPParameters   map[string]string `json:"data_groups"`
	CommandName      string            `json:"command_name"`
	CommandArguments string            `json:"command_arguments"`
}
