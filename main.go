package main

import (
	"fmt"
	"log"

	"github.com/docopt/docopt-go"
	"github.com/kovetskiy/pulseaudio"
	"github.com/reconquest/karma-go"
)

var (
	version = "[manual build]"
	usage   = "pa-switch-profile " + version + `

Usage:
  pa-switch-profile [options] <card> <profile>...
  pa-switch-profile -h | --help
  pa-switch-profile --version

Options:
  -h --help  Show this screen.
  --version  Show version.
`
)

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(err)
	}

	client, err := pulseaudio.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	card, err := getCard(client, args["<card>"].(string))
	if err != nil {
		log.Fatalln(err)
	}

	err = switchProfile(client, card, args["<profile>"].([]string))
	if err != nil {
		log.Fatalln(err)
	}
}

func switchProfile(
	client *pulseaudio.Client,
	card *pulseaudio.Card,
	profiles []string,
) error {
	var afterActive bool
	var next string
	for _, profile := range profiles {
		if !afterActive {
			if card.ActiveProfile.Name == profile {
				afterActive = true
			}

			continue
		}

		next = profile
		break
	}

	if next == "" {
		next = profiles[0]
	}

	err := client.SetCardProfile(card.Index, next)
	if err != nil {
		return karma.Format(
			err,
			"unable to set card profile to %s", next,
		)
	}

	return nil
}

func getCard(client *pulseaudio.Client, query string) (*pulseaudio.Card, error) {
	cards, err := client.Cards()
	if err != nil {
		log.Fatalln(err)
	}

	var targetCard pulseaudio.Card
	for _, card := range cards {
		if card.Name == query {
			targetCard = card
			break
		}

		if fmt.Sprint(card.Index) == query {
			targetCard = card
			break
		}
	}

	return &targetCard, nil
}
