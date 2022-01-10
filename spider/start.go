package spider

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/gocolly/colly/v2"
	buildingM "github.com/windyzoe/study-house/modules/building"
	houseM "github.com/windyzoe/study-house/modules/house"
)

// 爬虫
func Start() {
	// houseCh := make(chan int)
	// go getHousePageCount(houseCh)
	// houseSpider(houseCh, "杨浦")

	// buildingCh := make(chan int)
	// go getBuildingPageCount(buildingCh)
	// buildingSpider(buildingCh)

	// 映射脚本
	// getSchool()
	// service.GetBuildingList()
}

func getHousePageCount(ch chan int) {
	c := colly.NewCollector(
		colly.AllowedDomains("sh.lianjia.com"),
	)
	c.OnHTML(".house-lst-page-box", func(e *colly.HTMLElement) {
		s := e.Attr("page-data")
		var mapResult map[string]int
		if err := json.Unmarshal([]byte(s), &mapResult); err != nil {
			log.Printf("Error  %v\n", err)
		}
		log.Info().Msg(`开始发送`)
		log.Info().Msgf("%s", mapResult[`totalPage`])
		ch <- mapResult[`totalPage`]
		log.Info().Msg(`发送完毕`)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error %s: %v\n", r.Request.URL, err)
	})

	c.Visit("https://sh.lianjia.com/ershoufang/yangpu/")
}

func houseSpider(ch chan int, district string) {
	pageCount := <-ch
	listCollector := colly.NewCollector(
		colly.AllowedDomains("sh.lianjia.com"),
	)
	listCollector.OnHTML(".info", func(e *colly.HTMLElement) {
		var house houseM.House
		houseInfos := strings.Split(e.ChildText(".houseInfo"), " | ")
		positions := strings.Split(e.ChildText(".positionInfo a"), " ")
		house.Name = e.ChildText(".title a")
		house.Building = positions[0]
		house.Region = positions[1]
		house.District = district
		house.Area, _ = strconv.ParseFloat(strings.Split(houseInfos[1], "平")[0], 10)
		house.Time, _ = strconv.ParseInt(strings.Split(houseInfos[5], "年")[0], 10, 0)
		floorCountMatchs := regexp.MustCompile(`\d+`).FindAllString(houseInfos[4], -1)
		house.FloorCount, _ = strconv.ParseInt(floorCountMatchs[0], 10, 0)
		unitPriceMatchs := regexp.MustCompile(`\d+`).FindAllString(e.ChildText(".unitPrice span"), -1)
		house.UnitPrice, _ = strconv.ParseInt(unitPriceMatchs[0]+unitPriceMatchs[1], 10, 0)
		house.TotalPrice, _ = strconv.ParseInt(e.ChildText(".totalPrice span"), 10, 0)
		log.Printf("%#v\n", house)
		houseM.UpdateHouse(house)
	})
	listCollector.OnError(func(r *colly.Response, err error) {
		log.Printf("Error %s: %v\n", r.Request.URL, err)
	})
	for i := 1; i <= pageCount; i++ {
		listCollector.Visit("https://sh.lianjia.com/ershoufang/yangpu/pg" + strconv.FormatInt(int64(i), 10) + "/")
	}
}

func getBuildingPageCount(ch chan int) {
	c := colly.NewCollector(
		colly.AllowedDomains("sh.lianjia.com"),
	)
	c.OnHTML(".house-lst-page-box", func(e *colly.HTMLElement) {
		s := e.Attr("page-data")
		var mapResult map[string]int
		if err := json.Unmarshal([]byte(s), &mapResult); err != nil {
			log.Printf("Error  %v\n", err)
		}
		log.Info().Msg(`开始发送`)
		log.Info().Msgf("%s", mapResult[`totalPage`])
		ch <- mapResult[`totalPage`]
		log.Info().Msg(`发送完毕`)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error %s: %v\n", r.Request.URL, err)
	})

	c.Visit("https://sh.lianjia.com/xiaoqu/yangpu/")
}

func buildingSpider(ch chan int) {
	pageCount := <-ch
	listCollector := colly.NewCollector(
		colly.AllowedDomains("sh.lianjia.com"),
	)
	buildingCollector := colly.NewCollector(
		colly.AllowedDomains("sh.lianjia.com"),
	)
	buildingCollector.OnHTML("body", func(e *colly.HTMLElement) {
		var building buildingM.Building
		aliasInfos := strings.Split(e.ChildText(".detailDesc"), ")")
		buidingInfos := e.ChildTexts(".xiaoquDetailbreadCrumbs .fl a")
		infos := e.ChildTexts(".xiaoquInfoContent")
		building.Name = e.ChildText(".detailTitle")
		building.Alias = aliasInfos[1]
		building.Region = strings.Split(buidingInfos[3], "小区")[0]
		building.District = strings.Split(buidingInfos[2], "小区")[0]
		building.Time, _ = strconv.ParseInt(strings.Split(infos[0], "年")[0], 10, 0)
		building.BuildingCount, _ = strconv.ParseInt(strings.Split(infos[5], "栋")[0], 10, 0)
		building.HouseCount, _ = strconv.ParseInt(strings.Split(infos[6], "户")[0], 10, 0)
		log.Printf("%#v\n", building)
		buildingM.UpdateBuilding(building)
	})
	listCollector.OnHTML(".title a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		buildingCollector.Visit(href)
	})
	listCollector.OnError(func(r *colly.Response, err error) {
		log.Printf("Error %s: %v\n", r.Request.URL, err)
	})
	for i := 1; i <= pageCount; i++ {
		listCollector.Visit("https://sh.lianjia.com/xiaoqu/yangpu/pg" + strconv.FormatInt(int64(i), 10) + "/")
	}
}
