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
	CommitComment      GitHubEventType = iota
	Fork               GitHubEventType = iota
	Issue              GitHubEventType = iota
	IssueComment       GitHubEventType = iota
	Milestone          GitHubEventType = iota
	PullRequest        GitHubEventType = iota
	PullRequestReview  GitHubEventType = iota
	PullRequestComment GitHubEventType = iota
	Push               GitHubEventType = iota
	Vulnerability      GitHubEventType = iota
	Release            GitHubEventType = iota
	Security           GitHubEventType = iota
)

type GitHubEvent struct {
	event              GitHubEventType
	commitComment      github.CommitCommentPayload
	fork               github.ForkPayload
	issue              github.IssuesPayload
	issueComment       github.IssueCommentPayload
	milestone          github.MilestonePayload
	push               github.PushPayload
	pullRequest        github.PullRequestPayload
	pullRequestReview  github.PullRequestReviewPayload
	pullRequestComment github.PullRequestReviewCommentPayload
	vulnerability      github.RepositoryVulnerabilityAlertPayload
	release            github.ReleasePayload
	security           github.SecurityAdvisoryPayload
}

//func (g *GitHub) Init(port uint16, cert, key string) error {
func (g *GitHub) Init(ghc GitHubConfig, tlsc TLSConfig) error {
	log.Infof("Preparing GitHub webhook listener at port %d", ghc.Port)
	g.Port = ghc.Port
	g.Events = make(chan GitHubEvent)

	hook, _ := github.New(github.Options.Secret(ghc.Secret))

	http.HandleFunc(ghc.URI, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.ReleaseEvent, github.PushEvent,
			github.CommitCommentEvent, github.IssuesEvent, github.IssueCommentEvent,
			github.ForkEvent, github.MilestoneEvent, github.PullRequestEvent,
			github.PullRequestReviewEvent, github.RepositoryVulnerabilityAlertEvent,
			github.SecurityAdvisoryEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				log.Infof("Received payload for a different event: %+v", err.Error())
			}
		}
		switch payload.(type) {
		case github.CommitCommentPayload:
			g.CommitComment(payload.(github.CommitCommentPayload))
		case github.ForkPayload:
			g.Fork(payload.(github.ForkPayload))
		case github.IssuesPayload:
			g.Issue(payload.(github.IssuesPayload))
		case github.IssueCommentPayload:
			g.IssueComment(payload.(github.IssueCommentPayload))
		case github.MilestonePayload:
			g.Milestone(payload.(github.MilestonePayload))
		case github.PushPayload:
			g.Push(payload.(github.PushPayload))
		case github.PullRequestPayload:
			g.PullRequest(payload.(github.PullRequestPayload))
		case github.PullRequestReviewPayload:
			g.PullRequestReview(payload.(github.PullRequestReviewPayload))
		case github.PullRequestReviewCommentPayload:
			g.PullRequestComment(payload.(github.PullRequestReviewCommentPayload))
		case github.RepositoryVulnerabilityAlertPayload:
			g.Vulnerability(payload.(github.RepositoryVulnerabilityAlertPayload))
		case github.ReleasePayload:
			g.Release(payload.(github.ReleasePayload))
		case github.SecurityAdvisoryPayload:
			g.SecurityAdvisory(payload.(github.SecurityAdvisoryPayload))
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

func (g *GitHub) CommitComment(p github.CommitCommentPayload) error {
	if g.verifyProject(p.Repository.FullName) != nil {
		log.Warnf("Payload came from unverified project: %+v", p)
		if g.Discord != nil {
			g.Discord.sendLog("Repository event from unverified project")
		}
		return fmt.Errorf("unknown repository")
	}

	event := &GitHubEvent{
		event:         CommitComment,
		commitComment: p,
	}
	g.Events <- *event
	return nil
}

func (g *GitHub) Fork(p github.ForkPayload) error {
	if g.verifyProject(p.Repository.FullName) != nil {
		log.Warnf("Payload came from unverified project: %+v", p)
		if g.Discord != nil {
			g.Discord.sendLog("Repository event from unverified project")
		}
		return fmt.Errorf("unknown repository")
	}
	event := &GitHubEvent{
		event: Fork,
		fork:  p,
	}
	g.Events <- *event
	return nil
}

func (g *GitHub) Issue(p github.IssuesPayload) error {
	if g.verifyProject(p.Repository.FullName) != nil {
		log.Warnf("Payload came from unverified project: %+v", p)
		if g.Discord != nil {
			g.Discord.sendLog("Repository event from unverified project")
		}
		return fmt.Errorf("unknown repository")
	}
	event := &GitHubEvent{
		event: Issue,
		issue: p,
	}
	g.Events <- *event
	return nil
}

func (g *GitHub) IssueComment(p github.IssueCommentPayload) error {
	if g.verifyProject(p.Repository.FullName) != nil {
		log.Warnf("Payload came from unverified project: %+v", p)
		if g.Discord != nil {
			g.Discord.sendLog("Repository event from unverified project")
		}
		return fmt.Errorf("unknown repository")
	}
	event := &GitHubEvent{
		event:        IssueComment,
		issueComment: p,
	}
	g.Events <- *event
	return nil
}

func (g *GitHub) Milestone(p github.MilestonePayload) error {
	if g.verifyProject(p.Repository.FullName) != nil {
		log.Warnf("Payload came from unverified project: %+v", p)
		if g.Discord != nil {
			g.Discord.sendLog("Repository event from unverified project")
		}
		return fmt.Errorf("unknown repository")
	}
	event := &GitHubEvent{
		event:     Milestone,
		milestone: p,
	}
	g.Events <- *event
	return nil
}

func (g *GitHub) PullRequest(p github.PullRequestPayload) error {
	if g.verifyProject(p.Repository.FullName) != nil {
		log.Warnf("Payload came from unverified project: %+v", p)
		if g.Discord != nil {
			g.Discord.sendLog("Repository event from unverified project")
		}
		return fmt.Errorf("unknown repository")
	}
	event := &GitHubEvent{
		event:       PullRequest,
		pullRequest: p,
	}
	g.Events <- *event
	return nil
}

func (g *GitHub) PullRequestReview(p github.PullRequestReviewPayload) error {
	if g.verifyProject(p.Repository.FullName) != nil {
		log.Warnf("Payload came from unverified project: %+v", p)
		if g.Discord != nil {
			g.Discord.sendLog("Repository event from unverified project")
		}
		return fmt.Errorf("unknown repository")
	}
	event := &GitHubEvent{
		event:             PullRequestReview,
		pullRequestReview: p,
	}
	g.Events <- *event
	return nil
}

func (g *GitHub) PullRequestComment(p github.PullRequestReviewCommentPayload) error {
	if g.verifyProject(p.Repository.FullName) != nil {
		log.Warnf("Payload came from unverified project: %+v", p)
		if g.Discord != nil {
			g.Discord.sendLog("Repository event from unverified project")
		}
		return fmt.Errorf("unknown repository")
	}
	event := &GitHubEvent{
		event:              PullRequestComment,
		pullRequestComment: p,
	}
	g.Events <- *event
	return nil
}

func (g *GitHub) Vulnerability(p github.RepositoryVulnerabilityAlertPayload) error {
	event := &GitHubEvent{
		event:         Vulnerability,
		vulnerability: p,
	}
	g.Events <- *event
	return nil
}

func (g *GitHub) SecurityAdvisory(p github.SecurityAdvisoryPayload) error {
	event := &GitHubEvent{
		event:    Security,
		security: p,
	}
	g.Events <- *event
	return nil
}
