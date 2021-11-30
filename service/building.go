package service

import (
	"encoding/json"
	"log"
)

type Building struct {
	Name          string
	Region        string
	District      string
	Time          int64
	BuildingCount int64
	HouseCount    int64
	Decription    string
	Alias         string
	Id            int64
}

func UpdateBuilding(building Building) {
	data, err := json.Marshal(building)
	if err != nil {
		log.Fatal(err)
	}
	var mapResult map[string]interface{}
	if err := json.Unmarshal(data, &mapResult); err != nil {
		log.Fatal(err)
	}
	mapValues := []map[string]interface{}{mapResult}
	UpdateWithCommit("Building", mapValues)
}
