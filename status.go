package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"time"
)

// Status updates status post on Discord with actual data
// after some periods of time
type Status struct {
	ChannelID string
	MessageID string

	LastUpdate time.Time
	StartTime  time.Time
	Discord    *Discord
}

func (s *Status) Init(discord *Discord) error {
	if discord == nil {
		return fmt.Errorf("discord is nil")
	}
	s.Discord = discord
	s.StartTime = time.Now()

	s.ClearStatusMessages()

	s.MessageID = ""

	return nil
}

func (s *Status) Run() error {

	for {
		if time.Since(s.LastUpdate) > time.Duration(time.Second*30) {
			s.LastUpdate = time.Now()
			if err := s.UpdateStatus(); err != nil {
				log.Errorf("Failed to update status: %s", err.Error())
			}
		}

		time.Sleep(time.Millisecond * 100)
	}

	return nil
}

func (s *Status) UpdateStatus() error {
	if s.Discord == nil {
		return fmt.Errorf("discord is nil")
	}

	msg := new(discordgo.MessageEmbed)
	msg.Title = "Status"
	msg.Author = new(discordgo.MessageEmbedAuthor)
	msg.Author.Name = "WatchDog"
	msg.Author.IconURL = "https://savageking.io/img/savage-king-games.png"
	msg.Author.URL = "https://savageking.io"

	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Master Server Uptime",
		Value: s.GetUptime(),
	})

	if s.MessageID == "" {
		newMsg, err := s.Discord.sendEmbed(s.Discord.StatusChannel, msg)
		if err != nil {
			return err
		}
		s.MessageID = newMsg.ID
		return nil
	}

	_, err := s.Discord.editEmbed(s.Discord.StatusChannel, s.MessageID, msg)
	return err
}

func (s *Status) ClearStatusMessages() error {
	if s.Discord == nil {
		return fmt.Errorf("discord is nil")
	}

	msg, err := s.Discord.getMessages(s.Discord.StatusChannel)
	if err != nil {
		return nil
	}

	for _, m := range msg {
		if m == "" {
			continue
		}

		s.Discord.deleteMessage(s.Discord.StatusChannel, m)
	}

	return nil
}

func (s *Status) GetUptime() string {
	return time.Since(s.StartTime).String()
}
