package influxdb

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client"
	"github.com/spf13/viper"
	"log"
	"net/url"
)

type Client struct {
	host  string
	port int
	database string
	client *client.Client
}

func (i *Client) Init() {
	i.host = viper.GetString("influxdb.host")
	i.port = viper.GetInt("influxdb.port")
	i.database = viper.GetString("influxdb.database")

	host, err := url.Parse(fmt.Sprintf("http://%s:%d", i.host, i.port))
	if err != nil {
		log.Fatal(err)
	}
	i.client, err = client.NewClient(client.Config{URL: *host})
	if err != nil {
		log.Fatal(err)
	}
}

func (i *Client) AddPoint(point client.Point) {

	points := make([]client.Point,1)
	points[0] = point

	bps := client.BatchPoints{
		Points:          points,
		Database:        i.database,
		RetentionPolicy: "autogen",
	}

	_, err := i.client.Write(bps)
	if err != nil {
		log.Fatal(err)
	}
}
