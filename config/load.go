package config

import (
	"fmt"
	"strconv"
	"time"
)

const (
	defaultExecuteInterval          = "5m"
	defaultHitThreshold             = 5
	defaultMinimumNumberLength      = 5
	defaultActionChainMame          = "*default"
	defaultActionChainHoldoffPeriod = 0
	defaultActionChainRunCount      = 0
)

// Loaded After config.Load(...) is called, this variable holds the final configuration values in it's appropriate types
var Loaded *loadedValues

// Load SetsUp the loading of the configuration from specified file and handle outputs, keeps config.Loaded nil in case of some error happen
func Load(configDir string, configFileName string, validateOnly bool) error {

	if err := doLoad(configDir, configFileName, validateOnly); err != nil {
		Loaded = nil
		return err
	}

	return nil

}

// doLoad Loads configuration from specified file
func doLoad(configDir string, configFileName string, validateOnly bool) error {

	if err := parse(configDir, configFileName); err != nil {
		return err
	}

	if hasErrors, errors := validate(); hasErrors == true {
		return fmt.Errorf("Failed Validation. Errors: %s", errors) // TODO: Do something with the list of errors returned...
	}

	Loaded = new(loadedValues)

	// General Section

	// Softswitch Section
	Loaded.Softswitch.System = parsed.Softswitch.System
	Loaded.Softswitch.Version = parsed.Softswitch.Version
	Loaded.Softswitch.CDRsSource = parsed.Softswitch.CDRsSource

	// CDRs Sources
	Loaded.CDRsSources = *parsed.CDRsSources

	// Monitors
	Loaded.Monitors.SimultaneousCalls.Enabled = parsed.Monitors.SimultaneousCalls.Enabled
	executeInterval, err := time.ParseDuration(parsed.Monitors.SimultaneousCalls.ExecuteInterval)
	if err != nil {
		return fmt.Errorf("error parsing duration for \"execute_interval\" in \"monitors/simultaneous_calls\"")
	}
	Loaded.Monitors.SimultaneousCalls.ExecuteInterval = executeInterval
	Loaded.Monitors.SimultaneousCalls.HitThreshold = parsed.Monitors.SimultaneousCalls.HitThreshold
	Loaded.Monitors.SimultaneousCalls.MinimumNumberLength = parsed.Monitors.SimultaneousCalls.MinimumNumberLength
	Loaded.Monitors.SimultaneousCalls.ActionChainName = parsed.Monitors.SimultaneousCalls.ActionChainName
	Loaded.Monitors.SimultaneousCalls.ActionChainHoldoffPeriod = parsed.Monitors.SimultaneousCalls.ActionChainHoldoffPeriod
	Loaded.Monitors.SimultaneousCalls.MaxActionChainRunCount = parsed.Monitors.SimultaneousCalls.MaxActionChainRunCount

	Loaded.Monitors.DangerousDestinations.Enabled = parsed.Monitors.DangerousDestinations.Enabled
	executeInterval, err = time.ParseDuration(parsed.Monitors.DangerousDestinations.ExecuteInterval)
	if err != nil {
		return fmt.Errorf("error parsing duration for \"execute_interval\" in \"monitors/dangerous_destinations\"")
	}
	Loaded.Monitors.DangerousDestinations.ExecuteInterval = executeInterval
	Loaded.Monitors.DangerousDestinations.HitThreshold = parsed.Monitors.DangerousDestinations.HitThreshold
	Loaded.Monitors.DangerousDestinations.MinimumNumberLength = parsed.Monitors.DangerousDestinations.MinimumNumberLength
	Loaded.Monitors.DangerousDestinations.ActionChainName = parsed.Monitors.DangerousDestinations.ActionChainName
	Loaded.Monitors.DangerousDestinations.ActionChainHoldoffPeriod = parsed.Monitors.DangerousDestinations.ActionChainHoldoffPeriod
	Loaded.Monitors.DangerousDestinations.MaxActionChainRunCount = parsed.Monitors.DangerousDestinations.MaxActionChainRunCount
	if considerFromLast, err := time.ParseDuration(parsed.Monitors.DangerousDestinations.ConsiderCDRsFromLast); err != nil {
		considerFromLastUInt, err := strconv.Atoi(parsed.Monitors.DangerousDestinations.ConsiderCDRsFromLast)
		if err != nil {
			return fmt.Errorf("error converting value of \"consider_cdrs_from_last\" to int as a fallback from parseable time.Duration in \"monitors/dangerous_destinations\"")
		}
		considerFromLastDuration, err := time.ParseDuration(strconv.FormatUint(uint64(considerFromLastUInt*24), 10) + "h")
		if err != nil {
			return fmt.Errorf("error creating duration for \"consider_cdrs_from_last\" in \"monitors/dangerous_destinations\"")
		}
		Loaded.Monitors.DangerousDestinations.ConsiderCDRsFromLast = considerFromLastDuration
	} else {
		Loaded.Monitors.DangerousDestinations.ConsiderCDRsFromLast = considerFromLast
	}
	Loaded.Monitors.DangerousDestinations.PrefixList = parsed.Monitors.DangerousDestinations.PrefixList
	Loaded.Monitors.DangerousDestinations.MatchRegex = parsed.Monitors.DangerousDestinations.MatchRegex
	Loaded.Monitors.DangerousDestinations.IgnoreRegex = parsed.Monitors.DangerousDestinations.IgnoreRegex

	Loaded.Monitors.ExpectedDestinations.Enabled = parsed.Monitors.ExpectedDestinations.Enabled
	executeInterval, err = time.ParseDuration(parsed.Monitors.ExpectedDestinations.ExecuteInterval)
	if err != nil {
		return fmt.Errorf("error parsing duration for \"execute_interval\" in \"monitors/expected_destinations\"")
	}
	Loaded.Monitors.ExpectedDestinations.ExecuteInterval = executeInterval
	Loaded.Monitors.ExpectedDestinations.HitThreshold = parsed.Monitors.ExpectedDestinations.HitThreshold
	Loaded.Monitors.ExpectedDestinations.MinimumNumberLength = parsed.Monitors.ExpectedDestinations.MinimumNumberLength
	Loaded.Monitors.ExpectedDestinations.ActionChainName = parsed.Monitors.ExpectedDestinations.ActionChainName
	Loaded.Monitors.ExpectedDestinations.ActionChainHoldoffPeriod = parsed.Monitors.ExpectedDestinations.ActionChainHoldoffPeriod
	Loaded.Monitors.ExpectedDestinations.MaxActionChainRunCount = parsed.Monitors.ExpectedDestinations.MaxActionChainRunCount
	if considerFromLast, err := time.ParseDuration(parsed.Monitors.ExpectedDestinations.ConsiderCDRsFromLast); err != nil {
		considerFromLastUInt, err := strconv.Atoi(parsed.Monitors.ExpectedDestinations.ConsiderCDRsFromLast)
		if err != nil {
			return fmt.Errorf("error converting value of \"consider_cdrs_from_last\" to int as a fallback from parseable time.Duration in \"monitors/expected_destinations\"")
		}
		considerFromLastDuration, err := time.ParseDuration(strconv.FormatUint(uint64(considerFromLastUInt*24), 10) + "h")
		if err != nil {
			return fmt.Errorf("error creating duration for \"consider_cdrs_from_last\" in \"monitors/expected_destinations\"")
		}
		Loaded.Monitors.ExpectedDestinations.ConsiderCDRsFromLast = considerFromLastDuration
	} else {
		Loaded.Monitors.ExpectedDestinations.ConsiderCDRsFromLast = considerFromLast
	}
	Loaded.Monitors.ExpectedDestinations.PrefixList = parsed.Monitors.ExpectedDestinations.PrefixList
	Loaded.Monitors.ExpectedDestinations.MatchRegex = parsed.Monitors.ExpectedDestinations.MatchRegex
	Loaded.Monitors.ExpectedDestinations.IgnoreRegex = parsed.Monitors.ExpectedDestinations.IgnoreRegex

	Loaded.Monitors.SmallDurationCalls.Enabled = parsed.Monitors.SmallDurationCalls.Enabled
	executeInterval, err = time.ParseDuration(parsed.Monitors.SmallDurationCalls.ExecuteInterval)
	if err != nil {
		return fmt.Errorf("error parsing duration for \"execute_interval\" in \"monitors/small_duration_calls\"")
	}
	Loaded.Monitors.SmallDurationCalls.ExecuteInterval = executeInterval
	Loaded.Monitors.SmallDurationCalls.HitThreshold = parsed.Monitors.SmallDurationCalls.HitThreshold
	Loaded.Monitors.SmallDurationCalls.MinimumNumberLength = parsed.Monitors.SmallDurationCalls.MinimumNumberLength
	Loaded.Monitors.SmallDurationCalls.ActionChainName = parsed.Monitors.SmallDurationCalls.ActionChainName
	Loaded.Monitors.SmallDurationCalls.ActionChainHoldoffPeriod = parsed.Monitors.SmallDurationCalls.ActionChainHoldoffPeriod
	Loaded.Monitors.SmallDurationCalls.MaxActionChainRunCount = parsed.Monitors.SmallDurationCalls.MaxActionChainRunCount
	if considerFromLast, err := time.ParseDuration(parsed.Monitors.SmallDurationCalls.ConsiderCDRsFromLast); err != nil {
		considerFromLastUInt, err := strconv.Atoi(parsed.Monitors.SmallDurationCalls.ConsiderCDRsFromLast)
		if err != nil {
			return fmt.Errorf("error converting value of \"consider_cdrs_from_last\" to int as a fallback from parseable time.Duration in \"monitors/small_duration_calls\"")
		}
		considerFromLastDuration, err := time.ParseDuration(strconv.FormatUint(uint64(considerFromLastUInt*24), 10) + "h")
		if err != nil {
			return fmt.Errorf("error creating duration for \"consider_cdrs_from_last\" in \"monitors/small_duration_calls\"")
		}
		Loaded.Monitors.SmallDurationCalls.ConsiderCDRsFromLast = considerFromLastDuration
	} else {
		Loaded.Monitors.SmallDurationCalls.ConsiderCDRsFromLast = considerFromLast
	}
	durationThreshold, err := time.ParseDuration(parsed.Monitors.SmallDurationCalls.DurationThreshold)
	if err != nil {
		return fmt.Errorf("error parsing duration for \"duration_threshold\" in \"monitors/small_duration_calls\"")
	}
	Loaded.Monitors.SmallDurationCalls.DurationThreshold = durationThreshold

	// Actions
	// fmt.Println(parsed.Actions.Email)
	// fmt.Println(*parsed.Actions.Email)
	if parsed.Actions.Email != nil {
		Loaded.Actions.Email = new(actionEmail)
		Loaded.Actions.Email.Enabled = parsed.Actions.Email.Enabled
		Loaded.Actions.Email.Username = parsed.Actions.Email.Username
		Loaded.Actions.Email.Password = parsed.Actions.Email.Password
		Loaded.Actions.Email.Message = parsed.Actions.Email.Message
	}
	if parsed.Actions.Call != nil {
		Loaded.Actions.Call = new(actionCall)
		Loaded.Actions.Call.Enabled = parsed.Actions.Call.Enabled
	}
	if parsed.Actions.HTTP != nil {
		Loaded.Actions.HTTP = new(actionHTTP)
		Loaded.Actions.HTTP.Enabled = parsed.Actions.HTTP.Enabled
	}
	if parsed.Actions.LocalCommands != nil {
		Loaded.Actions.LocalCommands = new(actionLocalCommands)
		Loaded.Actions.LocalCommands.Enabled = parsed.Actions.LocalCommands.Enabled
	}

	// Action Chains
	Loaded.ActionChains = make(map[string][]actionChainAction)
	for k, v := range *parsed.ActionChains {
		Loaded.ActionChains[k] = v
	}

	// Data Groups
	Loaded.DataGroups = make(map[string]dataGroup)
	for k, v := range *parsed.DataGroups {
		Loaded.DataGroups[k] = v
	}

	return nil

}

type loadedValues struct {
	General      general
	Softswitch   softswitchInfo
	CDRsSources  cdrsSources
	Monitors     monitors
	Actions      actions
	ActionChains actionChains
	DataGroups   dataGroups
}

type general struct{}

type softswitchInfo struct {
	System     string
	Version    string
	CDRsSource string
}

// Monitors ...
type monitors struct {
	SimultaneousCalls     monitorSimultaneousCalls
	DangerousDestinations monitorDangerousDestinations
	ExpectedDestinations  monitorExpectedDestinations
	SmallDurationCalls    monitorSmallCallDurations
}

type monitorBase struct {
	Enabled                  bool
	ExecuteInterval          time.Duration
	HitThreshold             uint32
	MinimumNumberLength      uint32
	ActionChainName          string
	ActionChainHoldoffPeriod uint32
	MaxActionChainRunCount   uint32
}

type monitorSimultaneousCalls struct {
	monitorBase
}

type monitorDangerousDestinations struct {
	monitorBase
	ConsiderCDRsFromLast time.Duration
	PrefixList           []string
	MatchRegex           string
	IgnoreRegex          string
}

type monitorExpectedDestinations struct {
	monitorBase
	ConsiderCDRsFromLast time.Duration
	PrefixList           []string
	MatchRegex           string
	IgnoreRegex          string
}

type monitorSmallCallDurations struct {
	monitorBase
	ConsiderCDRsFromLast time.Duration
	DurationThreshold    time.Duration
}

type actions struct {
	Email         *actionEmail
	Call          *actionCall
	HTTP          *actionHTTP
	LocalCommands *actionLocalCommands
}

type actionEmail struct {
	Enabled  bool
	Username string
	Password string
	Message  string
}

type actionCall struct {
	Enabled bool
}

type actionHTTP struct {
	Enabled bool
}

type actionLocalCommands struct {
	Enabled bool
}
