package reddit

type SearchResponse struct {
	Data SearchResponseData `json:"data"`
}

type SearchResponseData struct {
	After string         `json:"after"`
	Posts []PostResponse `json:"children"`
}

type PostResponse struct {
	Data PostResponseData `json:"data"`
}

type PostResponseData struct {
	Score int64  `json:"score"`
	Url   string `json:"url"`
}
