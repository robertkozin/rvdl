package redditbot

import "strings"



func (bot *RedditBot) ProcessMentions(unreads []Message) (reads []string) {

	var mentions []Message
	var postIds strings.Builder

	for _, m := range unreads {
		if !(m.WasComment && m.Type == "username_mention") {
			continue
		}

		contextParts := strings.Split(m.Context, "/")
		if len(contextParts) < 5 {
			continue
		}

		m.PostId = "t3_" + contextParts[4]

		postIds.WriteString("t3_" + contextParts[4])
		postIds.WriteByte(',')


		mentions = append(mentions, m)
	}

	var posts []Post
	err := bot.private.GetListing("/api/info", V{"id": postIds.String()}, &posts)
	if err != nil {
		return
	}

	for _, p := range posts {
		if isVideo(&p) {
			continue
		}

		if p.Hidden {
			// send pm
		} else {

		}

		//reply
	}
}

func (bot *RedditBot) ReplyToMention(m *Message, p *Post) {
	err := bot.Reply(m.Name, "xd")
	if err != nil {

	}
}

func (bot *RedditBot) MessageToMention(m *Message, p *Post) {
	err := bot.Compose(m.Author, "nodnd", "lmaoo")
	if err != nil {

	}
}

func getPermalink(p *Post) string {
	return "https://www.rvdl.com" + strings.TrimSuffix(p.Permalink, "/") + ".mp4"
}

func isVideo(p *Post) bool {
	if p.Media.RedditVideo.DashUrl != "" {
		return true
	} else if len(p.CrossPostParentList) > 0 && p.CrossPostParentList[0].Media.RedditVideo.DashUrl != "" {
		return true
	} else if p.Preview.RedditVideoPreview.DashUrl != "" {
		return true
	} else if len(p.Preview.Images) > 0 && p.Preview.Images[0].Variants.Mp4.Source.URL != "" {
		return true
	} else {
		return false
	}
}