package spider

// 杨浦的数据解析
import (
	"encoding/csv"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	buildingM "github.com/windyzoe/study-house/modules/building"
	schoolM "github.com/windyzoe/study-house/modules/school"
	"github.com/windyzoe/study-house/util"
)

// 解析school building的映射csv,更新到数据库,启动一起即可
func getSchool() {
	opencast, err := os.Open("./yangpu2021.csv")
	if err != nil {
		log.Error().Err(err)
	}
	defer opencast.Close()
	//创建csv读取接口实例
	ReadCsv := csv.NewReader(opencast)
	//读取所有内容
	ReadAll, err := ReadCsv.ReadAll()
	if err != nil {
		log.Error().Err(err)
	}
	length := len(ReadAll)

	set := make(map[string]string) // New empty set
	for i := 0; i < length; i++ {
		set[ReadAll[i][0]] = ReadAll[i][1]
	}
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	for _, v := range keys {
		var school schoolM.School
		school.Size = 1
		school.Name = v
		school.District = "杨浦"
		// service.UpdateSchool(school)
	}
	allBuildings := buildingM.GetAllBuilding()
	rels := computeRelSchoolBuilding(ReadAll, allBuildings)
	for _, rel := range rels {
		schoolM.UpdateRelSchoolBuilding(rel)
	}
}

func computeRelSchoolBuilding(schools [][]string, buildings []buildingM.Building) []schoolM.Rel_School_Building {
	mapResult := make(map[string]string)
	var rels []schoolM.Rel_School_Building
	for _, b := range buildings {
		alias := strings.Split(b.Alias, ",")
		alias = append(alias, b.Name)
		uniqAlias := util.UniqStrings(alias)
		isHaveRel := false
		// 对每一个alias匹配
		for _, ali := range uniqAlias {
			// 匹配学校
			for _, line := range schools {
				schoolName := line[0]
				buildingLine := line[1]
				check := checkIsHaveRel(ali, buildingLine)
				if check {
					isHaveRel = true
					// log.Printf("----buildname: %s= %s , ----REL %s : %s\n", b.Name, schoolName, ali, buildingLine)
					_, ok := mapResult[b.Name+","+schoolName]
					// 去掉小学小区都重复的
					if !ok {
						mapResult[b.Name+","+schoolName] = schoolName + "," + ali + "," + buildingLine
						rel := schoolM.Rel_School_Building{
							School:     schoolName,
							Building:   b.Name,
							Year:       2021,
							Decription: "别名:" + ali + ",官方映射:" + buildingLine,
						}
						rels = append(rels, rel)
					}
				}
			}
		}
		if !isHaveRel {
			// log.Printf("未找到%#v\n", b)
		}
	}
	return rels
}

func checkIsHaveRel(addressAlias string, buildingLine string) bool {
	rangs := regexp.MustCompile(`\d+\-\d+`).FindAllString(buildingLine, -1)
	enmus := regexp.MustCompile(`、`).FindAllString(buildingLine, -1)
	excepts := regexp.MustCompile(`\(除`).FindAllString(buildingLine, -1)
	afters := regexp.MustCompile(`以后`).FindAllString(buildingLine, -1)
	//xx路
	if len(rangs) == 0 && len(enmus) == 0 && len(excepts) == 0 && len(afters) == 0 {
		return singleCompute(addressAlias, buildingLine)
	}
	//xx路1-100
	if len(rangs) == 1 && len(enmus) == 0 && len(excepts) == 0 && len(afters) == 0 {
		return rangeCompute(addressAlias, buildingLine)
	}
	//xx路100号及以后
	if len(rangs) == 0 && len(enmus) == 0 && len(excepts) == 0 && len(afters) == 1 {
		return afterCompute(addressAlias, buildingLine)
	}
	// 带除了
	if len(excepts) > 0 {
		lineTrue := strings.Split(buildingLine, "(除")[0]
		lineExcept := strings.Split(buildingLine, "(除")[1]
		subRanges := regexp.MustCompile(`\d+\-\d+`).FindAllString(lineTrue, -1)
		subAfters := regexp.MustCompile(`以后`).FindAllString(lineTrue, -1)
		// xx路
		if len(subRanges) == 0 && len(subAfters) == 0 {
			inCompute := singleCompute(addressAlias, lineTrue)
			if !inCompute {
				return false
			}
			inExcept := checkInExcept(lineExcept, addressAlias)
			if inExcept {
				// log.Info().Msg(addressAlias, buildingLine)
			}
			return !inExcept
		}
		// xx路1-100
		if len(subRanges) == 1 && len(subAfters) == 0 {
			inCompute := rangeCompute(addressAlias, lineTrue)
			if !inCompute {
				return false
			}
			inExcept := checkInExcept(lineExcept, addressAlias)
			if inExcept {
				// log.Info().Msg(addressAlias, buildingLine)
			}
			return !inExcept
		}
		// xx路xx号以后
		if len(subRanges) == 0 && len(subAfters) == 1 {
			inCompute := afterCompute(addressAlias, lineTrue)
			if !inCompute {
				return false
			}
			inExcept := checkInExcept(lineExcept, addressAlias)
			if inExcept {
				// log.Info().Msg(addressAlias, buildingLine)
			}
			return !inExcept
		}

	}
	return false
}

func checkInExcept(lineExcept string, addressAlias string) bool {
	exceptRanges := regexp.MustCompile(`\d+\-\d+`).FindAllString(lineExcept, -1)
	allNumbers := regexp.MustCompile(`\d+`).FindAllString(lineExcept, -1)
	addressNumbers := regexp.MustCompile(`\d+`).FindAllString(addressAlias, -1)
	// 这个没数字就没意义了,算true
	if addressNumbers == nil {
		return false
	}
	addressNumberString := addressNumbers[0]
	// 除了里 不带 -
	if len(exceptRanges) == 0 {
		for _, v := range allNumbers {
			if v == addressNumberString {
				return true
			}
		}
	}
	if len(exceptRanges) == 1 {
		numberStrs := strings.Split(exceptRanges[0], "-")
		start, _ := strconv.ParseInt(numberStrs[0], 10, 64)
		end, _ := strconv.ParseInt(numberStrs[1], 10, 64)
		addressNumber, _ := strconv.ParseInt(addressNumberString, 10, 64)
		if (addressNumber > start && addressNumber < end) || addressNumber == start || addressNumber == end {
			return true
		}
	}
	return false
}

func singleCompute(addressAlias string, buildingLine string) bool {
	if strings.Index(addressAlias, buildingLine) != -1 {
		return true
	}
	return false
}

func afterCompute(addressAlias string, buildingLine string) bool {
	// 路
	addressStrs := regexp.MustCompile("^[\u4e00-\u9fa5]+").FindAllString(addressAlias, -1)
	addressBuildingStrs := regexp.MustCompile("^[\u4e00-\u9fa5]+").FindAllString(buildingLine, -1)
	if addressStrs == nil || addressBuildingStrs == nil || addressStrs[0] != addressBuildingStrs[0] {
		return false
	}
	addressNumberString := regexp.MustCompile(`\d+`).FindAllString(addressAlias, -1)
	if len(addressNumberString) == 0 {
		return false
	}
	addressNumber, _ := strconv.ParseInt(addressNumberString[0], 10, 64)
	afterNumberString := regexp.MustCompile(`\d+`).FindAllString(buildingLine, -1)
	if len(afterNumberString) == 0 {
		return false
	}
	afterNumber, _ := strconv.ParseInt(afterNumberString[0], 10, 64)
	// log.Info().Msg(addressNumber, afterNumber, addressAlias, buildingLine)
	if addressNumber >= afterNumber {
		return true
	}

	return false
}

func rangeCompute(addressAlias string, buildingLine string) bool {
	rangs := regexp.MustCompile(`\d+\-\d+`).FindAllString(buildingLine, -1)
	// 路
	addressStrs := regexp.MustCompile("^[\u4e00-\u9fa5]+").FindAllString(addressAlias, -1)
	addressBuildingStrs := regexp.MustCompile("^[\u4e00-\u9fa5]+").FindAllString(buildingLine, -1)
	if addressStrs == nil || addressBuildingStrs == nil || addressStrs[0] != addressBuildingStrs[0] {
		return false
	}
	// 数字
	addressNumberString := regexp.MustCompile(`\d+`).FindAllString(addressAlias, -1)
	if len(addressNumberString) == 0 {
		return false
	}
	addressNumber, _ := strconv.ParseInt(addressNumberString[0], 10, 64)
	numberStrs := strings.Split(rangs[0], "-")
	start, _ := strconv.ParseInt(numberStrs[0], 10, 64)
	end, _ := strconv.ParseInt(numberStrs[1], 10, 64)
	if (addressNumber > start && addressNumber < end) || addressNumber == start || addressNumber == end {
		return true
	}
	return false
}
