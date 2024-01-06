package main

import (
	"encoding/json"
	"log"

	"github.com/spf13/cobra"

	"github.com/ddeville/kattungar-notify/internal/store"
)

var rootCmd = &cobra.Command{
	Use: "kattungar-notify-admin",
}

func init() {
	var cmdListDevices = &cobra.Command{
		Use:   "list-devices",
		Short: "List all devices",
		Run: func(_ *cobra.Command, _ []string) {
			res := makeRequest("GET", "https://notify.home.kattungar.net/admin/device", nil, nil)

			var devices []store.Device
			defer res.Body.Close()
			if err := json.NewDecoder(res.Body).Decode(&devices); err != nil {
				log.Fatal(err)
			}

			printDevices(devices)
		},
	}
	rootCmd.AddCommand(cmdListDevices)

	var cmdAddDevice = &cobra.Command{
		Use:   "add-device",
		Short: "Add a new device",
		Run: func(cmd *cobra.Command, _ []string) {
			key, _ := cmd.Flags().GetString("key")
			name, _ := cmd.Flags().GetString("name")

			body, err := json.Marshal(store.Device{
				Key:  key,
				Name: name,
			})
			if err != nil {
				log.Fatalln(err)
			}

			res := makeRequest("POST", "https://notify.home.kattungar.net/admin/device", body, nil)

			var device store.Device
			defer res.Body.Close()
			if err := json.NewDecoder(res.Body).Decode(&device); err != nil {
				log.Fatal(err)
			}

			printDevices([]store.Device{device})
		},
	}
	cmdAddDevice.Flags().String("key", "", "Key of the device")
	cmdAddDevice.Flags().String("name", "", "Name of the device")
	cmdAddDevice.MarkFlagRequired("key")
	cmdAddDevice.MarkFlagRequired("name")
	rootCmd.AddCommand(cmdAddDevice)

	var cmdDeleteDevice = &cobra.Command{
		Use:   "delete-device",
		Short: "Delete a device",
		Run: func(cmd *cobra.Command, _ []string) {
			key, _ := cmd.Flags().GetString("key")

			body, err := json.Marshal(store.Device{
				Key: key,
			})
			if err != nil {
				log.Fatalln(err)
			}

			_ = makeRequest("DELETE", "https://notify.home.kattungar.net/admin/device", body, nil)
			log.Printf("Deleted device with key: %s\n", key)
		},
	}
	cmdDeleteDevice.Flags().String("key", "", "Key of the device")
	cmdDeleteDevice.MarkFlagRequired("key")
	rootCmd.AddCommand(cmdDeleteDevice)

	var cmdUpdateDevice = &cobra.Command{
		Use:   "update-device",
		Short: "Update a device",
		Run: func(cmd *cobra.Command, _ []string) {
			key, _ := cmd.Flags().GetString("key")
			name, _ := cmd.Flags().GetString("name")

			body, err := json.Marshal(store.Device{
				Name: name,
			})
			if err != nil {
				log.Fatalln(err)
			}

			_ = makeRequest("PUT", "https://notify.home.kattungar.net/device/name", body, &key)
			log.Printf("Updated device name to \"%s\"\n", key)
		},
	}
	cmdUpdateDevice.Flags().String("key", "", "Key of the device")
	cmdUpdateDevice.Flags().String("name", "", "New name of the device")
	cmdUpdateDevice.MarkFlagRequired("key")
	cmdUpdateDevice.MarkFlagRequired("name")
	rootCmd.AddCommand(cmdUpdateDevice)
}
