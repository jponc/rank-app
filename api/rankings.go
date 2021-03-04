package api

type SayHelloRequest struct {
	Name string `json:"name"`
}

type SayHelloResponse struct {
	Message string `json:"message"`
}

type ProcessKeywordMessage struct {
	Keyword      string `json:"keyword"`
	SearchEngine string `json:"searchEngine"`
	Device       string `json:"device"`
	Count        int    `json:"count"`
}
