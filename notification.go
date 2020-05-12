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

	msg := new(discordgo.MessageEmbed)

	msg.Title = "Travis CI: " + packet.StatusMessage
	if packet.State == "started" {
		msg.Color = 14
	}

	return nil
}
