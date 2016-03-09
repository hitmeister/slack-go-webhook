package main

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"github.com/parnurzeal/gorequest"
	"log"
	"os"
	"strings"
)

var (
	Version                                                      = "dev"
	url, channel, text                                           string
	fallback, color, pretext, title, title_link, attachment_text string
	fields                                                       cli.StringSlice
)

func main() {
	app := cli.NewApp()
	app.Name = "slackhook"
	app.Usage = "Sends text messages to slack channel via webhook"
	app.Version = Version
	app.Authors = []cli.Author{
		{
			Name:  "Maksim Naumov",
			Email: "maksim.naumov@hitmeister.de",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "send",
			ShortName: "s",
			Usage:     "Send message via slack webhook",
			Action:    sendCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "url",
					Usage:       "Webhook URL",
					Destination: &url,
				},
				cli.StringFlag{
					Name:        "channel",
					Usage:       "Channel name",
					Destination: &channel,
				},
				cli.StringFlag{
					Name:        "text",
					Usage:       "Message text",
					Destination: &text,
				},
				cli.StringFlag{
					Name:        "fallback",
					Usage:       "A plain-text summary of the attachment",
					Destination: &fallback,
				},
				cli.StringFlag{
					Name:        "color",
					Value:       "#439FE0",
					Usage:       "Attachment color (good, warning, danger or any hex)",
					Destination: &color,
				},
				cli.StringFlag{
					Name:        "pretext",
					Usage:       "This is optional text that appears above the message attachment block",
					Destination: &pretext,
				},
				cli.StringFlag{
					Name:        "title",
					Usage:       "The title is displayed as larger, bold text near the top of a message attachment",
					Destination: &title,
				},
				cli.StringFlag{
					Name:        "title_link",
					Usage:       "By passing a valid URL in the title_link parameter (optional), the title text will be hyperlinked",
					Destination: &title_link,
				},
				cli.StringFlag{
					Name:        "attachment_text",
					Usage:       "This is the main text in a message attachment",
					Destination: &attachment_text,
				},
				cli.StringSliceFlag{
					Name:  "field",
					Usage: "Fields are defined as an array, and hashes contained within it will be displayed in a table inside the message attachment",
					Value: &fields,
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error on run app, %v", err)
	}
}

func sendCommand(ctx *cli.Context) {
	if url == "" {
		log.Fatal("Url is empty")
	}
	if channel == "" {
		log.Fatal("Channel is empty")
	}
	if text == "" {
		log.Fatal("Message text is empty")
	}

	// Normalize
	if channel[0] != '#' {
		channel = "#" + channel
	}

	attach := Attachment{
		Fallback:  &fallback,
		Color:     &color,
		Title:     &title,
		TitleLink: &title_link,
		Text:      &attachment_text,
	}

	for _, field := range fields {
		kv := strings.SplitN(field, ":", 2)
		if len(kv) > 1 && kv[0] != "" && kv[1] != "" {
			attach.Fields = append(attach.Fields, &Field{
				Title: kv[0],
				Value: kv[1],
				Short: true,
			})
		}
	}

	payload := make(map[string]interface{})
	payload["channel"] = channel
	payload["text"] = text
	payload["attachments"] = []Attachment{attach}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("error on encode request, %v", err)
	}

	_, _, errors := gorequest.New().Post(url).Send(string(data)).End()
	if len(errors) > 0 {
		log.Fatalf("error on send request, %#v", errors)
	}
}
