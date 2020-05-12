package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type NotificationBase struct {
}

// Notification subsystem
type Notification struct {
	discord *Discord
}

func (n *Notification) Init(discord *Discord) error {
	log.Infof("Initializing Notification Subsystem")
	if discord == nil {
		return fmt.Errorf("nil discord")
	}
	n.discord = discord

	return nil
}

func (n *Notification) Travis(packet *TravisPacket) error {
	if packet == nil {
		return fmt.Errorf("nil travis packet")
	}

	log.Debugf("Handling travis notification")
	msg := new(discordgo.MessageEmbed)

	log.Tracef("Travis Packet Status: %d", packet.Status)
	log.Tracef("Travis Packet State: %s", packet.State)
	log.Tracef("Travis Packet Status Message: %s", packet.StatusMessage)
	msg.Title = "Travis CI: " + packet.StatusMessage
	if packet.Status == 1 {
		if packet.State == "started" {
			msg.Color = 0xedfd00
		} else {
			msg.Color = 0xff0900
		}
	} else {
		msg.Color = 0x009b3a
	}

	msg.Author = &discordgo.MessageEmbedAuthor{
		URL:  packet.BuildURL,
		Name: packet.AuthorName,
		//IconURL: "https://travis-ci.com/images/logos/TravisCI-Mascot-blue.png",
	}

	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Started By",
		Value: packet.Type,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Message",
		Value: packet.Message,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Repository",
		Value: packet.Repository.Name,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Repository URL",
		Value: packet.Repository.URL,
	})

	n.discord.sendEmbed(n.discord.EventChannel, msg)

	return nil
}
