package main

import (
	"encoding/json"
	"log"

	"github.com/spf13/cobra"

	"github.com/ddeville/kattungar-notify/internal/client"
	"github.com/ddeville/kattungar-notify/internal/types"
)

var notifyCmd = &cobra.Command{
	Use:   "kattungar-notify",
	Short: "Send a notification to a device",
	Run: func(cmd *cobra.Command, _ []string) {
		key, _ := cmd.Flags().GetString("key")
		name, _ := cmd.Flags().GetString("name")
		title, _ := cmd.Flags().GetString("title")
		subtitle, _ := cmd.Flags().GetString("subtitle")
		body, _ := cmd.Flags().GetString("body")

		if len(key) == 0 && len(name) == 0 {
			log.Fatalln("You need to provide a device key or name!")
		}
		if len(key) > 0 && len(name) > 0 {
			log.Fatalln("You should only provide either a device key or name!")
		}

		if len(title) == 0 && len(subtitle) == 0 && len(body) == 0 {
			log.Fatalln("You need to provide a title, subtitle, or body!")
		}

		requestBody, err := json.Marshal(types.Notification{
			DeviceKey:  key,
			DeviceName: name,
			Title:      title,
			Subtitle:   subtitle,
			Body:       body,
		})
		if err != nil {
			log.Fatalln(err)
		}

		_ = client.MakeRequest("POST", "/notify", requestBody, nil)
		log.Println("Notification sent!")
	},
}

func init() {
	notifyCmd.Flags().String("key", "", "Key of the device (use either key or name but not both)")
	notifyCmd.Flags().String("name", "", "Name of the device (use either key or name but not both)")
	notifyCmd.Flags().String("title", "", "Title of the notification")
	notifyCmd.Flags().String("subtitle", "", "Subtitle of the notification")
	notifyCmd.Flags().String("body", "", "Body of the notification")
}
