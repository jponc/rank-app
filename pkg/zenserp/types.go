package zenserp

type QueryInfo struct {
	Query string `json:"q"`
	URL   string `json:"url"`
}

type ResultItem struct {
	Position    int    `json:"position"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type QueryResult struct {
	Query      QueryInfo    `json:"query"`
	ResulItems []ResultItem `json:"organic"`
}
