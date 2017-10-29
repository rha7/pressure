package printers

import (
	"encoding/json"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func jsonPrettyPrint(obj interface{}) string {
	v, _ := json.MarshalIndent(obj, "", "    ")
	return string(v)
}

func yamlPrettyPrint(obj interface{}) string {
	v, _ := yaml.Marshal(obj)
	return string(v)
}

func indentTextBlock(s string, prefix string) string {
	output := ""
	for _, line := range strings.Split(s, "\n") {
		output += prefix + line + "\n"
	}
	return output
}
