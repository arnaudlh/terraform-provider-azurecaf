package testutils

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

// GetResourceDefinitions loads resource definitions from the JSON file
func GetResourceDefinitions() map[string]interface{} {
	data, err := ioutil.ReadFile(filepath.Join("..", "..", "resourcedefinition.json"))
	if err != nil {
		panic(err)
	}

	var definitions map[string]interface{}
	if err := json.Unmarshal(data, &definitions); err != nil {
		panic(err)
	}

	return definitions
}
