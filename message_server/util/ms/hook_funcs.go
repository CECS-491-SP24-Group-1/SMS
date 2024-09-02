package ms

import (
	"reflect"
	"time"

	"github.com/go-viper/mapstructure/v2"
)

const (
	_TimeKey = "__ms__Time_Key"
)

var (
	//Controls the formatting of marshalled times.
	TimeFmt = time.RFC3339Nano
)

// Custom hook function to handle time.Time marshaling.
func timeToStringHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		//Skip non-time types
		if f != reflect.TypeOf(time.Time{}) && f != reflect.TypeOf(&time.Time{}) {
			return data, nil
		}

		//Setup the output time object
		var tim time.Time

		//Check if it's a time.Time
		if cnv, ok := data.(time.Time); ok {
			tim = cnv
		}

		//Check if it's a *time.Time
		if cnv, ok := data.(*time.Time); ok && cnv != nil {
			tim = *cnv
		}

		//Encode to a map if the output target is a map
		//`time.Time` requires this, but not `*time.Time`
		if t.Kind() == reflect.Map {
			mp := make(map[interface{}]interface{})
			mp[_TimeKey] = tim.Format(TimeFmt)
			return mp, nil
		}

		//Do not encode if the output target is not a string
		if t.Kind() != reflect.String {
			return data, nil
		}

		//Format the time and return it
		return tim.Format(TimeFmt), nil
	}
}
