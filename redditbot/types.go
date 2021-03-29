package redditbot

type Message struct {
	FirstMessage          *string `json:"first_message"`
	FirstMessageName      *string `json:"first_message_name"`
	Subreddit             *string `json:"subreddit"`
	Likes                 *bool `json:"likes"`
	Replies               string      `json:"replies"`
	AuthorFullname        string      `json:"author_fullname"`
	ID                    string      `json:"id"`
	Subject               string      `json:"subject"`
	AssociatedAwardingID  *string `json:"associated_awarding_id"`
	Score                 int         `json:"score"`
	Author                string      `json:"author"`
	NumComments           *int `json:"num_comments"`
	ParentID              *string `json:"parent_id"`
	SubredditNamePrefixed *string `json:"subreddit_name_prefixed"`
	New                   bool        `json:"new"`
	Type                  string      `json:"type"`
	Body                  string      `json:"body"`
	Dest                  string      `json:"dest"`
	WasComment            bool        `json:"was_comment"`
	BodyHTML              string      `json:"body_html"`
	Name                  string      `json:"name"`
	Created               float64     `json:"created"`
	CreatedUtc            float64     `json:"created_utc"`
	Context               string      `json:"context"`
	Distinguished         *string `json:"distinguished"`

	PostId	string
}

type Post struct {
	URL       string `json:"url"`
	Permalink string `json:"permalink"`
	Hidden bool `json:"hidden"`
	Media     struct {
		RedditVideo struct {
			DashUrl string `json:"dash_url"`
		} `json:"reddit_video"`
	} `json:"secure_media"`
	Preview struct {
		RedditVideoPreview struct {
			DashUrl string `json:"dash_url"`
		} `json:"reddit_video_preview"`
		Images             []struct {
			Variants struct {
				Mp4 struct {
					Source struct {
						URL string `json:"url"`
					} `jsonPermalink:"source"`
				} `json:"mp4"`
			} `json:"variants"`
		} `json:"images"`
	} `json:"preview"`
	CrossPostParentList []struct {
		Media struct {
			RedditVideo struct {
				DashUrl string `json:"dash_url"`
			}`json:"reddit_video"`
		} `json:"secure_media"`
	} `json:"crosspost_parent_list"`
}

// {"reason": "private", "message": "Forbidden", "error": 403}

// subreddit_type	"private"
