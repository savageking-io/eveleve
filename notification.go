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
		msg.Color = 0xff0900
	} else {
		if packet.State == "started" {
			msg.Color = 0xedfd00
		} else {
			msg.Color = 0x009b3a
		}
	}

	msg.Author = &discordgo.MessageEmbedAuthor{
		URL:     packet.BuildURL,
		Name:    packet.AuthorName,
		IconURL: "https://travis-ci.com/images/logos/TravisCI-Mascot-blue.png",
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
		Value: packet.Repository.OwnerName + "/" + packet.Repository.Name,
	})
	if packet.Repository.URL != "" {
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "Repository URL",
			Value: packet.Repository.URL,
		})
	}

	_, err := n.discord.sendEmbed(n.discord.EventChannel, msg)
	if err != nil {
		log.Errorf("Failed to send Travis Notification: %s", err.Error())
		return err
	}
	return nil
}

func (n *Notification) GitHub(e *GitHubEvent) error {
	switch e.event {
	case CommitComment:
		return n.githubCommitComment(e)
	case Fork:
	case Issue:
		return n.githubIssue(e)
	}

	return nil
}

func (n *Notification) githubCommitComment(e *GitHubEvent) error {
	//msg := new(discordgo.MessageEmbed)

	return nil
}

func (n *Notification) githubIssue(e *GitHubEvent) error {
	msg := new(discordgo.MessageEmbed)
	msg.Color = 0x2b1c39

	switch e.commitComment.Action {
	case "opened":
		msg.Title = fmt.Sprintf("New Issue %d has been created", e.issue.Issue.Number)
	case "edited":
		msg.Title = fmt.Sprintf("Issue %d has been updated", e.issue.Issue.Number)
	case "deleted":
		msg.Title = fmt.Sprintf("Issue %d has been deleted", e.issue.Issue.Number)
	case "closed":
		msg.Title = fmt.Sprintf("Issue %d has been closed", e.issue.Issue.Number)
	}

	msg.Description = "**" + e.issue.Issue.Title + "**\n"
	body := e.issue.Issue.Body
	if len(body) > 365 {
		body = body[:365] + "..."
	}
	msg.Description += body

	labels := ""

	for _, label := range e.issue.Issue.Labels {
		labels += label.Name + " "
	}
	if labels != "" {
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "Labels",
			Value: labels,
		})
	}

	assignees := ""
	for _, assignee := range e.issue.Issue.Assignees {
		assignees += assignee.Login + " "
	}
	if assignees != "" {
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "Assignees",
			Value: assignees,
		})
	}
	if e.issue.Issue.Milestone != nil {
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "Milestone",
			Value: e.issue.Issue.Milestone.Title,
		})
	}
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "State",
		Value: e.issue.Issue.Milestone.State,
	})

	msg.Author = &discordgo.MessageEmbedAuthor{
		Name:    e.issue.Issue.User.Login,
		IconURL: e.issue.Issue.User.AvatarURL,
		URL:     e.issue.Issue.User.URL,
	}

	n.discord.sendEmbed(n.discord.EventChannel, msg)

	return nil
}
