package scraper

type ScraperRequest struct {
	UserId string `json:"user_id"`
	Limit  int    `json:"limit"`
}

type ScraperResponse struct {
	Message          string `json:"message"`
	TotalFound       int    `json:"total_found"`
	TotalInserted    int    `json:"total_inserted"`
	TotalDuplicated  int    `json:"total_duplicated"`
	ExecutionTimeSec int    `json:"execution_time_sec"`
}

type RawJob struct {
	Title       string
	Company     string
	Description string
	Url         string
	Location    string
	Skills      []string
}

type JobDetail struct {
	Description string
	Company     string
	Location    string
	Skills      []string
}
