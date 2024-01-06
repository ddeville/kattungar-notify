package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "kattungar-notify-admin",
}

func setupCommands() {
	var cmdListDevices = &cobra.Command{
		Use:   "list-devices",
		Short: "List all devices",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("Listing devices...")
		},
	}
	rootCmd.AddCommand(cmdListDevices)

	var cmdAddDevice = &cobra.Command{
		Use:   "add-device",
		Short: "Add a new device",
		Run: func(cmd *cobra.Command, _ []string) {
			key, _ := cmd.Flags().GetString("key")
			name, _ := cmd.Flags().GetString("name")
			fmt.Printf("Adding device with key: %s and name: %s\n", key, name)
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
			fmt.Printf("Deleting device with key: %s\n", key)
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
			fmt.Printf("Updating device with key: %s and new name: %s\n", key, name)
		},
	}
	cmdUpdateDevice.Flags().String("key", "", "Key of the device")
	cmdUpdateDevice.Flags().String("name", "", "New name of the device")
	cmdUpdateDevice.MarkFlagRequired("key")
	cmdUpdateDevice.MarkFlagRequired("name")
	rootCmd.AddCommand(cmdUpdateDevice)
}

func main() {
	setupCommands()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
