package models

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	TraceId string `json:"traceId"`
	Path    string `json:"path"`
	TS      string `json:"time"`
}