package utils

// ContainsString returns true if slice contains targeted string
func ContainsString(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}
