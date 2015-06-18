package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
	influxdb_client "github.com/influxdb/influxdb/client"
)

const (
	DATABASE    = "intelmaker"
	MEASUREMENT = "audience"

	POSITIVE_VOTE = "positive"
	NEGATIVE_VOTE = "negative"

	MaxUint16      = ^uint16(0)
	CLAP_TRESSHOLD = 1850 // 1911
)

type Client struct {
	*influxdb_client.Client
}

func NewClient() *Client {
	host := os.Getenv("INFLUX_HOST")
	port, err := strconv.Atoi(os.Getenv("INFLUX_PORT"))
	if err != nil {
		log.Fatal("INFLUX_PORT must ben integer")
	}

	c, err := influxdb_client.NewClient(&influxdb_client.ClientConfig{
		Host:     fmt.Sprintf("%s:%d", host, port),
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PWD"),
		Database: DATABASE,
	})
	if err != nil {
		panic(err)
	}

	return &Client{c}
}

func (c Client) send(column string) {
	series := &influxdb_client.Series{
		Name:    MEASUREMENT,
		Columns: []string{column},
		Points: [][]interface{}{
			{1},
		},
	}
	log.Printf("Sending %s...", column)
	if err := c.WriteSeries([]*influxdb_client.Series{series}); err != nil {
		log.Println(err)
	}
}

func (c Client) Vote(value string) {
	c.send(value)
}

func (c Client) Clapping() {
	c.send("clapping")
}

func main() {
	gbot := gobot.NewGobot()

	edisonAdaptor := edison.NewEdisonAdaptor("edison")

	buttonPositive := gpio.NewButtonDriver(edisonAdaptor, "button_positive", "4")
	buttonNegative := gpio.NewButtonDriver(edisonAdaptor, "button_negative", "3")

	redLed := gpio.NewLedDriver(edisonAdaptor, "red_led", "7")
	greenLed := gpio.NewLedDriver(edisonAdaptor, "green_led", "8")
	blueLed := gpio.NewLedDriver(edisonAdaptor, "blue_led", "6")

	soundSensor := gpio.NewAnalogSensorDriver(edisonAdaptor, "sound_sensor", "0")

	client := NewClient()

	work := func() {
		gobot.On(buttonPositive.Event("push"), func(data interface{}) {
			go func() {
				blueLed.Off()
				greenLed.On()
				client.Vote(POSITIVE_VOTE)
				time.Sleep(1 * time.Second)
				greenLed.Off()
			}()
		})

		gobot.On(buttonNegative.Event("push"), func(data interface{}) {
			redLed.On()
			client.Vote(NEGATIVE_VOTE)
			time.Sleep(1 * time.Second)
			redLed.Off()
		})

		gobot.On(soundSensor.Event("data"), func(data interface{}) {
			level := uint16(
				gobot.ToScale(gobot.FromScale(float64(data.(int)), 0, float64(MaxUint16)), 0, float64(MaxUint16)),
			)
			if level > CLAP_TRESSHOLD {
				blueLed.On()
				client.Clapping()
				time.Sleep(1 * time.Second)
				blueLed.Off()
			}
		})
	}

	robot := gobot.NewRobot(
		"buttonBot",
		[]gobot.Connection{edisonAdaptor},
		[]gobot.Device{buttonPositive, buttonNegative, redLed, greenLed, blueLed, soundSensor},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
