package reddit

type RedditAPIResponse struct {
	Data RedditPosts `json:"data"`
}

type PostData struct {
	Title string `json:"title"`
	Score int    `json:"score"`
}

type Data struct {
	PostData PostData `json:"data"`
}

type RedditPosts struct {
	Posts []Data `json:"children"`
}
