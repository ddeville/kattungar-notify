package main

import (
	"encoding/json"
	"log"

	"github.com/spf13/cobra"

	"github.com/ddeville/kattungar-notify/internal/client"
	"github.com/ddeville/kattungar-notify/internal/types"
)

func init() {
	var cmdListDevices = &cobra.Command{
		Use:   "list-devices",
		Short: "List all devices",
		Run: func(_ *cobra.Command, _ []string) {
			apiKey := client.C.APIKey
			res := client.MakeRequest("GET", "/admin/device", nil, &apiKey)

			var devices []types.Device
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

			body, err := json.Marshal(types.Device{
				Key:  key,
				Name: name,
			})
			if err != nil {
				log.Fatalln(err)
			}

			apiKey := client.C.APIKey
			res := client.MakeRequest("POST", "/admin/device", body, &apiKey)

			var device types.Device
			defer res.Body.Close()
			if err := json.NewDecoder(res.Body).Decode(&device); err != nil {
				log.Fatal(err)
			}

			printDevices([]types.Device{device})
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

			body, err := json.Marshal(types.Device{
				Key: key,
			})
			if err != nil {
				log.Fatalln(err)
			}

			apiKey := client.C.APIKey
			_ = client.MakeRequest("DELETE", "/admin/device", body, &apiKey)
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

			body, err := json.Marshal(types.Device{
				Name: name,
			})
			if err != nil {
				log.Fatalln(err)
			}

			apiKey := client.C.APIKey
			_ = client.MakeRequest("PUT", "/device/name", body, &apiKey)
			log.Printf("Updated device name to \"%s\"\n", key)
		},
	}
	cmdUpdateDevice.Flags().String("key", "", "Key of the device")
	cmdUpdateDevice.Flags().String("name", "", "New name of the device")
	cmdUpdateDevice.MarkFlagRequired("key")
	cmdUpdateDevice.MarkFlagRequired("name")
	rootCmd.AddCommand(cmdUpdateDevice)

	var cmdListNotifications = &cobra.Command{
		Use:   "list-notifications",
		Short: "List all notifications for a device",
		Run: func(cmd *cobra.Command, _ []string) {
			key, _ := cmd.Flags().GetString("key")

			res := client.MakeRequest("GET", "/device/list_notifications", nil, &key)

			var notifications []types.Notification
			defer res.Body.Close()
			if err := json.NewDecoder(res.Body).Decode(&notifications); err != nil {
				log.Fatal(err)
			}

			printNotifications(notifications)
		},
	}
	cmdListNotifications.Flags().String("key", "", "Key of the device")
	cmdListNotifications.MarkFlagRequired("key")
	rootCmd.AddCommand(cmdListNotifications)
}
