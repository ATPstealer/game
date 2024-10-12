package controllers

type JSONResult struct {
	Code   int         `json:"code" `
	Text   string      `json:"text"`
	Data   interface{} `json:"data"`
	Values interface{} `json:"values"`
} //@name jsonResult
