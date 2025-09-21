package reddit

type RedditAPIResponse struct {
	Data RedditPosts `json:"data"`
}

type PostData struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Score     int    `json:"score"`
	Subreddit string `json:"subreddit"`
}

type Data struct {
	PostData PostData `json:"data"`
}

type RedditPosts struct {
	Posts []Data `json:"children"`
}

type FullPost struct {
	Post     PostData      `json:"post"`
	Comments []CommentData `json:"comments"`
}

type CommentAPIResponse []struct {
	Data struct {
		Children []struct {
			Data CommentData `json:"data"`
		} `json:"children"`
	} `json:"data"`
}
type CommentData struct {
	Author string `json:"author"`
	Body   string `json:"body"`
	Score  int    `json:"score"`
}
type CommentMessage struct {
	PostID  string      `json:"post_id"`
	Comment CommentData `json:"comment"`
}
