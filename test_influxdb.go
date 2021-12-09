package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client"
	"log"
	"math/rand"
	"net/url"
	"time"
)

func test_influxdb() {

	rand.Seed(42)

	host, err := url.Parse(fmt.Sprintf("http://%s:%d", "localhost", 8086))
	if err != nil {
		log.Fatal(err)
	}
	con, err := client.NewClient(client.Config{URL: *host})
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	fmt.Println(t.String())

	point := client.Point{
		Measurement: "shapes",
		Tags: map[string]string{
			"color": "blue",
			"shape": "square",
		},
		Fields: map[string]interface{}{
			"value": rand.Intn(1000),
		},
		Time: t,
	}

	points := make([]client.Point, 1)
	points[0] = point

	bps := client.BatchPoints{
		Points:          points,
		Database:        "mydb",
		RetentionPolicy: "autogen",
	}

	_, err = con.Write(bps)
	if err != nil {
		log.Fatal(err)
	}
}
