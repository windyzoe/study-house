package service

import (
	"encoding/json"
	"log"
)

// 奇怪,结构体字段必须走大写
type House struct {
	Name       string
	Building   string
	Region     string
	District   string
	Time       int64
	TotalPrice int64
	UnitPrice  int64
	Area       float64
	FloorCount int64
	Id         int64
}

func UpdateHouse(house House) {
	data, err := json.Marshal(house)
	if err != nil {
		log.Fatal(err)
	}
	var mapResult map[string]interface{}
	if err := json.Unmarshal(data, &mapResult); err != nil {
		log.Fatal(err)
	}
	houseMap := []map[string]interface{}{mapResult}
	UpdateWithCommit("House", houseMap)
}
