package config

import (
//"fmt"
//"github.com/andmar/fraudion/utils"
)

// Validated ...
var Validated bool

// Validate ...
func validate() (bool, []error) {

	//var errors []error

	// General

	// Softswitch
	/*
		"softswitch": {
			"brand": "*asterisk",
			"version": "1.8",
			"cdrs_source": "*db_mysql",
		},
		name: softswitch, mandatory
		values: bellow...", everything else is ignored in parsing

		brand: string, mandatory, one of *asterisk, *softswitch (FUTURE)
		version: string, optional, any string, validated on Parsing
		cdrs_source: string, mandatory, one of *db_mysql, *csv (FUTURE)
	*/
	/*fmt.Printf("\n\n\n\n\n\n\n")
	fmt.Println(Parsed.Softswitch)
	fmt.Println(Parsed.CDRsSources)
	if Parsed.Softswitch.Brand == "" && Parsed.Softswitch.CDRsSource == "" {
		errors = append(errors, fmt.Errorf("\"softswitch\" section not found OR all of it's values are empty"))
	} else {
		if Parsed.Softswitch.Brand == "" {
			errors = append(errors, fmt.Errorf("\"brand\" value in \"softswitch\" section is blank"))
		} else {
			if utils.StringInStringsSlice(Parsed.Softswitch.Brand, []string{"*asterisk", "*freeswitch"}) == false {
				errors = append(errors, fmt.Errorf("\"brand\" value in \"softswitch\" section must be one of: *asterisk, *freeswitch"))
			}
		}
		if Parsed.Softswitch.CDRsSource == "" {
			errors = append(errors, fmt.Errorf("\"cdrs_source\" value in \"softswitch\" section is blank"))
		} else {
			if utils.StringInStringsSlice(Parsed.Softswitch.CDRsSource, []string{"*db_mysql"}) == false {
				errors = append(errors, fmt.Errorf("\"cdrs_sources\" value in \"softswitch\" section must be one of: *db_mysql"))
			} else {
				if _, found := *Parsed.CDRsSources[Parsed.Softswitch.CDRsSource]; found == false {
					errors = append(errors, fmt.Errorf("\"cdrs_sources\" value in \"softswitch\" not configured in \"cdrs_sources\""))
				}
			}
		}
	}

	// CDRs Sources
	fmt.Printf("\n\n")
	fmt.Println(Parsed.CDRsSources)
	for key, cdrSource := range *Parsed.CDRsSources {
		if utils.StringInStringsSlice(key, []string{"*db_mysql"}) {
			errors = append(errors, fmt.Errorf("CDR Source type \"%s\" unknown in \"cdrs_sources\" section", key))
		} else {
			if cdrSource["user_name"] == "" && cdrSource["user_password"] == "" && cdrSource["database_name"] == "" && cdrSource["table_name"] == "" {
				errors = append(errors, fmt.Errorf("\"cdrs_sources\" section not found OR all of it's values are empty"))
			} else {
				if cdrSource["user_name"] == "" || cdrSource["user_password"] == "" || cdrSource["database_name"] == "" || cdrSource["table_name"] == "" {
					errors = append(errors, fmt.Errorf("Each source in \"cdrs_sources\" must have username, user_password, database_name, table_name set"))
				}
			}
		}
	}*/

	/*
		   "*db_mysql": {
		   	"user_name": "",
		   	"user_password": "",
		   	"database_name": "asteriskcdrdb",
		   	"table_name": "cdr",
		   	"mysql_options": "allowOldPasswords=1"
		   }
			 name: cdrs_sources, mandatory if the following monitors are enabled: *dangerous_destinations, *expected_destinations, *small_duration_calls
			 values: one of *db_mysql, *csv (FUTURE), everything else is ignored on parsing

			 *db_mysql
		   user_name: string, mandatory
		   user_password: mandatory
		   database_name: mandatory
		   table_name: mandatory
		   mysql_options: optional
	*/

	// Monitors
	/*
		   	"monitors": {
		       "*simultaneous_calls": {
		         "enabled": true,
		         "execute_interval": "5m",
		         "hit_threshold": 5,
		         "minimum_number_length": 5,
		         "action_chain_name": "*default",
		   			"action_chain_holdoff_period": 0,
		   			"action_chain_run_count": 0,
		       },

		       "*dangerous_destinations": {
		         "enabled": true,
		         "execute_interval": "1m",
		         "hit_threshold": 5,
		         "minimum_number_length": 5,
		         "action_chain_name": "*default",
		   			"action_chain_holdoff_period": 0,
		   			"action_chain_run_count": 0,
		   			"consider_cdrs_from_last": "5",
		         "prefix_list": ["351", "244", "91", "53", "256", "48"],
		         "match_regex": "([0-9]{0,8})?(0{2})?__prefix__[0-9]{5,}",
		         "ignore_regex": "^[0-9]{9}$"
		   		},

		       "*expected_destinations": {
		         "enabled": true,
		         "execute_interval": "5m",
		         "hit_threshold": 5,
		         "minimum_number_length": 10,
		         "action_chain_name": "*default",
		   			"action_chain_holdoff_period": 0,
		   			"action_chain_run_count": 0,
		   			"consider_cdrs_from_last": "5d",
		         "prefix_list": ["244"],
		         "match_regex": "([0-9]{0,8})?(0{2})?__prefix__[0-9]{5,}",
		         "ignore_regex": "^[0-9]{9}$"
		       },

		       "*small_duration_calls": {
		         "enabled": true,
		         "execute_interval": "5m",
		         "hit_threshold": 5,
		         "minimum_number_length": 5,
		         "action_chain_name": "*default",
		   			"action_chain_holdoff_period": 0,
		   			"action_chain_run_count": 0,
		   			"consider_cdrs_from_last": "5d",
		         "duration_threshold": "5s"
		       }
		     },
				 name: monitores
				 values: one of *simultaneous_calls, *dangerous_destinations, *small_duration_calls, *expected_destinations, everything else is ignored in Parsing

				 "common":
				 enabled: bool, mandatory
				 execute_interval: string with parseable time.Duration e.g. "5m", optional (has default value of ?)
				 hit_threshold: uint, optional (has default value of ?), > 0
				 minimum_number_length: uint, optional (has default value of ?), > 0
				 action_chain_name: string, optional (defaults to "*default"), string has to match one of the setup Action Chains
				 action_chain_holdoff_period: TBS
				 action_chain_run_count: uint, optional (has default value of ?), > 0

				 *simultaneous_calls

				 *dangerous_destinations
				 consider_cdrs_from_last: string, optional, parseable duration or number of days e.g. "5", "5h"
				 prefix_list: [string], mandatory e.g. ["351", "244", "91", "53", "256", "48"]
				 match_regex: string/regex, mandatory e.g. "([0-9]{0,8})?(0{2})?__prefix__[0-9]{5,}" __prefix__ will be replaced by each item on prefix_list on checking
				 ignore_regex: string/regex, mandatory e.g. "^[0-9]{9}$" TBS: __prefix__ will be replaced by each item on prefix_list on checking

				 *small_duration_calls
				 consider_cdrs_from_last: string, optional, parseable duration or number of days e.g. "5", "5h"
				 duration_threshold: string, optional, parseable duration e.g. "5s"

				 *expected_destinations
				 consider_cdrs_from_last: string, optional, parseable duration or number of days e.g. "5", "5h"
				 prefix_list: [string], mandatory e.g. ["351", "244", "91", "53", "256", "48"]
				 match_regex: string/regex, mandatory e.g. "([0-9]{0,8})?(0{2})?__prefix__[0-9]{5,}" __prefix__ will be replaced by each item on prefix_list on checking
				 ignore_regex: string/regex, mandatory e.g. "^[0-9]{9}$" TBS: __prefix__ will be replaced by each item on prefix_list on checking

	*/

	//fmt.Printf("Errors:\n%s\n", errors)

	// if len(errors) != 0 {
	// 	return true, errors
	// }

	return false, nil

	// Actions
	/*
		   	"actions": {

		       "*email": {
		         "enabled": true, // [Optional] If omitted this we consider it "disabled"
		         "gmail_username": "username@domain",
		         "gmail_password": "password",
		         "message": "This is a message, we support some __tags__ that we replace with information."
		       },
		       "*http": {
		         "enabled": true, // [Optional] If omitted we consider it "disabled"
		       },
		       "*call": {
		         "enabled": true // [Optional] If omitted we consider it "disabled"
		       },
		       "*local_commands": { // You can define your own command actions by giving them a name and a string that will be executed on the system! "*local_command x N"
		         "enabled": true // [Optional] If omitted we consider it "disabled"
		       }

		     },
		   	name: actions
		   	values: one of *email, *http, *local_comands, *call, everything else is ignored in Parsing

		   	*email
				enabled: bool, optional (has default value of false)
				gmail_username: string, mandatory if enabled
				gmail_password: string, mandatory if enabled
				message: string, optional

		   	*http
				enabled: bool, optional (has default value of false)

		   	*call
				enabled: bool, optional (has default value of false)

		   	*local_comands
				enabled: bool, optional (has default value of false)

	*/

	// Action Chains
	/*
		   "action_chains": {

		   		"*default": [
		   			{
		   				"action": "*email",
		   				"data_groups": ["DataGroupName", "DataGroup2Name"]
		   			},
		   			{
		   				"action": "*call",
		   				"data_groups": ["DataGroupName"]
		   			},
		   			{
		   				"action": "*localcommand",
		   				"data_groups": ["DataGroupName"]
		   			},
		   			// etc...
		   		],
		   		"OneRandomName": [
		   			{
		   				"action": "*call",
		   				"data_groups": ["DataGroupName", "DataGroup2Name"]
		   			}
		   		],
		   		// etc...

		   },
		   name: action_chains
		   values: *default or custom strings, default is mandatory if any of the monitors is enabled

			 each value is an array of:
			 action: string, mandatory, one of the defined actions
			 data_groups: array of string, mandatory, one or more of the defined Data Groups
	*/

	// Data Groups
	/*
					"data_groups": {

							"DataGroupName": {
								"phone_number": "003519347396460",
								"email_address": "username@domain",
								"http_url": "api.somedomain.com/fraudion_in",
								"http_method": "POST",
								"http_parameters": {
									"http_post_parameters_1_k": "http_post_parameters_1_v",
									"http_post_parameters_2_k": "http_post_parameters_2_v"
									// etc...
								},
								"command_name": "amportal",
								"command_arguments": "stop"
							}


					}
					name: data_groups
					values: *default or custom strings (DataGroup Name), default is mandatory if any of the monitors is enabled

					For each value:
					the mandatory fiels depend on the type of actions that use the data_group on action chains


			Dependencies:
			Data Groups
			Action Chains: Check if Action exists > Check if Data Group Exists, check if Data Group with Name has the required values for the Action
			Monitors: Check if Action Chain exists
		  CDRs Sources:
		  Softswitch: Check if CDRs Source exists

	*/

	//return nil
}
