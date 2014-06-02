package main

import (
	"github.com/codegangsta/cli"
	"github.com/savaki/mixpanel"
	"io"
	"log"
	"os"
)

func GetClient(c *cli.Context) *mixpanel.Client {
	apiKey := os.Getenv("MIXPANEL_API_KEY")
	if value := c.String("api-key"); value != "" {
		apiKey = value
	}
	if apiKey == "" {
		log.Fatalln("ERROR: mixpanel api key not set!")
	}

	apiSecret := os.Getenv("MIXPANEL_API_SECRET")
	if value := c.String("api-secret"); value != "" {
		apiSecret = value
	}
	if apiSecret == "" {
		log.Fatalln("ERROR: mixpanel api secret not set!")
	}

	return mixpanel.New(apiKey, apiSecret)
}

func main() {
	globalFlags := []cli.Flag{
		cli.StringFlag{"api-key, k", "", "Mixpanel API key; may be specified in env as MIXPANEL_API_KEY"},
		cli.StringFlag{"api-secret, s", "", "Mixpanel API secret; may be specified in env as MIXPANEL_API_SECRET"},
	}

	app := cli.NewApp()
	app.Name = "mixpanel"
	app.Usage = "a command line interface for MixPanel"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "download",
			Usage: "download raw json events from mixpanel",
			Flags: append(globalFlags, []cli.Flag{
				cli.StringFlag{"from, f", "", "from (inclusive) date"},
				cli.StringFlag{"to, t", "", "to (inclusive) date"},
				cli.StringFlag{"event, e", "", "[optional] event type to include; defaults to all events"},
			}...),
			Action: func(c *cli.Context) {
				client := GetClient(c)
				from := c.String("from")
				to := c.String("to")
				event := c.String("event")

				stream, err := client.Download(from, to, event)
				if err != nil {
					log.Fatalln(err)
				}

				defer stream.Close()
				io.Copy(os.Stdout, stream)
			},
		},
	}

	app.Run(os.Args)
}
