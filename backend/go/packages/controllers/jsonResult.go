package controllers

type JSONResult struct {
	Code   int         `json:"code" example:"0" validate:"required"`
	Text   string      `json:"text"`
	Data   interface{} `json:"data"`
	Values interface{} `json:"values"`
} //@name jsonResult

type JSONResultStringArray struct {
	Code   int         `json:"code"`
	Data   []string    `json:"data"`
	Text   string      `json:"text,omitempty"`
	Values interface{} `json:"values"`
} //@name jsonResultStringArray
