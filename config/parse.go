package config

import (
	"os"
	//"reflect"

	"encoding/json"
	"path/filepath"

	"github.com/DisposaBoy/JsonConfigReader"
)

// Parsed After config.Parse(...) is called this variable holds the values parsed from the config file specified
var Parsed *parsedValues

//var types = make(map[string]reflect.Type) // NOTE: The purpose of this is to use a local function bellow to simplify the JSON Parsing of the various sections

func init() {
	// types["generalJSON"] = reflect.TypeOf(generalJSON{})
	// types["softswitchJSON"] = reflect.TypeOf(softswitchJSON{})
	// types["monitorsJSON"] = reflect.TypeOf(monitorsJSON{})
	// types["actionsJSON"] = reflect.TypeOf(actionsJSON{})
}

// Parse Parses the config file at "configDir" with name "configFileName" and puts the value on the config.Parsed variable
func Parse(configDir string, configFileName string) error {

	// NOTE: JSON Related Help at https://github.com/DisposaBoy/JsonConfigReader, https://golang.org/pkg/encoding/json/, https://blog.golang.org/json-and-go
	// NOTE: ANything that can't be found will be saved to the JSON objects with empty values for the specified types

	configFile, err := os.Open(filepath.Join(configDir, configFileName))
	if err != nil {
		return err
	}
	defer configFile.Close()

	Parsed = new(parsedValues)

	// NOTE: Really? In the end it was only one line of code?
	if err = json.NewDecoder(JsonConfigReader.New(configFile)).Decode(Parsed); err != nil {
		Parsed = nil
		return err
	}

	return nil

	/*
		// NOTE: This was the old better way of Parsing the JSON config file

		configs := new(Parsed)

		var rawJSON map[string]*json.RawMessage // NOTE: Better than using the Lib example's Empty interface... https://tour.golang.org/methods/14
		if err = json.NewDecoder(JsonConfigReader.New(configFile)).Decode(&rawJSON); err != nil {
			fmt.Println("Error:", err)
		}

		//fmt.Println(configsQuick)
		os.Exit(-1)

		unmarshalJSONIntoObject := func(sectionName string, typeName string, getObjectPointer func(object interface{}) interface{}) error {
			reflectType, hasKey := types[typeName]
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
		//
		// if err := unmarshalIntoObject("general", "generalJSON",
		// 	func(object interface{}) interface{} {
		// 		configs.General = object.(generalJSON)
		// 		return &configs.General
		// 	}); err != nil {
		// 	log.LogO("ERROR", fmt.Sprintf("%s", err), marlog.OptionFatal)
		// }
		//
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
		//
		// if err := unmarshalJSONIntoField("cdrs_sources", &configs.CDRsSources); err != nil {
		// 	log.LogO("ERROR", fmt.Sprintf("%s", err), marlog.OptionFatal)
		// }
		//

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
	*/

}

// parsedValues ...
type parsedValues struct {
	General      generalJSON                    `json:"general"`
	Softswitch   softswitchJSON                 `json:"softswitch"`
	CDRsSources  map[string]map[string]string   `json:"cdrs_sources"`
	Monitors     monitorsJSON                   `json:"monitors"`
	Actions      actionsJSON                    `json:"actions"`
	ActionChains map[string][]actionChainAction `json:"action_chains"`
	DataGroups   map[string]dataGroup           `json:"data_groups"`
}

// generalJSON ...
type generalJSON struct {
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

// actionsJSON ...
type actionsJSON struct {
	Email         actionEmailJSON         `json:"*email"`
	Call          actionCallJSON          `json:"*call"`
	HTTP          actionHTTPJSON          `json:"*http"`
	LocalCommands actionLocalCommandsJSON `json:"*local_commands"`
}
type actionBaseJSON struct {
	Enabled bool
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
