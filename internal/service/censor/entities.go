package censor

type CheckReq struct {
	Text string `json:"text"`
}

type CheckResp struct {
	Status bool `json:"status"`
}
