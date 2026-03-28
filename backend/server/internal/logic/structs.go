package logic

type RegisterData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UnRegisterData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SuccessfulRegistartion struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

type SuccessfulLogin struct {
	Status string `json:"status"`
	Name   string `json:"name"`
	Rank   int    `json:"rank"`
}

type Leaders struct {
	Name string `json:"name"`
	Rank int    `json:"rank"`
}

type LeaderBoard struct {
	Data []Leaders `json:"data"`
}

type SuccessfulLeaderBoard struct {
	Status string      `json:"status"`
	LB     LeaderBoard `json:"lb"`
}
