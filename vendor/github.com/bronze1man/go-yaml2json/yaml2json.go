package yaml2json

import (
	"encoding/json"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v2"
)

func Convert(input []byte) ([]byte, error) {
	var yamlData interface{}
	err := yaml.Unmarshal(input, &yamlData)
	if err != nil {
		return nil, err
	}
	jsonData, err := convert(yamlData)
	if err != nil {
		return nil, err
	}
	output, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func convert(inputObj interface{}) (interface{}, error) {
	switch inputObj.(type) {
	case map[interface{}]interface{}:
		input := inputObj.(map[interface{}]interface{})
		output := make(map[string]interface{})
		for key, value := range input {
			var outputKey string
			switch key.(type) {
			case string:
				outputKey = key.(string)
			case int:
				outputKey = strconv.Itoa(key.(int))
			default:
				return nil, fmt.Errorf("Expected map key to be a string or int, but was %T", key)
			}
			outputValue, err := convert(value)
			if err != nil {
				return nil, err
			}
			output[outputKey] = outputValue
		}
		return output, nil
	case []interface{}:
		input := inputObj.([]interface{})
		output := make([]interface{}, len(input))
		for i, inputElem := range input {
			outputElem, err := convert(inputElem)
			if err != nil {
				return nil, err
			}
			output[i] = outputElem
		}
		return output, nil
	default:
		return inputObj, nil
	}
}
