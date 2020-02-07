package repository

//http://wiki.connect.qq.com/get_user_info
type QQConnectUserInfoResult struct {
	Ret       int    `json:"ret"`
	Msg       string `json:"msg"`
	IsLost    int    `json:"is_lost"`
	Nickname  string `json:"nickname"`
	Gender    string `json:"gender"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Year      string `json:"year"`
	Avatar30  string `json:"figureurl"`
	Avatar50  string `json:"figureurl_1"`
	Avatar100 string `json:"figureurl_2"`

	AvatarQQ40  string `json:"figureurl_qq_1"`
	AvatarQQ100 string `json:"figureurl_qq_2"`

	IsYellowVip     string `json:"is_yellow_vip"`
	Vip             string `json:"vip"`
	YellowVipLevel  string `json:"yellow_vip_level"`
	Level           string `json:"level"`
	IsYellowYearVip string `json:"is_yellow_year_vip"`
}
