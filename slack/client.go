package slack

import (
	"github.com/cecchisandrone/smarthome-server/model"
	"gopkg.in/resty.v1"
)

const apiUrl = "https://slack.com/api/"
const channelListSuffix = "channels.list"
const chatPostMessageSuffix = "chat.postMessage"
const AlarmChannel = "alarm"

type Client struct {
	Configuration model.Slack
}

type Response struct {
	Ok bool `json:"ok"`
}

type ChannelListResponse struct {
	*Response
	Channels []Channel `json:"channels"`
}

type Channel struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	IsChannel bool   `json:"is_channel"`
}

func (s *Client) GetChannelList() (ChannelListResponse, error) {

	channelList := ChannelListResponse{}
	_, err := resty.R().SetResult(&channelList).SetQueryParams(map[string]string{"token": s.Configuration.Token, "scope": channelListSuffix}).Get(apiUrl + channelListSuffix)
	return channelList, err
}

func (s *Client) SendMessageToChannel(channel string, message string) error {

	channelList, err := s.GetChannelList()
	if err == nil {
		for _, c := range channelList.Channels {
			if c.Name == channel {
				Response := Response{}
				_, err := resty.R().SetResult(&Response).SetQueryParams(map[string]string{"token": s.Configuration.Token, "channel": c.Id, "text": message, "as_user": "false", "username": "SmartHome", "icon_url": "https://dl.dropboxusercontent.com/u/1580227/icons/home.png"}).Get(apiUrl + chatPostMessageSuffix)
				return err
			}
		}
	}
	return err
}
