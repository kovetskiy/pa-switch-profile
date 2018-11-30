package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/docopt/docopt-go"
	"github.com/kovetskiy/pulseaudio"
	"github.com/reconquest/karma-go"
)

var (
	version = "[manual build]"
	usage   = "pa-switch-profile " + version + `

Usage:
  pa-switch-profile [options] <card> <profile>... [-i <match>]...
  pa-switch-profile -h | --help
  pa-switch-profile --version

Options:
  <card>               Index of card or 'active' to find active.
  -i --ignore <match>  Ignore specified card.
  -h --help            Show this screen.
  --version            Show version.
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

	profiles := args["<profile>"].([]string)

	card, err := getCard(
		client,
		args["<card>"].(string),
		args["--ignore"].([]string),
		profiles,
	)
	if err != nil {
		log.Fatalln(err)
	}

	err = switchProfile(client, card, profiles)
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
			if card.ActiveProfile == nil {
				panic("card without active profile specified")
			}
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

func getCard(
	client *pulseaudio.Client,
	query string,
	ignores []string,
	profiles []string,
) (*pulseaudio.Card, error) {
	cards, err := client.Cards()
	if err != nil {
		log.Fatalln(err)
	}

	var targetCard pulseaudio.Card
	var found bool
	for _, card := range cards {
		if card.Name == query {
			found = true
		}

		if fmt.Sprint(card.Index) == query {
			found = true
		}

		if query == "active" && card.ActiveProfile != nil {
			found = true
		}

		if found {
			for _, match := range ignores {
				ok, err := regexp.MatchString(match, card.Name)
				if err != nil {
					log.Printf("unable to match by pattern: %q: %v", err)
				}

				if ok {
					found = false
				}
			}
		}

		if found {
			targetCard = card
		}

		for _, required := range profiles {
			requiredFound := false
			for profile, _ := range card.Profiles {
				if profile == required {
					requiredFound = true
					break
				}
			}

			if !requiredFound {
				found = false
			}
		}

		if found {
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("unable to find card: %q", query)
	}

	return &targetCard, nil
}
