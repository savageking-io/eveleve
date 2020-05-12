package main

import (
	"fmt"
	"net"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type Master struct {
	Config   *Config
	GitHub   *GitHub
	Travis   *Travis
	Discord  *Discord
	Status   *Status
	Listener *net.TCPListener
}

func (m *Master) Init() error {
	log.Infof("Starting EvelEve in Master mode")
	log.Infof("App version: %s", AppVersion)

	if err := m.InitConfig(); err != nil {
		m.Config = nil
		log.Errorf("%s", err.Error())
	}

	if err := m.InitGitHub(); err != nil {
		log.Errorf("%s", err.Error())
	}

	if err := m.InitTravis(); err != nil {
		log.Errorf("%s", err.Error())
	}

	if err := m.InitDiscord(); err != nil {
		log.Errorf("%s", err.Error())
	}

	if err := m.InitAPI(); err != nil {
		log.Errorf("%s", err.Error())
	}

	return nil

}

func (m *Master) InitConfig() error {
	m.Config = new(Config)
	if err := m.Config.Init("/etc/eveleve/config.yaml"); err != nil {
		m.Config = nil
		return fmt.Errorf("Failed to initialize configuration subsystem: %s", err.Error())
	}
	return nil
}

func (m *Master) InitGitHub() error {
	if m.Config == nil {
		return fmt.Errorf("Skipping GitHub initialization due to an empty configuration")
	}
	m.GitHub = new(GitHub)
	if err := m.GitHub.Init(m.Config.GitHub, m.Config.TLS); err != nil {
		m.GitHub = nil
		return fmt.Errorf("Failed to initialize GitHub subsystem: %s", err.Error())
	}
	m.GitHub.SetProjects(m.Config.Projects)
	return nil
}

func (m *Master) InitTravis() error {
	if m.Config == nil {
		return fmt.Errorf("Skipping Travis initialziation due to an empty configuration")
	}
	m.Travis = new(Travis)
	if err := m.Travis.Init(&m.Config.Travis); err != nil {
		m.Travis = nil
		return fmt.Errorf("Failed to initialize Travis subsystem: %s", err.Error())
	}
	return nil
}

func (m *Master) InitDiscord() error {
	if m.Config == nil {
		return fmt.Errorf("Skipping Discord initialization due to empty configuration")
	}
	m.Discord = new(Discord)
	if err := m.Discord.Init(m.Config.Discord); err != nil {
		m.Discord = nil
		return fmt.Errorf("%s", err.Error())
	}

	if m.Discord != nil {
		m.Discord.sendLog("EvelEve Bot Online. All Systems Nominal")

		for i := 0; i < 255; i++ {
			msg := new(discordgo.MessageEmbed)
			msg.Title = fmt.Sprintf("Testing color %d", i)
			msg.Color = i
			msg.Author = new(discordgo.MessageEmbedAuthor)
			msg.Author.Name = "Ivan Ivanovich"

			m.Discord.sendEmbed(m.Config.Discord.LogChannel, msg)
			time.Sleep(time.Second * 2)
		}

		if m.GitHub != nil {
			m.GitHub.addNotificationSubsystem(m.Discord)
		}
	}

	log.Infof("Initializing Status Subsystem")
	m.Status = new(Status)
	if err := m.Status.Init(m.Discord); err != nil {
		log.Errorf("Failed to initialize Status Subsystem: %s", err.Error())
	}

	return nil
}

func (m *Master) InitAPI() error {
	return nil
}

func (m *Master) Run() error {

	log.Infof("Running Status Subsystem")
	go m.Status.Run()
	go m.Travis.Run()

	for {
		if m.Discord == nil || m.GitHub == nil || m.Config == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		select {
		case cmd := <-m.Discord.Commands:
			handleCommand(cmd)
		case gevent := <-m.GitHub.Events:
			log.Infof("New repo event: %+v", gevent)
		case tevent := <-m.Travis.Events:
			log.Infof("New Travis Event: %+v", tevent)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}

	return nil
}

func handleCommand(command Command) error {
	switch command.Cmd {
	case "!projects":

	}

	return nil
}

func handleEvent(event GitHubEvent) error {

	return nil
}
