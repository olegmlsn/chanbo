package chabot

import (
	"fmt"
	"os"
)

type Chabot struct {
	BotApiKey   string
	ChannelName string
}

func New() *Chabot {
	return &Chabot{
		BotApiKey:   os.Getenv("BOT_API_KEY"),
		ChannelName: os.Getenv("CHANNEL_NAME"),
	}
}

func (c Chabot) Bla() {
	fmt.Println(c.BotApiKey, c.ChannelName)
}
