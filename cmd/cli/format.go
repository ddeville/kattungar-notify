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
