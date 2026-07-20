package models

type Resource struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Region       string   `json:"region"`
	ARN          string   `json:"arn"`
	Dependencies []string `json:"dependencies"`
	Tags         map[string]string `json:"tags,omitempty"`
}
