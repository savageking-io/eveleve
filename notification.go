package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v5/github"
)

type NotificationField struct {
	Name   string `yaml:"name"`
	Value  string `yaml:"value"`
	Inline bool   `yaml:"inline"`
}

type NotificationConfig struct {
	Title       string              `yaml:"title"`
	Color       int                 `yaml:"color"`
	Description string              `yaml:"description"`
	Fields      []NotificationField `yaml:"fields"`
	Footer      struct {
		Text      string `yaml:"text"`
		Icon      string `yaml:"icon"`
		ProxyIcon string `yaml:"proxy_icon"`
	} `yaml:"footer"`
	Image struct {
		URL      string `yaml:"url"`
		ProxyURL string `yaml:"proxy_url"`
		Width    int    `yaml:"width"`
		Height   int    `yaml:"height"`
	} `yaml:"image"`
	Author struct {
		URL       string `yaml:"url"`
		Name      string `yaml:"name"`
		Icon      string `yaml:"icon"`
		ProxyIcon string `yaml:"proxy_icon"`
	} `yaml:"author"`
	Provider struct {
		URL  string `yaml:"url"`
		Name string `yaml:"name"`
	} `yaml:"provider"`
	Thumbnail struct {
		URL      string `yaml:"url"`
		ProxyURL string `yaml:"proxy_url"`
		Width    int    `yaml:"width"`
		Height   int    `yaml:"height"`
	} `yaml:"thumbnail"`
	Video struct {
		URL      string `yaml:"url"`
		ProxyURL string `yaml:"proxy_url"`
		Width    int    `yaml:"width"`
		Height   int    `yaml:"height"`
	} `yaml:"video"`
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

func (n *Notification) ParseGitHub(e *GitHubEvent) error {
	switch e.event {
	case Issue:
	}

	return nil
}

func (n *Notification) parseIssue(p *github.IssuesPayload) (map[string]string, error) {
	if p == nil {
		return nil, fmt.Errorf("nil issue")
	}
	data := make(map[string]string)

	data["{issue.action}"] = p.Action
	data["{issue.assignees}"] = ""
	for _, assignee := range p.Issue.Assignees {
		data["{issue.assignees}"] += assignee.Login + " "
	}
	data["{issue.labels}"] = ""
	for _, label := range p.Issue.Labels {
		data["{issue.labels}"] += label.Name + " "
	}
	data["{issue.milestone.title}"] = p.Issue.Milestone.Title
	data["{issue.milestone.number}"] = fmt.Sprintf("%d", p.Issue.Milestone.Number)
	data["{issue.author.name}"] = p.Issue.User.Login
	data["{issue.author.icon}"] = p.Issue.User.AvatarURL
	data["{issue.author.url}"] = p.Issue.User.URL

	return data, nil
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
	default:
		return nil
	}

	msg.Description = e.commitComment.Action + "**" + e.issue.Issue.Title + "**\n"
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
		Value: e.issue.Issue.State,
	})

	msg.Author = &discordgo.MessageEmbedAuthor{
		Name:    e.issue.Issue.User.Login,
		IconURL: e.issue.Issue.User.AvatarURL,
		URL:     e.issue.Issue.User.URL,
	}

	n.discord.sendEmbed(n.discord.EventChannel, msg)

	return nil
}

func (n *Notification) githubPush(e *GitHubEvent) error {
	msg := new(discordgo.MessageEmbed)
	msg.Color = 0x2b1c39
	msg.Title = e.push.Sender.Login + " sent new commits to " + e.push.Repository.FullName
	msg.Author = &discordgo.MessageEmbedAuthor{
		URL:     e.push.Sender.HTMLURL,
		Name:    e.push.Sender.Login,
		IconURL: e.push.Sender.AvatarURL,
	}
	for _, c := range e.push.Commits {
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  c.Message,
			Value: c.ID + " by " + c.Committer.Username,
		})
	}
	msg.Provider = &discordgo.MessageEmbedProvider{
		URL:  "https://github.com",
		Name: "GitHub",
	}

	n.discord.sendEmbed(n.discord.EventChannel, msg)

	return nil
}
