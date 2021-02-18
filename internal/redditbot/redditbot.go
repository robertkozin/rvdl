package redditbot

import (
	"github.com/robertkozin/rvdl/pkg/reddit"
	"strings"
)

type RedditBot struct {
	private reddit.Client
	rss     reddit.Client
}

type V = reddit.V

func NewRedditBot(clientId, clientSecret, username, password, userAgent string) (*RedditBot, error) {
	bot := &RedditBot{}

	client, err := reddit.NewPrivateClient(clientId, clientSecret, username, password, userAgent)
	if err != nil {
		return nil, err
	}

	bot.private = client

	bot.rss = reddit.NewRssClient(clientId, username, userAgent) // TOKEN

	return bot, nil
}

func (bot *RedditBot) Process() {
	unreads, err := bot.GetUnreads()
	if err != nil {
		return
	}

	bot.ProcessMentions(unreads)
}



func (bot *RedditBot) GetUnreads() ([]Message, error) {
	var messages []Message
	err := bot.rss.Get(
		"/message/unread",
		nil,
		&messages,
	)

	return messages, err
}

func (bot *RedditBot) GetPosts(ids []string) ([]Post, error) {
	var posts []Post
	err := bot.private.GetListing(
		"/api/info",
			V{"id": strings.Join(ids, ",")},
			posts,
		)

	return posts, err
}

func (bot *RedditBot) HidePosts(ids []string) error {
	err := bot.private.Post(
		"/api/hide",
		V{"id": strings.Join(ids, ",")},
		nil, // TODO Post handle nil reciever
		 // TODO figure out endpoint response
		)

	return err
}

func (bot *RedditBot) Reply(thingId string, text string) error {
	err := bot.private.Post(
		"/api/comment",
		V{
			"api_type": "json",
			"text" : text,
			"thing_id": thingId,
		},
		nil,
		)

	return err
}

func (bot *RedditBot) Compose(to, subject, text string) error {
	err := bot.private.Post(
		"/api/compose",
		V{
			"api_type": "json",
			"to": to,
			"subject": subject,
			"text": text,
		},
		nil,
		)

	return err
}
