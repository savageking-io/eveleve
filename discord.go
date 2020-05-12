package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Command struct {
	Cmd    string
	Params []string
}

type Discord struct {
	Token         string
	LogChannel    string
	EventChannel  string
	StatusChannel string
	Session       *discordgo.Session
	Commands      chan Command
}

func (d *Discord) Init(config DiscordConfig) error {
	var err error
	d.Commands = make(chan Command)
	log.Infof("Initializing Discord Bot")
	d.Token = config.Token
	d.LogChannel = config.LogChannel
	d.EventChannel = config.EventChannel
	d.StatusChannel = config.StatusChannel

	d.Session, err = discordgo.New("Bot " + d.Token)
	if err != nil {
		return err
	}

	d.Session.Token = "Bot " + d.Token
	d.Session.State.User, err = d.Session.User("@me")
	if err != nil {
		return fmt.Errorf("Failed to retrieve user data: %s", err)
	}

	d.Session.AddHandler(d.messageCreate)

	err = d.Session.Open()
	if err != nil {
		return fmt.Errorf("Failed to open connection to Discord: %s", err.Error())
	}

	return nil
}

func (d *Discord) messageCreate(s *discordgo.Session, msg *discordgo.MessageCreate) {
	parts := strings.Split(msg.Content, " ")
	if len(msg.Content) == 0 {
		return
	}
	if len(parts) == 0 {
		return
	}

	if parts[0][0] == '!' {
		var c Command
		c.Cmd = parts[0]
		if len(parts) > 1 {
			c.Params = parts[1:]
		}
		d.Commands <- c
	}
}

func (d *Discord) sendLog(text string) {
	d.sendMessage(text, d.LogChannel)
}

func (d *Discord) sendEvent(text string) {
	d.sendMessage(text, d.EventChannel)
}

func (d *Discord) sendMessage(text, channelID string) {
	buffer := []string{}

	// Discord can receive message up to 2000 characters long.
	if len(text) > 1999 {
		lines := strings.Split(text, "\n")
		b := ""
		for _, line := range lines {
			if len(b)+len(line)+len("\n") > 1999 {
				buffer = append(buffer, b)
				b = line + "\n"
				continue
			}
			b += line + "\n"
		}
		buffer = append(buffer, b)
	} else {
		buffer = append(buffer, text)
	}

	for _, str := range buffer {
		_, err := d.Session.ChannelMessageSend(channelID, str)
		if err != nil {
			log.Errorf("Failed to send message to %s: %s", channelID, err.Error())
		}
	}
}

func (d *Discord) sendEmbed(channelID string, data *discordgo.MessageEmbed) (*discordgo.Message, error) {
	log.Tracef("Sending Message Embed to %s: %+v", channelID, data)
	return d.Session.ChannelMessageSendEmbed(channelID, data)
}

func (d *Discord) editEmbed(channelID, msgID string, data *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return d.Session.ChannelMessageEditEmbed(channelID, msgID, data)
}

func (d *Discord) getMessages(channelID string) ([]string, error) {
	messages, err := d.Session.ChannelMessages(channelID, 10, "", "", "")
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, msg := range messages {
		result = append(result, msg.ID)
	}
	return result, nil
}

func (d *Discord) deleteMessage(channelID, messageID string) error {
	return d.Session.ChannelMessageDelete(channelID, messageID)
}
