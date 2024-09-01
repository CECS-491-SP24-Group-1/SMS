package ms

import "github.com/go-viper/mapstructure/v2"

// Struct -> Map adapter for `MapStructure` that recursively marshalls embedded structs.
func MSRecursiveMarshal[T any](input T, output *map[string]interface{}, tagName string) error {
	//Marshal to a map
	err := msMarshaller(
		input, output, tagName,
		mapstructure.ComposeDecodeHookFunc(mapstructure.RecursiveStructToMapHookFunc(), timeToStringHookFunc()),
	)
	if err != nil {
		return err
	}

	//Apply post-marshal filters
	filterTMaps(output)
	return nil
}

// Struct -> Map adapter for `MapStructure` that hooks `TextMarshaller` for custom types.
func MSTextMarshal[T any](input T, output *map[string]interface{}, tagName string) error {
	//Marshal to a map
	err := msMarshaller(
		input, output, tagName,
		mapstructure.ComposeDecodeHookFunc(mapstructure.TextUnmarshallerHookFunc(), timeToStringHookFunc()),
	)
	if err != nil {
		return err
	}

	//Apply post-marshal filters
	filterTMaps(output)
	return nil
}

// Map -> Struct adapter for `MapStructure` that recursively marshalls embedded maps.
func MSRecursiveUnmarshal[T any](input map[string]interface{}, output *T, tagName string) error {
	return msMarshaller(
		input, output, tagName,
		mapstructure.ComposeDecodeHookFunc(mapstructure.RecursiveStructToMapHookFunc(), timeToStringHookFunc()),
	)
}

// Map -> Struct adapter for `MapStructure` that hooks `TextUnmarshaller` for custom types.
func MSTextUnmarshal[T any](input map[string]interface{}, output *T, tagName string) error {
	return msMarshaller(
		input, output, tagName,
		mapstructure.ComposeDecodeHookFunc(mapstructure.TextUnmarshallerHookFunc(), timeToStringHookFunc()),
	)
}

// Backend function for `MSTextMarshal()` and `MSTextUnmarshal()`
func msMarshaller[T any, U any](input T, output *U, tagName string, hookFunc mapstructure.DecodeHookFunc) error {
	//Setup the decoder config
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		TagName:          tagName,
		DecodeHook:       hookFunc,
	}

	//Setup the decoder
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	//Decode the input to the output pointer
	return decoder.Decode(input)
}

/*
Filters "time maps" in output maps. This gets around the issue where
marshaling `time.Time` requires a map to work properly.
*/
func filterTMaps(m *map[string]interface{}) {
	//Loop over the map
	for k, value := range *m {
		//If the current value is a map, then continue
		if innerMap, ok := value.(map[string]interface{}); ok {
			//Check if the map has a value for the "time key"
			if v, exists := innerMap[_TimeKey]; exists {
				//Strip away the embedded map and replace the map with the actual value
				delete(*m, k)
				(*m)[k] = v
			}
		}
	}
}
