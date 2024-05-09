package main

import (
	"fmt"
	"main/ipinfo"
	"main/webhook"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	viper.SetDefault("DELAY", "1h")

	delay := viper.GetString("DELAY")
	wh_url := viper.GetString("WEBHOOK_URL")
	if wh_url == "" {
		panic("WEBHOOK_URL is not set")
	}

	errChan := make(chan error)
	runner := cron.New()

	var prevInfo *ipinfo.Info

	// Add cron job
	_, err := runner.AddFunc("@every "+delay, func() {
		// Get info
		info, err := ipinfo.GetInfo()
		if err != nil {
			errChan <- err
			return
		}

		// Set default embed title and description
		title := "New ip detected"
		description := fmt.Sprintf("||`%s`||", info.IP)

		// Check info is changed
		if prevInfo.IP == info.IP {
			return
		}

		msg := &webhook.Message{
			Embeds: &[]webhook.Embed{
				{
					Title:       &title,
					Description: &description,
				},
			},
		}

		// Send message
		err = webhook.SendMessage(wh_url, msg)
		if err != nil {
			errChan <- err
			return
		}

		// Save info
		prevInfo = info
	})
	if err != nil {
		panic(err)
	}

	// Start first run

	// Get info
	info, err := ipinfo.GetInfo()
	if err != nil {
		errChan <- err
		return
	}

	// Set default embed title and description
	title := "Started with ip"
	description := fmt.Sprintf("||`%s`||", info.IP)

	// Check info is changed
	if prevInfo != nil && prevInfo.IP == info.IP {
		return
	}

	msg := &webhook.Message{
		Embeds: &[]webhook.Embed{
			{
				Title:       &title,
				Description: &description,
			},
		},
	}

	// Send message
	err = webhook.SendMessage(wh_url, msg)
	if err != nil {
		errChan <- err
		return
	}

	// Save info
	prevInfo = info

	// Start runner
	runner.Start()

	for {
		err := <-errChan
		if err != nil {
			panic(err)
		}
	}
}
