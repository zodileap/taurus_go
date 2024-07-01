package maputil

import "fmt"

func InterfaceToString(fieldsInterface map[string]interface{}) map[string]string {
	fieldsString := make(map[string]string, len(fieldsInterface))
	for k, v := range fieldsInterface {
		fieldsString[k] = fmt.Sprintf("%v", v)
	}
	return fieldsString
}
