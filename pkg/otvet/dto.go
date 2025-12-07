package otvet

// CreatePostRequest represents the request body for creating a post
type CreatePostRequest struct {
	Content   *Content `json:"content"`
	ID        *int     `json:"id,omitempty"`
	Poll      *Poll    `json:"poll,omitempty"`
	Spaces    []Space  `json:"spaces,omitempty"`
	Tags      []Tag    `json:"tags,omitempty"`
	Title     string   `json:"title"`
	TopicType int      `json:"topic_type"`
	Version   int      `json:"version"`
	VisibleTo int      `json:"visible_to"`
	AuthorID  *int     `json:"author_id,omitempty"`
}

// Content represents the content structure
type Content struct {
	Type    string        `json:"type"`
	Content []ContentNode `json:"content"`
}

// ContentNode represents a node in the content tree
type ContentNode struct {
	Attrs   map[string]interface{} `json:"attrs,omitempty"`
	Content []interface{}          `json:"content,omitempty"`
	Marks   []Mark                 `json:"marks,omitempty"`
	Text    string                 `json:"text,omitempty"`
	Type    string                 `json:"type"`
}

// Mark represents a mark in the content
type Mark struct {
	Attrs map[string]interface{} `json:"attrs,omitempty"`
	Type  string                 `json:"type"`
}

// Poll represents a poll structure
type Poll struct {
	ID       int          `json:"id"`
	Multiple bool         `json:"multiple"`
	Polls    []PollOption `json:"polls"`
	Quiz     bool         `json:"quiz"`
	Title    string       `json:"title"`
}

// PollOption represents a poll option
type PollOption struct {
	Correct            bool   `json:"correct"`
	ID                 int    `json:"id"`
	Title              string `json:"title"`
	VotedByCurrentUser bool   `json:"voted_by_current_user"`
	Votes              int    `json:"votes"`
}

// Space represents a space structure
type Space struct {
	Adult         bool           `json:"adult"`
	Banner        string         `json:"banner,omitempty"`
	Description   string         `json:"description,omitempty"`
	Icon          string         `json:"icon,omitempty"`
	ID            int            `json:"id"`
	IsPrime       bool           `json:"is_prime"`
	IsSubscribed  bool           `json:"is_subscribed"`
	OrderSpace    int            `json:"order_space,omitempty"`
	Path          string         `json:"path,omitempty"`
	SpaceCounters *SpaceCounters `json:"space_counters,omitempty"`
	Title         string         `json:"title,omitempty"`
}

// SpaceCounters represents space counters
type SpaceCounters struct {
	SubscriptionCount int `json:"subscription_count"`
	TopicCount        int `json:"topic_count"`
}

// Tag represents a tag structure
type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
}

// CreatePostResponse represents the response from creating a post
type CreatePostResponse struct {
	Result *PostResult `json:"result"`
}

// PostResult represents the post data in the response
type PostResult struct {
	ID               int64             `json:"id"`
	Title            string            `json:"title"`
	Content          *Content          `json:"content"`
	Author           *Author           `json:"author"`
	RepliesCount     int               `json:"replies_count"`
	RepliesViewCount int               `json:"replies_view_count"`
	Tags             *[]Tag            `json:"tags"`
	Spaces           []Space           `json:"spaces"`
	TopicType        int               `json:"topic_type"`
	Poll             *Poll             `json:"poll,omitempty"`
	ReactionCounter  []ReactionCounter `json:"reaction_counter"`
	CreatedAt        string            `json:"created_at"`
	UpdateAt         string            `json:"update_at"`
	IsBookmarked     bool              `json:"isBookmarked"`
	VisibleTo        int               `json:"visible_to"`
}

// Author represents the author information
type Author struct {
	ID         int64  `json:"id"`
	Nick       string `json:"nick"`
	Avatar     string `json:"avatar"`
	Username   string `json:"username"`
	Level      int    `json:"level"`
	CreatedAt  string `json:"created_at"`
	UserStatus int    `json:"user_status"`
}

// ReactionCounter represents a reaction counter
type ReactionCounter struct {
	// Add fields as needed when API documentation is available
	// For now, keeping it flexible
}

// PredictTagsSpacesRequest represents the request for tags/spaces prediction
type PredictTagsSpacesRequest struct {
	Data []PredictDataItem `json:"data"`
}

// PredictDataItem represents a single item in the prediction request
type PredictDataItem struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// PredictTagsSpacesResponse represents the response from tags/spaces prediction
type PredictTagsSpacesResponse []PredictResult

// PredictResult represents a single prediction result
type PredictResult struct {
	ID       int      `json:"id"`
	Tags     []string `json:"tags"`
	Spaces   []int    `json:"spaces"`
	MajorCat int      `json:"major_cat"`
}
