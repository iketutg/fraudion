package config

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"encoding/json"

	"github.com/DisposaBoy/JsonConfigReader"
	v "github.com/gima/govalid/v1"
)

var configSchema = v.Object(

	// +INFO: https://github.com/gima/govalid

	v.ObjKV("general", v.Object(
		v.ObjKV("hostname", v.String(v.StrMin(5))),
	)),

	v.ObjKV("softswitch", v.Object(
		v.ObjKV("type", v.String(v.StrIs("*asterisk"))),
		v.ObjKV("version", v.String()),
		v.ObjKV("cdrs_source", v.Object(
			v.ObjKV("type", v.String(v.StrIs("*database"))),
			v.ObjKV("dbms", v.String()),
			v.ObjKV("user_name", v.String()),
			v.ObjKV("user_password", v.String()),
			v.ObjKV("database_name", v.String()),
			v.ObjKV("table_name", v.String()),
		)),
	)),

	// NOTE: Example of using this lib to make different possible combinations of fields on a section, when "live_calls_data_source" can be freeswitch the options list may be different...
	// v.ObjKV("live_calls_data_source", v.Or(
	// 	v.Object(
	// 		v.ObjKV("type", v.Or(v.String(v.StrIs("*asterisk")))),
	// 		v.ObjKV("version", v.String())),
	// 	v.Object(
	// 		v.ObjKV("type", v.Or(v.String(v.StrIs("*freeswitch")))),
	// 		v.ObjKV("sversion", v.String()),
	// 	),
	// )),

	v.ObjKV("monitors", v.Optional(v.Object(
		v.ObjKV("simultaneous_calls", v.Object(
			v.ObjKV("enabled", v.Boolean()),
			v.ObjKV("execute_interval", v.Function(validatorParseableDuration)),
			v.ObjKV("hit_threshold", v.Number(v.NumMin(1.0))),
			v.ObjKV("minimum_number_length", v.Number(v.NumMin(1.0))),
			v.ObjKV("action_chain_name", v.String()),
		)),

		v.ObjKV("dangerous_destinations", v.Optional(v.Object(
			v.ObjKV("enabled", v.Boolean()),
			v.ObjKV("execute_interval", v.Function(validatorParseableDuration)),
			v.ObjKV("hit_threshold", v.Number(v.NumMin(1.0))),
			v.ObjKV("minimum_number_length", v.Number(v.NumMin(1.0))),
			v.ObjKV("action_chain_name", v.String()),

			v.ObjKV("consider_cdrs_from_last", v.Function(validatorParseableDurationOrInt)),
			v.ObjKV("prefix_list", v.Array(v.ArrEach(v.String()))),
			v.ObjKV("match_regex", v.Function(validatorCompilableRegex)),
			v.ObjKV("ignore_regex", v.Function(validatorCompilableRegex)),
		))),

		v.ObjKV("expected_destinations", v.Optional(v.Object(
			v.ObjKV("enabled", v.Boolean()),
			v.ObjKV("execute_interval", v.Function(validatorParseableDuration)),
			v.ObjKV("hit_threshold", v.Number(v.NumMin(1.0))),
			v.ObjKV("minimum_number_length", v.Number(v.NumMin(1.0))),
			v.ObjKV("action_chain_name", v.String()),

			v.ObjKV("consider_cdrs_from_last", v.Function(validatorParseableDurationOrInt)),
			v.ObjKV("prefix_list", v.Array(v.ArrEach(v.String()))),
			v.ObjKV("match_regex", v.Function(validatorCompilableRegex)),
			v.ObjKV("ignore_regex", v.Function(validatorCompilableRegex)),
		))),

		v.ObjKV("small_duration_calls", v.Optional(v.Object(
			v.ObjKV("enabled", v.Boolean()),
			v.ObjKV("execute_interval", v.Function(validatorParseableDuration)),
			v.ObjKV("hit_threshold", v.Number(v.NumMin(1.0))),
			v.ObjKV("minimum_number_length", v.Number(v.NumMin(1.0))),
			v.ObjKV("action_chain_name", v.String()),

			v.ObjKV("consider_cdrs_from_last", v.Function(validatorParseableDurationOrInt)),
			v.ObjKV("duration_threshold", v.Function(validatorParseableDuration)),
			v.ObjKV("match_regex", v.Function(validatorCompilableRegex)),
			v.ObjKV("ignore_regex", v.Function(validatorCompilableRegex))),
		))),
	)),

	v.ObjKV("actions", v.Optional(v.Object(
		v.ObjKV("email", v.Optional(v.Object(
			v.ObjKV("enabled", v.Boolean()),
			v.ObjKV("recurrent", v.Boolean()),
			v.ObjKV("type", v.Or(v.String(v.StrIs("*gmail")))),
			v.ObjKV("username", v.String()),
			v.ObjKV("password", v.String()),
			v.ObjKV("title", v.String()),
			v.ObjKV("body", v.String()),
		))),

		v.ObjKV("local_commands", v.Optional(v.Object(
			v.ObjKV("enabled", v.Boolean()),
			v.ObjKV("recurrent", v.Optional(v.Boolean())),
		))),
	))),

	v.ObjKV("action_chains", v.Optional(v.Object(
		v.ObjKeys(v.String()),
		v.ObjValues(v.Array(v.ArrEach(v.Object(
			v.ObjKV("action_name", v.Or(v.String(v.StrIs("*email")), v.String(v.StrIs("*local_commands")))),
			v.ObjKV("data_groups", v.Array(v.ArrEach(v.String()))),
		)))),
	))),

	v.ObjKV("data_groups", v.Optional(v.Object(
		v.ObjKeys(v.String()),
		v.ObjValues(v.Object(
			// TODO/Future: Validate a Phone Number
			v.ObjKV("phone_number", v.Optional(v.String())),
			// TODO/Future: Validate an e-mail Address
			v.ObjKV("email_address", v.Optional(v.String())),
			// TODO/Future: Validate an URL
			v.ObjKV("http_url", v.Optional(v.String())),
			v.ObjKV("http_method", v.Optional(v.Or(v.String(v.StrIs("POST")), v.String(v.StrIs("GET"))))),
			v.ObjKV("http_parameters", v.Optional(v.Object(
				v.ObjKeys(v.String()),
				v.ObjValues(v.String()),
			))),
			v.ObjKV("command_name", v.Optional(v.String())),
			v.ObjKV("command_arguments", v.Optional(v.String())),
		)),
	))),
)

// ValidateFromFile ...
func ValidateFromFile(configFile *os.File) error {

	var data interface{}
	if err := json.NewDecoder(JsonConfigReader.New(configFile)).Decode(&data); err != nil {
		return err
	}

	return validateWithShema(data, configSchema)

}

// ValidateFromURL ...
func ValidateFromURL(reader io.Reader) error {

	var data interface{}
	if err := json.NewDecoder(JsonConfigReader.New(reader)).Decode(&data); err != nil {
		return err
	}

	fmt.Println(data)

	return validateWithShema(data, configSchema)

}

func validateWithShema(data interface{}, schema v.Validator) error {

	path, err := schema.Validate(data)
	if err == nil {
		return nil
	}

	return fmt.Errorf("failed validation at %s with error %s\n", path, err)

}

// Validation Funcions
func validatorParseableDuration(data interface{}) (path string, err error) {

	path = "validatorParseableDuration"

	validate, ok := data.(string)
	if !ok {
		return path, fmt.Errorf("expected string, got %v", reflect.TypeOf(data))
	}

	if _, err := time.ParseDuration(validate); err != nil {
		return path, fmt.Errorf("expected parseable time.Duration, got %s which isn't", validate)
	}

	return "", nil

}

func validatorParseableDurationOrInt(data interface{}) (path string, err error) {

	path = "validatorParseableDurationOrInt"

	validate, ok := data.(string)
	if !ok {
		return path, fmt.Errorf("expected string, got %v", reflect.TypeOf(data))
	}

	_, errNotDuration := time.ParseDuration(validate)
	integerValue, errNotInt := strconv.Atoi(validate)

	if errNotDuration != nil && errNotInt != nil {
		return path, fmt.Errorf("expected parseable time.Duration OR string convertible to int, got %s which isn't", validate)
	}

	if errNotDuration != nil && errNotInt == nil {
		if integerValue <= 0 {
			return path, fmt.Errorf("expected > 0 int, got %s", validate)
		}
	}

	return "", nil

}

func validatorCompilableRegex(data interface{}) (path string, err error) {

	path = "validatorCompilableRegex"

	validate, ok := data.(string)
	if !ok {
		return path, fmt.Errorf("expected string, got %v", reflect.TypeOf(data))
	}

	if _, err := regexp.Compile(validate); err != nil {
		return path, fmt.Errorf("expected compilable regex string, got %s which isn't", validate)
	}

	return "", nil

}
