package model

type RawJob struct {
	Title       string
	Company     string
	Description string
	Url         string
	Location    string
	Salary      string
	Skills      []string
	Source      string
}

type JobDetail struct {
	Description string
	Location    string
	Skills      []string
	Company     string
	Salary      string
}

type ScrapeRequest struct {
	UserID  string `json:"user_id"`
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit"`
}

type ScrapeResponse struct {
	Message          string `json:"message"`
	TotalFound       int    `json:"total_found"`
	TotalInserted    int    `json:"total_inserted"`
	TotalDuplicated  int    `json:"total_duplicated"`
	ExecutionTimeSec int    `json:"execution_time_sec"`
}

type ProxyAuth struct {
	Server string
	User   string
	Pass   string
}
