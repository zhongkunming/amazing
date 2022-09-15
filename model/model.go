package model

type CheckInResult struct {
	Ret           int               `json:"ret"`
	Msg           string            `json:"msg"`
	UnFlowTraffic int64             `json:"unflowtraffic"`
	Traffic       string            `json:"traffic"`
	TrafficInfo   map[string]string `json:"trafficInfo"`
}

type LoginResult struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}
