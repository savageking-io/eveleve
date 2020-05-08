package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v5/github"
	"net/http"
)

// GitHub listens for github hooks and performs actions
type GitHub struct {
	Port     uint16 // Port for webhooks
	Events   chan GitHubEvent
	Discord  *Discord
	Projects []string
}

type GitHubEventType uint8

const (
	Release GitHubEventType = iota
	Push    GitHubEventType = iota
)

type GitHubEvent struct {
	event   GitHubEventType
	release github.ReleasePayload
	push    github.PushPayload
}

//func (g *GitHub) Init(port uint16, cert, key string) error {
func (g *GitHub) Init(ghc GitHubConfig, tlsc TLSConfig) error {
	log.Infof("Preparing GitHub webhook listener at port %d", ghc.Port)
	g.Port = ghc.Port
	g.Events = make(chan GitHubEvent)

	hook, _ := github.New(github.Options.Secret(ghc.Secret))

	http.HandleFunc(ghc.URI, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.ReleaseEvent, github.PushEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				log.Infof("Received payload for a different event: %+v", err.Error())
			}
		}
		switch payload.(type) {

		case github.ReleasePayload:
			g.Release(payload.(github.ReleasePayload))
		case github.PushPayload:
			g.Push(payload.(github.PushPayload))
		}
	})
	go http.ListenAndServeTLS(fmt.Sprintf(":%d", g.Port), tlsc.Cert, tlsc.Key, nil)
	return nil
}

func (g *GitHub) SetProjects(projects []string) {
	url := "github.com/"
	g.Projects = g.Projects[:0]
	log.Infof("Setting GitHub projects")
	for _, project := range projects {
		if project[0:len(url)] == url {
			g.Projects = append(g.Projects, project[len(url):])
			log.Infof("Adding project %s", project[len(url):])
			continue
		}
		log.Infof("Ignoring not a GitHub project: %s", project)
	}
	log.Infof("%d project added in total", len(g.Projects))
}

func (g *GitHub) addNotificationSubsystem(d *Discord) {
	g.Discord = d
}

func (g *GitHub) Release(payload github.ReleasePayload) {

}

func (g *GitHub) verifyProject(name string) error {
	for _, project := range g.Projects {
		if project == name {
			return nil
		}
	}
	return fmt.Errorf("Unknown repository")
}

func (g *GitHub) Push(payload github.PushPayload) {
	// Verify it's one of our projects
	if g.verifyProject(payload.Repository.FullName) != nil {
		log.Warnf("Payload came from unverified project: %+v", payload)
		if g.Discord != nil {
			g.Discord.sendLog("Repository event from unverified project")
		}
		return
	}

	if g.Discord == nil {
		log.Errorf("Can't send discord notification: discord is nil")
		return
	}

	msg := new(discordgo.MessageEmbed)
	msg.Title = "New Push Event in " + payload.Repository.Name
	msg.Color = 14

	msg.Author = new(discordgo.MessageEmbedAuthor)
	msg.Author.Name = payload.Pusher.Name
	msg.Author.IconURL = payload.Sender.AvatarURL
	msg.Author.URL = payload.Sender.URL

	for _, c := range payload.Commits {
		f := new(discordgo.MessageEmbedField)
		f.Name = c.ID
		f.Value = c.Message
		msg.Fields = append(msg.Fields, f)
	}

	g.Discord.sendEmbed(g.Discord.EventChannel, msg)
}
