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

// Load Loads configuration from specified file
func Load() error {

	Loaded = new(loadedValues)

	// General Section

	// Softswitch Section
	Loaded.Softswitch.Brand = Parsed.Softswitch.Brand
	Loaded.Softswitch.Version = Parsed.Softswitch.Version
	Loaded.Softswitch.CDRsSource = Parsed.Softswitch.CDRsSource

	// CDRs Sources
	Loaded.CDRsSources = make(map[string]map[string]string)
	for k, v := range *Parsed.CDRsSources {
		Loaded.CDRsSources[k] = v
	}

	// Monitors
	Loaded.Monitors.SimultaneousCalls.Enabled = Parsed.Monitors.SimultaneousCalls.Enabled
	executeInterval, err := time.ParseDuration(Parsed.Monitors.SimultaneousCalls.ExecuteInterval)
	if err != nil {
		return fmt.Errorf("error parsing duration for \"execute_interval\" in \"monitors/simultaneous_calls\"")
	}
	Loaded.Monitors.SimultaneousCalls.ExecuteInterval = executeInterval
	Loaded.Monitors.SimultaneousCalls.HitThreshold = Parsed.Monitors.SimultaneousCalls.HitThreshold
	Loaded.Monitors.SimultaneousCalls.MinimumNumberLength = Parsed.Monitors.SimultaneousCalls.MinimumNumberLength
	Loaded.Monitors.SimultaneousCalls.ActionChainName = Parsed.Monitors.SimultaneousCalls.ActionChainName
	Loaded.Monitors.SimultaneousCalls.ActionChainHoldoffPeriod = Parsed.Monitors.SimultaneousCalls.ActionChainHoldoffPeriod
	Loaded.Monitors.SimultaneousCalls.MaxActionChainRunCount = Parsed.Monitors.SimultaneousCalls.MaxActionChainRunCount

	Loaded.Monitors.DangerousDestinations.Enabled = Parsed.Monitors.DangerousDestinations.Enabled
	executeInterval, err = time.ParseDuration(Parsed.Monitors.DangerousDestinations.ExecuteInterval)
	if err != nil {
		return fmt.Errorf("error parsing duration for \"execute_interval\" in \"monitors/dangerous_destinations\"")
	}
	Loaded.Monitors.DangerousDestinations.ExecuteInterval = executeInterval
	Loaded.Monitors.DangerousDestinations.HitThreshold = Parsed.Monitors.DangerousDestinations.HitThreshold
	Loaded.Monitors.DangerousDestinations.MinimumNumberLength = Parsed.Monitors.DangerousDestinations.MinimumNumberLength
	Loaded.Monitors.DangerousDestinations.ActionChainName = Parsed.Monitors.DangerousDestinations.ActionChainName
	Loaded.Monitors.DangerousDestinations.ActionChainHoldoffPeriod = Parsed.Monitors.DangerousDestinations.ActionChainHoldoffPeriod
	Loaded.Monitors.DangerousDestinations.MaxActionChainRunCount = Parsed.Monitors.DangerousDestinations.MaxActionChainRunCount
	if considerFromLast, err := time.ParseDuration(Parsed.Monitors.DangerousDestinations.ConsiderCDRsFromLast); err != nil {
		considerFromLastUInt, err := strconv.Atoi(Parsed.Monitors.DangerousDestinations.ConsiderCDRsFromLast)
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
	Loaded.Monitors.DangerousDestinations.PrefixList = Parsed.Monitors.DangerousDestinations.PrefixList
	Loaded.Monitors.DangerousDestinations.MatchRegex = Parsed.Monitors.DangerousDestinations.MatchRegex
	Loaded.Monitors.DangerousDestinations.IgnoreRegex = Parsed.Monitors.DangerousDestinations.IgnoreRegex

	Loaded.Monitors.ExpectedDestinations.Enabled = Parsed.Monitors.ExpectedDestinations.Enabled
	executeInterval, err = time.ParseDuration(Parsed.Monitors.ExpectedDestinations.ExecuteInterval)
	if err != nil {
		return fmt.Errorf("error parsing duration for \"execute_interval\" in \"monitors/expected_destinations\"")
	}
	Loaded.Monitors.ExpectedDestinations.ExecuteInterval = executeInterval
	Loaded.Monitors.ExpectedDestinations.HitThreshold = Parsed.Monitors.ExpectedDestinations.HitThreshold
	Loaded.Monitors.ExpectedDestinations.MinimumNumberLength = Parsed.Monitors.ExpectedDestinations.MinimumNumberLength
	Loaded.Monitors.ExpectedDestinations.ActionChainName = Parsed.Monitors.ExpectedDestinations.ActionChainName
	Loaded.Monitors.ExpectedDestinations.ActionChainHoldoffPeriod = Parsed.Monitors.ExpectedDestinations.ActionChainHoldoffPeriod
	Loaded.Monitors.ExpectedDestinations.MaxActionChainRunCount = Parsed.Monitors.ExpectedDestinations.MaxActionChainRunCount
	if considerFromLast, err := time.ParseDuration(Parsed.Monitors.ExpectedDestinations.ConsiderCDRsFromLast); err != nil {
		considerFromLastUInt, err := strconv.Atoi(Parsed.Monitors.ExpectedDestinations.ConsiderCDRsFromLast)
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
	Loaded.Monitors.ExpectedDestinations.PrefixList = Parsed.Monitors.ExpectedDestinations.PrefixList
	Loaded.Monitors.ExpectedDestinations.MatchRegex = Parsed.Monitors.ExpectedDestinations.MatchRegex
	Loaded.Monitors.ExpectedDestinations.IgnoreRegex = Parsed.Monitors.ExpectedDestinations.IgnoreRegex

	Loaded.Monitors.SmallDurationCalls.Enabled = Parsed.Monitors.SmallDurationCalls.Enabled
	executeInterval, err = time.ParseDuration(Parsed.Monitors.SmallDurationCalls.ExecuteInterval)
	if err != nil {
		return fmt.Errorf("error parsing duration for \"execute_interval\" in \"monitors/small_duration_calls\"")
	}
	Loaded.Monitors.SmallDurationCalls.ExecuteInterval = executeInterval
	Loaded.Monitors.SmallDurationCalls.HitThreshold = Parsed.Monitors.SmallDurationCalls.HitThreshold
	Loaded.Monitors.SmallDurationCalls.MinimumNumberLength = Parsed.Monitors.SmallDurationCalls.MinimumNumberLength
	Loaded.Monitors.SmallDurationCalls.ActionChainName = Parsed.Monitors.SmallDurationCalls.ActionChainName
	Loaded.Monitors.SmallDurationCalls.ActionChainHoldoffPeriod = Parsed.Monitors.SmallDurationCalls.ActionChainHoldoffPeriod
	Loaded.Monitors.SmallDurationCalls.MaxActionChainRunCount = Parsed.Monitors.SmallDurationCalls.MaxActionChainRunCount
	if considerFromLast, err := time.ParseDuration(Parsed.Monitors.SmallDurationCalls.ConsiderCDRsFromLast); err != nil {
		considerFromLastUInt, err := strconv.Atoi(Parsed.Monitors.SmallDurationCalls.ConsiderCDRsFromLast)
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
	durationThreshold, err := time.ParseDuration(Parsed.Monitors.SmallDurationCalls.DurationThreshold)
	if err != nil {
		return fmt.Errorf("error parsing duration for \"duration_threshold\" in \"monitors/small_duration_calls\"")
	}
	Loaded.Monitors.SmallDurationCalls.DurationThreshold = durationThreshold

	// Actions
	fmt.Println(Parsed.Actions.Email)
	fmt.Println(*Parsed.Actions.Email)
	if Parsed.Actions.Email != nil {
		Loaded.Actions.Email = new(actionEmail)
		Loaded.Actions.Email.Enabled = Parsed.Actions.Email.Enabled
		Loaded.Actions.Email.Username = Parsed.Actions.Email.Username
		Loaded.Actions.Email.Password = Parsed.Actions.Email.Password
		Loaded.Actions.Email.Message = Parsed.Actions.Email.Message
	}
	if Parsed.Actions.Call != nil {
		Loaded.Actions.Call = new(actionCall)
		Loaded.Actions.Call.Enabled = Parsed.Actions.Call.Enabled
	}
	if Parsed.Actions.HTTP != nil {
		Loaded.Actions.HTTP = new(actionHTTP)
		Loaded.Actions.HTTP.Enabled = Parsed.Actions.HTTP.Enabled
	}
	if Parsed.Actions.LocalCommands != nil {
		Loaded.Actions.LocalCommands = new(actionLocalCommands)
		Loaded.Actions.LocalCommands.Enabled = Parsed.Actions.LocalCommands.Enabled
	}

	// Action Chains
	Loaded.ActionChains = make(map[string][]actionChainAction)
	for k, v := range *Parsed.ActionChains {
		Loaded.ActionChains[k] = v
	}

	// Data Groups
	Loaded.DataGroups = make(map[string]dataGroup)
	for k, v := range *Parsed.DataGroups {
		Loaded.DataGroups[k] = v
	}

	return nil

}

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
	Brand      string
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
