/*
Convert YAML to JSON.

Inspired by https://github.com/bronze1man/yaml2json, a command line
tool with the same effect.

	import (
		"encoding/json"
		"io/ioutil"
		"os"

		"github.com/peter-edge/go-yaml2json"
	)

	func ReadYamlToJson(yamlFilePath string) (interface{}, error) {
		yamlFile, err := os.Open(yamlFilePath)
		if err != nil {
			return nil, err
		}
		defer yamlFile.Close()
		yamlFileData, err := ioutil.ReadAll(yamlFile)
		if err != nil {
			return nil, err
		}
		jsonData, err := yaml2json.Convert(yamlFileData)
		if err != nil {
			return nil, err
		}
		var obj interface{}
		err = json.Unmarshal(jsonData, &obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}
*/
package yaml2json
