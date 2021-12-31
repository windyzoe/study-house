package modal

type BuildingRsp struct {
	Name          string
	Region        string
	District      string
	Time          int64
	BuildingCount int64
	HouseCount    int64
	Decription    string
	Alias         string
	Id            int64
	Schools       string
}

type BuildingListReq struct {
	PageNumber int64 `json:"pageNumber" form:"pageNumber"`
	PageSize   int64 `json:"pageSize" form:"pageSize"`
}
