package models

import "time"

type SLARequest struct {
	SLA       int64     `json:"sla" form:"sla" binding:"required" `
	StartTime time.Time `json:"start_time" form:"start_time"  binding:"required"`
}

type SLAResponse struct {
	SLA50  string `json:"sla_50"`
	SLA75  string `json:"sla_75"`
	SLA100 string `json:"sla_100"`
}
