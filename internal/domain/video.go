package domain

// Video represents the core video entity returned by the application.
type Video struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	ViewCount    string `json:"viewCount"`
	LikeCount    string `json:"likeCount"`
	ChannelTitle string `json:"channelTitle"`
	PublishedAt  string `json:"publishedAt"`
}
