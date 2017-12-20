package service

import "github.com/cecchisandrone/smarthome-server/slack"

type Notification struct {
	ConfigurationService *Configuration `inject:""`
	SlackClient          slack.Client
}

func (n *Notification) Init() {
	configuration := n.ConfigurationService.GetCurrent()
	n.SlackClient = slack.Client{configuration.Slack}
}

func (n *Notification) SendSlackMessage(channel string, message string) error {
	return n.SlackClient.SendMessageToChannel(channel, message)
}
