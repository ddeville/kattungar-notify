package main

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/ddeville/kattungar-notify/internal/store"
)

func printDevices(devices []store.Device) {
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

func printNotifications(notifications []store.Notification) {
	if len(notifications) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Id", "DeviceName", "Title", "Subtitle", "Body"})

		for _, n := range notifications {
			table.Append([]string{fmt.Sprintf("%d", n.Id), n.DeviceName, n.Title, n.Subtitle, n.Body})
		}
		table.Render()
	} else {
		log.Println("No notifications sent to device!")
	}
}
