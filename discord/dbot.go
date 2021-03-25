package discord

import (
	discord "github.com/bwmarrin/discordgo"
	"github.com/robertkozin/rvdl/pkg/util"
	"github.com/segmentio/encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type GuildId = string
type UserId = string

type DiscordBot struct {
	Id string
	Session *discord.Session

	owners map[GuildId]UserId

	settingsFile string
	settings map[GuildId]*GuildSettings
}

type GuildSettings struct {
	WhenMentioned bool `json:"wm"`
	HideParentEmbed bool `json:"hpe"`
}

func NewRvdlDiscordBot(token string, settingsFile string) (*DiscordBot, error) {
	bot := &DiscordBot{settingsFile: settingsFile}

	err := bot.ReadSettings()
	if err != nil {
		return nil, err
	}

	session, _ := discord.New("Bot " + token)

	session.Identify.Intents = discord.MakeIntent(discord.IntentsGuildMessages)
	session.Identify.LargeThreshold = 50
	session.Identify.GuildSubscriptions = false
	session.SyncEvents = true
	session.StateEnabled = false

	session.AddHandler(bot.ready)
	session.AddHandler(bot.messageCreate)

	err = session.Open()
	if err != nil {
		return nil, err
	}

	bot.Session = session

	return bot, nil
}

func (bot *DiscordBot) ReadSettings() error {
	config, err := os.OpenFile(bot.settingsFile, os.O_CREATE | os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer config.Close()

	err = json.NewDecoder(config).Decode(&bot.settings)
	if err != nil {
		return err
	}

	return nil
}

func (bot *DiscordBot) WriteSettings() error {
	config, err := os.OpenFile(bot.settingsFile, os.O_CREATE | os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer config.Close()

	err = json.NewEncoder(config).Encode(&bot.settings)
	if err != nil {
		return err
	}

	return nil
}

func (bot *DiscordBot) Close() error {
	_ = bot.WriteSettings()
	_ = bot.Session.Close()

	return nil
}

func (bot *DiscordBot) ready(s *discord.Session, ready *discord.Ready) {
	bot.Id = ready.User.ID

	err := bot.Session.UpdateListeningStatus("@me")
	if err != nil {
		// TODO: log?
	}

	for _, g := range ready.Guilds {
		bot.owners[g.ID] = g.OwnerID
	}
}

func (bot *DiscordBot) messageCreate(s *discord.Session, m *discord.MessageCreate) {
	if m.Type != discord.MessageTypeDefault {
		return
	}

	// Ignore messages from self
	if m.Author.ID == bot.Id {
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

	go func() {
		req, err := http.NewRequest("GET", "https://www.rvdl.com/"+link, nil)
		if err != nil {
			return
		}

		req.Header.Set("User-Agent", userAgent)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}

		_, _ = io.Copy(ioutil.Discard, res.Body)
		_ = res.Body.Close()

		if res.Header.Get("Video-Found") != "?1" {
			return
		}

		_, _ = s.ChannelMessageSend(m.ChannelID, util.UrlRawString(res.Request.URL))
	}()
}
