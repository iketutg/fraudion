package config

import (
	"os"

	"encoding/json"
	"path/filepath"

	"github.com/DisposaBoy/JsonConfigReader"
)

// Parsed After config.Parse(...) is called, this variable holds the values parsed from the JSON config file specified
var parsed *parsedValues

//var parsed interface{}

// Parse Parses the config file at "configDir" with name "configFileName" and puts the value on the config.Parsed variable
func parse(configDir string, configFileName string) error {

	// NOTE: JSON Related Help at https://github.com/DisposaBoy/JsonConfigReader, https://golang.org/pkg/encoding/json/, https://blog.golang.org/json-and-go
	// NOTE: Anything that can't be found will be saved to the JSON objects with empty values for the specified types

	configFile, err := os.Open(filepath.Join(configDir, configFileName))
	if err != nil {
		return err
	}
	defer configFile.Close()

	parsed = new(parsedValues)
	//parsed = new(map[string]interface{})

	// NOTE: Really? In the end it was only one line of code?
	if err = json.NewDecoder(JsonConfigReader.New(configFile)).Decode(parsed); err != nil {
		parsed = nil
		return err
	}

	return nil

}

type parsedValues struct {
	General      *generalJSON    `json:"general"`
	Softswitch   *softswitchJSON `json:"softswitch"`
	CDRsSources  *cdrsSources    `json:"cdrs_sources"`
	Monitors     *monitorsJSON   `json:"monitors"`
	Actions      *actionsJSON    `json:"actions"`
	ActionChains *actionChains   `json:"action_chains"`
	DataGroups   *dataGroups     `json:"data_groups"`
}

type generalJSON struct{}

type softswitchJSON struct {
	System     string
	Version    string
	CDRsSource map[string]string `json:"cdrs_source"`
}

//type cdrsSources map[string]map[string]string

type monitorsJSON struct {
	SimultaneousCalls     *monitorSimultaneousCallsJSON     `json:"*simultaneous_calls"`
	DangerousDestinations *monitorDangerousDestinationsJSON `json:"*dangerous_destinations"`
	ExpectedDestinations  *monitorExpectedDestinationsJSON  `json:"*expected_destinations"`
	SmallDurationCalls    *monitorSmallCallDurationsJSON    `json:"*small_duration_calls"`
}

type monitorBaseJSON struct {
	Enabled                  bool
	ExecuteInterval          string `json:"execute_interval"`
	HitThreshold             uint32 `json:"hit_threshold"`
	MinimumNumberLength      uint32 `json:"minimum_number_length"`
	ActionChainName          string `json:"action_chain_name"`
	ActionChainHoldoffPeriod uint32 `json:"action_chain_holdoff_period"`
	MaxActionChainRunCount   uint32 `json:"action_chain_run_count"`
}

type monitorSimultaneousCallsJSON struct {
	monitorBaseJSON
}
type monitorDangerousDestinationsJSON struct {
	monitorBaseJSON
	ConsiderCDRsFromLast string   `json:"consider_cdrs_from_last"`
	PrefixList           []string `json:"prefix_list"`
	MatchRegex           string   `json:"match_regex"`
	IgnoreRegex          string   `json:"ignore_regex"`
}
type monitorExpectedDestinationsJSON struct {
	monitorBaseJSON
	ConsiderCDRsFromLast string   `json:"consider_cdrs_from_last"`
	PrefixList           []string `json:"prefix_list"`
	MatchRegex           string   `json:"match_regex"`
	IgnoreRegex          string   `json:"ignore_regex"`
}
type monitorSmallCallDurationsJSON struct {
	monitorBaseJSON
	ConsiderCDRsFromLast string `json:"consider_cdrs_from_last"`
	DurationThreshold    string `json:"duration_threshold"`
}

type actionsJSON struct {
	Email         *actionEmailJSON         `json:"*email"`
	Call          *actionCallJSON          `json:"*call"`
	HTTP          *actionHTTPJSON          `json:"*http"`
	LocalCommands *actionLocalCommandsJSON `json:"*local_commands"`
}
type actionBaseJSON struct {
	Enabled   bool
	Recurrent bool
}
type actionEmailJSON struct {
	actionBaseJSON
	Username string `json:"gmail_username"`
	Password string `json:"gmail_password"`
	Message  string `json:"message"`
}
type actionCallJSON struct {
	actionBaseJSON
}
type actionHTTPJSON struct {
	actionBaseJSON
}
type actionLocalCommandsJSON struct {
	actionBaseJSON
}

//
// type actionChains map[string][]actionChainAction
//
// type actionChainAction struct {
// 	ActionName     string   `json:"action"`
// 	DataGroupNames []string `json:"data_groups"`
// }
//
// type dataGroups map[string]dataGroup
//
// type dataGroup struct {
// 	PhoneNumber      string            `json:"phone_number"`
// 	EmailAddress     string            `json:"data_groups"`
// 	HTTPURL          string            `json:"http_url"`
// 	HTTPMethod       string            `json:"http_method"`
// 	HTTPParameters   map[string]string `json:"data_groups"`
// 	CommandName      string            `json:"command_name"`
// 	CommandArguments string            `json:"command_arguments"`
// }
