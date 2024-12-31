package main

import (
	"log"
)

func main() {
	if err := notifyCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
