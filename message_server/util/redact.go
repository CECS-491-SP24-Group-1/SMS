package util

import (
	"encoding/json"
)

/*
Marshalls a struct to JSON, removing any fields mentioned in a blacklist.
The field names are based off those mentioned in the json struct tags.
This function first marshals the struct to JSON, then unmarshals the JSON
to an intermediate map. Keep in mind that marshalling occurs three times:
initial marshal to JSON, then to an intermediate map with redacted fields
removed, and back out to JSON, so performance is likely impacted.
*/
func RedactJsonDM[T any](target T, whitelist bool, fieldNames ...string) ([]byte, error) {
	//Step 1: Marshal the struct to JSON
	jsonFp, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}

	//Step 2: Unmarshal the JSON to a map
	imap := make(map[string]interface{})
	if err := json.Unmarshal(jsonFp, &imap); err != nil {
		return nil, err
	}

	//Redact fields from the map
	imap = RedactMap(imap, whitelist, fieldNames...)

	//Step 3: Marshal back to JSON and return the result
	return json.Marshal(imap)
}

/*
Marshalls a struct to JSON, removing any fields mentioned in a blacklist.
The field names are based off those mentioned in the json struct tags.
This function uses `mapstructure` to marshal the struct to an intermediate
map.
*/
func RedactJsonMS[T any](target T, whitelist bool, fieldNames ...string) ([]byte, error) {
	//Step 1: Marshal the struct to a map using mapstructure
	imap := make(map[string]interface{})
	if err := MSRecursiveMarshal(target, &imap, "json"); err != nil {
		return nil, err
	}

	//Redact fields from the map
	imap = RedactMap(imap, whitelist, fieldNames...)

	//Step 2: Marshal back to JSON and return the result
	return json.Marshal(imap)
}

/*
Redacts fields from a map given a list of fields. In blacklist mode, this
entails removing all mentioned fields. In whitelist mode, this entails
removing all fields that are not mentioned. This function does not mutate
the source map.
*/
func RedactMap[T any](target map[string]T, whitelist bool, fieldNames ...string) map[string]T {
	// Check if the names are a whitelist or blacklist
	imap2 := target
	if whitelist {
		//Clear the map
		imap2 = make(map[string]T)

		//Copy only specified fields over
		for _, fieldName := range fieldNames {
			val, ok := target[fieldName]
			if ok {
				imap2[fieldName] = val
			}
		}
	} else {
		//Redact the specified fields
		for _, fieldName := range fieldNames {
			delete(imap2, fieldName) //This function is a NOP if the field doesn't exist
		}
	}
	return imap2
}
