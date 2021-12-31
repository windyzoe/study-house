package houseM

import (
	"github.com/windyzoe/study-house/db"
	"github.com/windyzoe/study-house/util"
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
	mapResult := util.StructToMap(house)
	houseMap := []map[string]interface{}{mapResult}
	db.UpdateWithCommit("House", houseMap)
}

func GetHouseList(pageNumber int64, pageSize int64) (pageData db.PageData) {
	mapper := map[string]string{
		"Id":         "b.Id",
		"Building":   "b.Building",
		"Name":       "b.Name",
		"Region":     "b.Region",
		"District":   "b.District",
		"Time":       "b.Time",
		"TotalPrice": "b.TotalPrice",
		"Area":       "b.Area",
		"UnitPrice":  "b.UnitPrice",
	}
	pageData = db.QueryPage(mapper, "House b", pageNumber, pageSize)
	return pageData
}
