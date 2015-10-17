package conf

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Custom errors
var (
	InvalidType      = errors.New("Invalid type given to the parser")
	InvalidFieldType = errors.New("Invalid field type")
	CannotParse      = errors.New("Some data cannot be parsed")
	LargeVal         = errors.New("Given value is too large")
)

// Prepaired data for filling into struct
type prepData struct {
	name      string         // This is arg name
	defVal    string         // This is default value
	desc      string         // This is description for arg
	field     *reflect.Value // Field in config struct
	flagValue *string        // This is parsed value from cmd-line
}

// This is the main argument parser function
func GetArguments(v interface{}) error {
	// Getting value
	val := reflect.ValueOf(v)

	// If the given value is not a pointer and not a struct returning error
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return InvalidType
	}

	// Parsing data
	if err := flagParser(val); err != nil {
		return err
	}

	return nil
}

// This func will parse and export data from command-line args into struct
func flagParser(v reflect.Value) error {
	// Accessing to values
	val := v.Elem()

	// Place to store all required flags
	reqFlags := []*ReqVal{}

	// Place to store prepaired and checked fields
	prepFields := []prepData{}

	// Crossing over all fields
	for i := 0; i < val.NumField(); i++ {
		// Getting field's value and type
		field := val.Field(i)
		tp := val.Type().Field(i)

		// Checking if field is exported
		if !field.CanSet() {
			continue
		}

		// Getting values for name, description and defualt tags
		descTag := tp.Tag.Get("description")
		nameTag := tp.Tag.Get("name")
		defValTag := tp.Tag.Get("default")

		// Checking for description and name tags
		if len(descTag) == 0 || len(nameTag) == 0 {
			return errors.New(fmt.Sprintf("Field: '%s', doesn't have description or name tag", tp.Name))
		}

		// Result argument
		var arg *string

		// Appending required flags to store
		if req, _ := strconv.ParseBool(tp.Tag.Get("required")); req {
			// Initializing req value
			reqVal := ReqVal{Default: defValTag, Name: nameTag}
			reqFlags = append(reqFlags, &reqVal)

			// Getting value from arg
			flag.Var(&reqVal, nameTag, descTag)
			arg = reqVal.Get()
		} else {
			arg = flag.String(nameTag, defValTag, descTag)
		}

		// Filling prepaired data
		resData := prepData{
			defVal:    defValTag,
			desc:      descTag,
			field:     &field,
			name:      nameTag,
			flagValue: arg,
		}

		// Appending prepaired data to list
		prepFields = append(prepFields, resData)

	}

	// Parsing argument values
	flag.Parse()

	// Checking if all values are parsed
	if !flag.Parsed() {
		return CannotParse
	}

	// Checking required arguments
	for _, r := range reqFlags {
		if !r.IsDefined() {
			return errors.New(fmt.Sprintf("Flag '%s', wasn't provided", r.Name))
		}
	}

	// At last filling config struct
	return fillConfig(prepFields)
}

// This func we need to fill config struct fields
func fillConfig(data []prepData) error {
	for _, d := range data {
		// Crossing over types
		switch d.field.Interface().(type) {
		// Parsing boolean values
		case bool:
			val, err := strconv.ParseBool(*d.flagValue)
			if err != nil {
				return err
			}
			d.field.SetBool(val)
		// Parsing datetime values
		case time.Duration:
			val, err := time.Parse("2006-02-01 15:04:05", *d.flagValue)
			if err != nil {
				return err
			}
			d.field.Set(reflect.ValueOf(val))
		// Parsing float values
		case float32, float64:
			val, err := strconv.ParseFloat(*d.flagValue, 64)
			if err != nil {
				return err
			}
			if d.field.OverflowFloat(val) {
				return LargeVal
			}
			d.field.SetFloat(val)
		// Parsing integer values
		case int, int8, int32, int64:
			val, err := strconv.ParseInt(*d.flagValue, 10, 64)
			if err != nil {
				return err
			}
			if d.field.OverflowInt(val) {
				return LargeVal
			}
			d.field.SetInt(val)
		// Parsing unsigned integers
		case uint, uint8, uint16, uint32, uint64:
			val, err := strconv.ParseUint(*d.flagValue, 10, 64)
			if err != nil {
				return err
			}
			if d.field.OverflowUint(val) {
				return LargeVal
			}
			d.field.SetUint(val)
		// Parsing string
		case string:
			d.field.SetString(*d.flagValue)
		// Handle unknown type
		default:
			return InvalidFieldType
		}
	}
	return nil
}
