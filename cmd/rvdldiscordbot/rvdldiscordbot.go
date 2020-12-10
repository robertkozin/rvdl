package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"rvdl/internal/discordbot"
	"rvdl/pkg/util"
	"syscall"
)

var DiscordToken = util.EnvString("RVDL_DISCORD_TOKEN", "")

func main() {


	go func() {
		if _, err := discordbot.SetupDiscord(DiscordToken); err != nil {
			log.Fatalf("discord: %s\n", err)
		}
	}()

	fmt.Println("Discord Bot Started")

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done

	fmt.Println("Discord Bot Stopped")


	err := discordbot.TeardownDiscord()
	if err != nil {
		log.Fatalln(err)
	}
}
