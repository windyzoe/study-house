package buildingM

import (
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/windyzoe/study-house/db"
	"github.com/windyzoe/study-house/util"
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
	mapResult := util.StructToMap(building)
	mapValues := []map[string]interface{}{mapResult}
	db.UpdateWithCommit("Building", mapValues)
}

func GetAllBuilding() []Building {
	var buildings []Building
	keys := util.StructKeys(Building{})
	mapBuildings := db.QueryAll("Building", keys, -1)
	for _, v := range mapBuildings {
		data, err := json.Marshal(v)
		if err != nil {
			log.Error().Err(err)
		}
		var building Building
		if err := json.Unmarshal(data, &building); err != nil {
			log.Error().Err(err)
		}
		buildings = append(buildings, building)
	}
	return buildings
}
func GetBuildingList(pageNumber int64, pageSize int64) (pageData db.PageData) {
	mapper := map[string]string{
		"Id":            "b.Id",
		"Name":          "b.Name",
		"Region":        "b.Region",
		"District":      "b.District",
		"Decription":    "b.Decription",
		"Time":          "b.Time",
		"Alias":         "b.Alias",
		"BuildingCount": "b.BuildingCount",
		"HouseCount":    "b.HouseCount",
		"Schools":       "group_concat(rsb.School) as Schools",
	}
	pageData = db.QueryPage(mapper, "Building b left join Rel_School_Building rsb on b.Name = rsb.Building GROUP BY b.Id", pageNumber, pageSize)
	return pageData
}
