package rest

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
