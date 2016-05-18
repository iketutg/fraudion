package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	v "github.com/gima/govalid/v1"
)

// Validated ...
var validated bool

// Validate ...
func validate() error {

	validatorParseableDuration := func(data interface{}) (path string, err error) {

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

	validatorParseableDurationOrInt := func(data interface{}) (path string, err error) {

		path = "validatorCompilableRegex"

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

	validatorCompilableRegex := func(data interface{}) (path string, err error) {

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

	schema := v.Object(

		v.ObjKV("general", v.Optional(v.Object())),

		v.ObjKV("softswitch", v.Object(
			v.ObjKV("type", v.String()),
			v.ObjKV("version", v.Optional(v.String())),
			v.ObjKV("cdrs_source", v.Object(
				v.ObjKV("type", v.String()),
				v.ObjKV("name", v.String()),
			)),
		)),

		v.ObjKV("cdrs_sources", v.Object(
			v.ObjKeys(v.String()),
			v.ObjValues(v.Object(
				v.ObjKV("type", v.String()),
				v.ObjKV("dbms", v.String()),
				v.ObjKV("user_name", v.String()),
				v.ObjKV("user_password", v.String()),
				v.ObjKV("database_name", v.String()),
				v.ObjKV("table_name", v.String()),
			)),
		)),

		v.ObjKV("monitors", v.Object(
			v.ObjKV("simultaneous_calls", v.Object(
				v.ObjKV("enabled", v.Boolean()),
				v.ObjKV("execute_interval", v.Optional(v.Function(validatorParseableDuration))),
				v.ObjKV("hit_threshold", v.Optional(v.Number(v.NumMin(1.0)))),
				v.ObjKV("minimum_number_length", v.Optional(v.Number(v.NumMin(1.0)))),
				v.ObjKV("action_chain_name", v.Optional(v.String())),
			)),

			v.ObjKV("dangerous_destinations", v.Optional(v.Object(
				v.ObjKV("enabled", v.Boolean()),
				v.ObjKV("execute_interval", v.Optional(v.Function(validatorParseableDuration))),
				v.ObjKV("hit_threshold", v.Optional(v.Number(v.NumMin(1.0)))),
				v.ObjKV("minimum_number_length", v.Optional(v.Number(v.NumMin(1.0)))),
				v.ObjKV("action_chain_name", v.Optional(v.String())),

				v.ObjKV("consider_cdrs_from_last", v.Function(validatorParseableDurationOrInt)),
				v.ObjKV("prefix_list", v.Array(v.ArrEach(v.String()))),
				v.ObjKV("match_regex", v.Function(validatorCompilableRegex)),
				v.ObjKV("ignore_regex", v.Function(validatorCompilableRegex)),
			))),

			v.ObjKV("expected_destinations", v.Optional(v.Object(
				v.ObjKV("enabled", v.Boolean()),
				v.ObjKV("execute_interval", v.Optional(v.Function(validatorParseableDuration))),
				v.ObjKV("hit_threshold", v.Optional(v.Number(v.NumMin(1.0)))),
				v.ObjKV("minimum_number_length", v.Optional(v.Number(v.NumMin(1.0)))),
				v.ObjKV("action_chain_name", v.Optional(v.String())),

				v.ObjKV("consider_cdrs_from_last", v.Function(validatorParseableDurationOrInt)),
				v.ObjKV("prefix_list", v.Array(v.ArrEach(v.String()))),
				v.ObjKV("match_regex", v.Function(validatorCompilableRegex)),
				v.ObjKV("ignore_regex", v.Function(validatorCompilableRegex)),
			))),

			v.ObjKV("small_duration_calls", v.Optional(v.Object(
				v.ObjKV("enabled", v.Boolean()),
				v.ObjKV("execute_interval", v.Optional(v.Function(validatorParseableDuration))),
				v.ObjKV("hit_threshold", v.Optional(v.Number(v.NumMin(1.0)))),
				v.ObjKV("minimum_number_length", v.Optional(v.Number(v.NumMin(1.0)))),
				v.ObjKV("action_chain_name", v.Optional(v.String())),

				v.ObjKV("consider_cdrs_from_last", v.Function(validatorParseableDurationOrInt)),
				v.ObjKV("duration_threshold", v.Function(validatorParseableDuration))),
			))),
		),

		v.ObjKV("actions", v.Optional(v.Object(
			v.ObjKV("email", v.Optional(v.Object(
				v.ObjKV("enabled", v.Boolean()),
				v.ObjKV("recurrent", v.Optional(v.Boolean())),
				v.ObjKV("gmail_username", v.String()),
				v.ObjKV("gmail_password", v.String()),
				v.ObjKV("title", v.Optional(v.String())),
				v.ObjKV("body", v.Optional(v.String())),
			))),

			v.ObjKV("local_commands", v.Optional(v.Object(
				v.ObjKV("enabled", v.Boolean()),
				v.ObjKV("recurrent", v.Optional(v.Boolean())),
			))),
		))),

		v.ObjKV("action_chains", v.Optional(v.Object(
			v.ObjKeys(v.String()),
			v.ObjValues(v.Array(v.ArrEach(v.Object(
				v.ObjKV("action", v.Or(v.String(v.StrIs("*email")), v.String(v.StrIs("*local_commands")))),
				v.ObjKV("data_groups", v.Array(v.ArrEach(v.String()))),
			)))),
		))),

		v.ObjKV("data_groups", v.Optional(v.Object(
			v.ObjKeys(v.String()),
			v.ObjValues(v.Object(
				v.ObjKV("phone_number", v.String()),  // TODO: Validate a Phone Number
				v.ObjKV("email_address", v.String()), // TODO: Validate an e-mail Address
				v.ObjKV("http_url", v.String()),      // TODO: Validate an URL
				v.ObjKV("http_method", v.Or(v.String(v.StrIs("POST")), v.String(v.StrIs("GET")))),
				v.ObjKV("http_parameters", v.Object(
					v.ObjKeys(v.String()),
					v.ObjValues(v.String()),
				)),
				v.ObjKV("command_name", v.String()),
				v.ObjKV("command_arguments", v.String()),
			)),
		))),
	)

	path, err := schema.Validate(parsed)
	if err == nil {
		validated = true
		return nil
	}

	return fmt.Errorf("Failed Validation at %s with error %s.\n", path, err)

}
