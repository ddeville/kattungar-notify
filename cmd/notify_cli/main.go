package main

import (
	"log"

	"github.com/ddeville/kattungar-notify/internal/client"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(client.InitConfig)
}

func main() {
	if err := notifyCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
