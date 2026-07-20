package models

type ServiceStatus struct {
	ServiceName string         `json:"service_name"`
	Counts      map[string]int `json:"counts"`
}

type StatusReport struct {
	ScanTime string          `json:"scan_time"`
	Services []ServiceStatus `json:"services"`
}
