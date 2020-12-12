package discordbot

import (
	"github.com/bwmarrin/discordgo"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"rvdl/pkg/util"
	"strings"
)

var Discord *discordgo.Session
var botID = ""

var helpMessage = `
	I will reply with an mp4 link to any reddit link posted with a video
`

var redditLink = regexp.MustCompile(`\S*(?:reddit\.com|redd\.it)/\S+`)

var ifModifiedSince = "Wed, 02 Dec 2020 00:00:00 GMT"
var userAgent = "rvdl"

func SetupDiscord(secret string) (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + secret)
	if err != nil {
		return nil, err
	}

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsDirectMessages | discordgo.IntentsGuildMessages)
	discord.Identify.LargeThreshold = 50
	discord.Identify.GuildSubscriptions = false
	discord.Identify.Presence = discordgo.GatewayStatusUpdate{
		Since: 0,
		Game: discordgo.Activity{
			Name: "@me",
			Type: discordgo.ActivityTypeCustom,
		},
		Status: string(discordgo.StatusOnline),
		AFK:    false,
	}

	discord.SyncEvents = true
	discord.StateEnabled = false

	discord.AddHandler(ready)
	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		return nil, err
	}

	Discord = discord

	return discord, nil

	//TODO: return close channel, or actually just make a close method
}

func TeardownDiscord() error {
	return Discord.Close()
}

func ready(s *discordgo.Session, ready *discordgo.Ready) {
	botID = ready.User.ID
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from self
	if m.Author.ID == botID {
		return
	}

	// Reply with a help message when mentioned
	if len(m.Mentions) == 1 && m.Mentions[0].ID == botID {
		_, _ = s.ChannelMessageSend(m.ChannelID, helpMessage)
		return
	}

	// Handle reddit links
	if !(strings.Contains(m.Content, "reddit.com/") || strings.Contains(m.Content, "redd.it/")) {
		return
	}

	link := redditLink.FindString(m.Content)

	if strings.Contains(link, ".png") || strings.Contains(link, ".gif") || strings.Contains(link, ".jpg") {
		return
	}

	req, err := http.NewRequest("GET", "https://www.rvdl.com/"+link, nil)
	if err != nil {
		return
	}

	req.Header.Set("If-Modified-Since", ifModifiedSince)
	req.Header.Set("User-Agent", userAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	_, _ = io.Copy(ioutil.Discard, res.Body)
	_ = res.Body.Close()

	if (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotModified) || res.Header.Get("Content-Type") != "video/mp4" {
		return
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, util.UrlRawString(res.Request.URL))
}
