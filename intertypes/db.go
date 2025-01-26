package intertypes

type LoginAttemptInfo struct {
	UserAgent string `json:"username"`
	IP        string `json:"ip"`
}
