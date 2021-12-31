package schoolM

import (
	"log"

	"github.com/windyzoe/study-house/db"
	"github.com/windyzoe/study-house/util"
)

type School struct {
	Name       string
	Region     string
	District   string
	Decription string
	Size       int64
	Awesome    int64
	Id         int64
}

type Rel_School_Building struct {
	Building   string
	Decription string
	School     string
	Year       int64
	Id         int64
}

func UpdateSchool(school School) {
	mapResult := util.StructToMap(school)
	mapValues := []map[string]interface{}{mapResult}
	log.Print(mapValues)
	db.UpdateWithCommit("School", mapValues)
}

func UpdateRelSchoolBuilding(rel Rel_School_Building) {
	mapResult := util.StructToMap(rel)
	mapValues := []map[string]interface{}{mapResult}
	log.Print(mapValues)
	db.UpdateWithCommit("Rel_School_Building", mapValues)
}

func GetSchoolList(pageNumber int64, pageSize int64) (pageData db.PageData) {
	mapper := map[string]string{
		"Id":         "b.Id",
		"Name":       "b.Name",
		"Region":     "b.Region",
		"District":   "b.District",
		"Decription": "b.Decription",
		"Size":       "b.Size",
		"Awesome":    "b.Awesome",
	}
	pageData = db.QueryPage(mapper, "School b", pageNumber, pageSize)
	return pageData
}
