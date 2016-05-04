package config

import (
	"fmt"
	"os"
	"reflect"

	"encoding/json"
	"path/filepath"

	"github.com/andmar/marlog"

	"github.com/DisposaBoy/JsonConfigReader"
)

var parsedConfig *Parsed
var typeRegistry = make(map[string]reflect.Type)

func init() {
	typeRegistry["generalJSON"] = reflect.TypeOf(generalJSON{})
	typeRegistry["softswitchJSON"] = reflect.TypeOf(softswitchJSON{})
	typeRegistry["monitorsJSON"] = reflect.TypeOf(monitorsJSON{})
	typeRegistry["actionsJSON"] = reflect.TypeOf(actionsJSON{})
}

// Parse ...
func Parse(configDir string, configFileName string) error {

	log := marlog.MarLog
	configs := new(Parsed)

	// From JSON config file to map[string]"RawJSON"
	configFile, err := os.Open(filepath.Join(configDir, configFileName))
	if err != nil {
		os.Exit(-1)
	}
	defer configFile.Close()

	// NOTE: JSON Related Help at https://github.com/DisposaBoy/JsonConfigReader, https://golang.org/pkg/encoding/json/, https://blog.golang.org/json-and-go

	var rawJSON map[string]*json.RawMessage // NOTE: Better than using the Lib example's Empty interface example... https://tour.golang.org/methods/14
	if err = json.NewDecoder(JsonConfigReader.New(configFile)).Decode(&rawJSON); err != nil {

	}

	unmarshalJSONIntoObject := func(sectionName string, typeName string, getObjectPointer func(object interface{}) interface{}) error {
		reflectType, hasKey := typeRegistry[typeName]
		if hasKey == false {
			return fmt.Errorf("Destination object is nil\n")
		}
		sectionRawJSON, hasKey := rawJSON[sectionName]
		if hasKey == false {
			return fmt.Errorf("Section not found in JSON\n")
		}
		if err := json.Unmarshal(*sectionRawJSON, getObjectPointer(reflect.New(reflectType).Elem().Interface())); err != nil {
			return fmt.Errorf("Unmarshal error: %s\n", err)
		}
		return nil
	}
	// NOTE: Usage example:
	/*
		if err := unmarshalIntoObject("general", "generalJSON",
			func(object interface{}) interface{} {
				configs.General = object.(generalJSON)
				return &configs.General
			}); err != nil {
			log.LogO("ERROR", fmt.Sprintf("%s", err), marlog.OptionFatal)
		}
	*/
	unmarshalJSONIntoField := func(sectionName string, field interface{}) error {
		sectionRawJSON, hasKey := rawJSON[sectionName]
		if hasKey == false {
			return fmt.Errorf("Section not found in JSON\n")
		}
		if err := json.Unmarshal(*sectionRawJSON, field); err != nil {
			return fmt.Errorf("Unmarshal error: %s\n", err)
		}
		return nil
	}
	// NOTE: Usage example:
	/*
		if err := unmarshalJSONIntoField("cdrs_sources", &configs.CDRsSources); err != nil {
			log.LogO("ERROR", fmt.Sprintf("%s", err), marlog.OptionFatal)
		}
	*/

	// ** General Section
	if err := unmarshalJSONIntoObject("general", "generalJSON",
		func(object interface{}) interface{} {
			configs.General = object.(generalJSON)
			return &configs.General
		}); err != nil {
		log.LogO("ERROR", fmt.Sprintf("%s", err), marlog.OptionFatal)
	}

	// ** Softswitch Section
	if err := unmarshalJSONIntoObject("softswitch", "softswitchJSON",
		func(object interface{}) interface{} {
			configs.Softswitch = object.(softswitchJSON)
			return &configs.Softswitch
		}); err != nil {
		log.LogO("ERROR", fmt.Sprintf("%s", err), marlog.OptionFatal)
	}

	// ** CDRs Sources Section
	if err := unmarshalJSONIntoField("cdrs_sources", &configs.CDRsSources); err != nil {
		log.LogO("ERROR", fmt.Sprintf("%s", err), marlog.OptionFatal)
	}

	// ** Monitors
	if err := unmarshalJSONIntoObject("monitors", "monitorsJSON",
		func(object interface{}) interface{} {
			configs.Monitors = object.(monitorsJSON)
			return &configs.Monitors
		}); err != nil {
		log.LogO("ERROR", fmt.Sprintf("%s", err), marlog.OptionFatal)
	}

	// ** Actions Section
	if err := unmarshalJSONIntoObject("actions", "actionsJSON",
		func(object interface{}) interface{} {
			configs.Actions = object.(actionsJSON)
			return &configs.Actions
		}); err != nil {
		log.LogO("ERROR", fmt.Sprintf("%s", err), marlog.OptionFatal)
	}

	// ** Action Chains Section
	if err := unmarshalJSONIntoField("action_chains", &configs.ActionChains); err != nil {
		log.LogO("ERROR", fmt.Sprintf("%s", err), marlog.OptionFatal)
	}

	// ** Data Groups Section
	if err := unmarshalJSONIntoField("data_groups", &configs.DataGroups); err != nil {
		log.LogO("ERROR", fmt.Sprintf("%s", err), marlog.OptionFatal)
	}

	fmt.Println("General:", configs.General)
	fmt.Println("Softswitch:", configs.Softswitch)
	fmt.Println("CDRs Sources:", configs.CDRsSources)
	fmt.Println("Monitors:", configs.Monitors)
	fmt.Println("Actions:", configs.Actions)
	fmt.Println("Action Chains:", configs.ActionChains)
	fmt.Println("Data Groups:", configs.DataGroups)

	return nil

}

// Parsed ...
type Parsed struct {
	General      generalJSON
	Softswitch   softswitchJSON
	CDRsSources  map[string]map[string]string
	Monitors     monitorsJSON
	Actions      actionsJSON
	ActionChains map[string][]actionChainAction
	DataGroups   map[string]dataGroup
}

// generalJSON ...
type generalJSON struct {
	CDRsSource                            string `json:"cdrs_source"`
	DefaultTriggerExecuteInterval         string `json:"default_trigger_execute_interval"`
	DefaultHitThreshold                   uint32 `json:"default_hit_threshold"`
	DefaultMinimumDestinationNumberLength uint32 `json:"default_minimum_destination_number_length"`
	DefaultActionChainHoldoffPeriod       string `json:"default_action_chain_holdoff_period"`
	DefaultActionChainRunCount            uint32 `json:"default_action_chain_run_count"`
}

// softswitchJSON ...
type softswitchJSON struct {
	Brand      string
	Version    string
	CDRsSource string `json:"cdrs_source"`
}

// monitoresJSON ...
type monitorsJSON struct {
	SimultaneousCalls     monitorSimultaneousCallsJSON     `json:"*simultaneous_calls"`
	DangerousDestinations monitorDangerousDestinationsJSON `json:"*dangerous_destinations"`
	ExpectedDestinations  monitorExpectedDestinationsJSON  `json:"*expected_destinations"`
	SmallDurationCalls    monitorSmallCallDurationsJSON    `json:"*small_duration_calls"`
}

type monitorBase struct {
	Enabled                  bool   `json:"enabled"`
	ExecuteInterval          string `json:"execute_interval"`
	HitThreshold             uint32 `json:"hit_threshold"`
	MinimumNumberLength      uint32 `json:"minimum_number_length"`
	ActionChainName          string `json:"action_chain_name"`
	ActionChainHoldoffPeriod uint32 `json:"action_chain_holdoff_period"`
	MaxActionChainRunCount   uint32 `json:"action_chain_run_count"`
}

type monitorSimultaneousCallsJSON struct {
	monitorBase
}
type monitorDangerousDestinationsJSON struct {
	monitorBase
	ConsiderCDRsFromLast string   `json:"consider_cdrs_from_last"`
	PrefixList           []string `json:"prefix_list"`
	MatchRegex           string   `json:"match_regex"`
	IgnoreRegex          string   `json:"ignore_regex"`
}
type monitorExpectedDestinationsJSON struct {
	monitorBase
	ConsiderCDRsFromLast string   `json:"consider_cdrs_from_last"`
	PrefixList           []string `json:"prefix_list"`
	MatchRegex           string   `json:"match_regex"`
	IgnoreRegex          string   `json:"ignore_regex"`
}
type monitorSmallCallDurationsJSON struct {
	monitorBase
	ConsiderCDRsFromLast string `json:"consider_cdrs_from_last"`
	DurationThreshold    string `json:"duration_threshold"`
}

// actionsJSON ...
type actionsJSON struct {
	Email         actionEmailJSON         `json:"*email"`
	Call          actionCallJSON          `json:"*call"`
	HTTP          actionHTTPJSON          `json:"*http"`
	LocalCommands actionLocalCommandsJSON `json:"*local_commands"`
}

type actionEmailJSON struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"gmail_username"`
	Password string `json:"gmail_password"`
	Message  string `json:"message"`
}

type actionCallJSON struct {
	Enabled bool `json:"enabled"`
}

type actionHTTPJSON struct {
	Enabled bool `json:"enabled"`
}

type actionLocalCommandsJSON struct {
	Enabled bool `json:"enabled"`
}
