package main

import (
	"errors"
	"github.com/smallfish/simpleyaml"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

// ReadYaml reads a file and parse it as YAML thanks to github.com/smallfish/simpleyaml
func ReadYaml(fileName string) (*simpleyaml.Yaml, error) {

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Err #%v ", err)
	}

	return simpleyaml.NewYaml(file)
}

// ParseDefinitions return a list of strings containing three columns , API , Definition and property
func ParseDefinitions(definitions map[string]string, yaml *simpleyaml.Yaml) ([][]string, error) {
	rows := make([][]string, 1)
	rows[0] = []string{"API", "Data_element", "property", "description"}
	for path, definition := range definitions {
		definition = strings.ReplaceAll(definition, "#/definitions/", "")
		log.Printf("%s - %s", path, definition)
		response, err := yaml.GetPath("definitions", definition, "properties").Map()
		if err != nil {
			return nil, err
		}
		keys := make([]string, 0, len(response))
		for k := range response {
			keys = append(keys, k.(string))
		}
		sort.Strings(keys)
		for _, value := range keys {
			descriptors := response[value]
			description := ""
			for key, value := range descriptors.(map[interface{}]interface{}) {
				if key == "description" {
					description = value.(string)
					break
				}
			}

			rows = append(rows, []string{path, definition, value, description})
		}
	}
	return rows, nil
}

// ParseGetPaths return a list of paths and associated swagger definition
func ParseGetPaths(yaml *simpleyaml.Yaml) (map[string]string, error) {

	if !yaml.Get("paths").IsFound() {
		return nil, errors.New("paths not found")
	}
	paths, err := yaml.Get("paths").Map()
	if err != nil {
		return nil, err
	}
	getMap := map[string]string{}
	for path, element := range paths {
		methods := element.(map[interface{}]interface{})
		for _, method := range methods {
			details := method.(map[interface{}]interface{})
			response := details["responses"]
			if response == nil {
				continue
			}
			status := response.(map[interface{}]interface{})[200]
			if status == nil {
				status = response.(map[interface{}]interface{})["200"]
			}
			// Check for int or string
			if status == nil {
				continue
			}
			success := status.(map[interface{}]interface{})["schema"]
			if success == nil {
				continue
			}
			schema := success.(map[interface{}]interface{})
			var definition string
			if schema["$ref"] == nil && schema["items"] != nil {
				items := schema["items"].(map[interface{}]interface{})["$ref"]
				definition = items.(string)
			} else if schema["$ref"] != nil {
				definition = schema["$ref"].(string)
			} else {
				continue
			}
			getMap[path.(string)] = definition
		}
	}
	return getMap, nil
}
