package config

import (
	"os"

	"encoding/json"

	"github.com/DisposaBoy/JsonConfigReader"
)

var parsed *parsedValues

func parseFromFile(configFile *os.File) error {

	// NOTE: Validation reads through the whole file so we need to get back to the beginning.
	configFile.Seek(0, 0)

	parsed = new(parsedValues)

	// NOTE: Really? In the end it was only one line of code?
	if err := json.NewDecoder(JsonConfigReader.New(configFile)).Decode(parsed); err != nil {
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
	Type           string
	Version        string
	CDRsSourceName string `json:"cdrs_source_name"`
}

//type cdrsSources map[string]map[string]string // TODO: Moved to the load.go file

type monitorsJSON struct {
	SimultaneousCalls     *monitorSimultaneousCallsJSON     `json:"simultaneous_calls"`
	DangerousDestinations *monitorDangerousDestinationsJSON `json:"dangerous_destinations"`
	ExpectedDestinations  *monitorExpectedDestinationsJSON  `json:"expected_destinations"`
	SmallDurationCalls    *monitorSmallCallDurationsJSON    `json:"small_duration_calls"`
}

type monitorBaseJSON struct {
	Enabled             bool
	ExecuteInterval     string `json:"execute_interval"`
	HitThreshold        uint32 `json:"hit_threshold"`
	MinimumNumberLength uint32 `json:"minimum_number_length"`
	ActionChainName     string `json:"action_chain_name"`
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
	Email         *actionEmailJSON
	LocalCommands *actionLocalCommandsJSON `json:"local_commands"`
}

type actionBaseJSON struct {
	Enabled   bool
	Recurrent bool
}

type actionEmailJSON struct {
	actionBaseJSON
	Type     string
	Username string
	Password string
	Title    string
	Body     string
}

type actionLocalCommandsJSON struct {
	actionBaseJSON
}

// type actionChains map[string][]actionChainAction  // TODO: Moved to the load.go file

// type actionChainAction struct {  // TODO: Moved to the load.go file
// 	ActionName     string   `json:"action"`
// 	DataGroupNames []string `json:"data_groups"`
// }

// type dataGroups map[string]dataGroup  // TODO: Moved to the load.go file

// type dataGroup struct {  // TODO: Moved to the load.go file
// 	PhoneNumber      string            `json:"phone_number"`
// 	EmailAddress     string            `json:"data_groups"`
// 	HTTPURL          string            `json:"http_url"`
// 	HTTPMethod       string            `json:"http_method"`
// 	HTTPParameters   map[string]string `json:"data_groups"`
// 	CommandName      string            `json:"command_name"`
// 	CommandArguments string            `json:"command_arguments"`
// }
