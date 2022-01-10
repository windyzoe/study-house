package util

import (
	"encoding/json"
	"reflect"

	"github.com/rs/zerolog/log"
)

func StructKeys(stru interface{}) []string {
	var list []string
	getType := reflect.TypeOf(stru)
	for i := 0; i < getType.NumField(); i++ {
		fieldType := getType.Field(i)
		list = append(list, fieldType.Name)
	}
	return list
}

func StructToMap(stru interface{}) map[string]interface{} {
	data, err := json.Marshal(stru)
	if err != nil {
		log.Error().Err(err)
	}
	var mapResult map[string]interface{}
	if err := json.Unmarshal(data, &mapResult); err != nil {
		log.Error().Err(err)
	}
	return mapResult
}

func UniqStrings(array []string) []string {
	set := make(map[string]string) // New empty set
	for _, v := range array {
		set[v] = "1"
	}
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	return keys
}
