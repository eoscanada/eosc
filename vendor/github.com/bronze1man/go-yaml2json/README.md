[![Codeship Status](http://img.shields.io/codeship/34b974b0-6dfa-0132-51b4-66f2bf861e14/master.svg?style=flat-square)](https://codeship.com/projects/57533)
[![API Documentation](http://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/peter-edge/go-yaml2json)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/peter-edge/go-yaml2json/blob/master/LICENSE)

Convert YAML to JSON.

## Installation
```bash
go get -u github.com/peter-edge/go-yaml2json
```

## Import
```go
import (
    "github.com/peter-edge/go-yaml2json"
)
```

Inspired by https://github.com/bronze1man/yaml2json, a command line tool with
the same effect.

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

## Usage

#### func  Convert

```go
func Convert(input []byte) ([]byte, error)
```
