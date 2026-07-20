package models

type Result struct {
	DeletedResources []Resource `json:"deleted_resources"`
	FailedResources  []Resource `json:"failed_resources"`
}
