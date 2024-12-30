package main

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/ddeville/kattungar-notify/internal/types"
)

func printDevices(devices []types.Device) {
	if len(devices) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Id", "Key", "Name", "Token"})

		for _, d := range devices {
			table.Append([]string{fmt.Sprintf("%d", d.Id), d.Key, d.Name, d.Token})
		}
		table.Render()
	} else {
		log.Println("No registered device!")
	}
}

func printNotifications(notifications []types.Notification) {
	if len(notifications) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Id", "Title", "Subtitle", "Body"})

		for _, n := range notifications {
			table.Append([]string{fmt.Sprintf("%d", n.Id), n.Title, n.Subtitle, n.Body})
		}
		table.Render()
	} else {
		log.Println("No notifications sent to device!")
	}
}
