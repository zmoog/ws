/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zmoog/ws/v2/feedback"
	"github.com/zmoog/ws/v2/ws"
	"github.com/zmoog/ws/v2/ws/identity"
)

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Manage devices",
	Long:  `Manage devices in your account.`,
}

// listCmd represents the list command
var listDevicesCmd = &cobra.Command{
	Use:   "list",
	Short: "List devices",
	Long:  `List the devices in your account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		identityManager := identity.NewManager(
			identity.Config{
				Username:  viper.GetString("username"),
				Password:  viper.GetString("password"),
				WebApiKey: viper.GetString("web_api_key"),
			},
		)
		client := ws.NewClient(
			identityManager,
			viper.GetString("api_endpoint"),
		)

		devices, err := client.ListDevices()
		if err != nil {
			return fmt.Errorf("failed to list devices: %v", err)
		}

		_ = feedback.PrintResult(deviceResult{Devices: devices})

		return nil
	},
}

type deviceResult struct {
	Devices []ws.Device `json:"devices"`
}

func (r deviceResult) Table() string {

	table := pterm.TableData{}
	table = append(table, []string{
		"Name",
		"Serial Number",
		"Firmware Available",
		"Firmware Installed",
		"Type",
		"Last Heartbeat",
	})

	for _, device := range r.Devices {
		table = append(table, []string{
			device.Name,
			device.SerialNumber,
			device.FirmwareAvailable,
			device.FirmwareInstalled,
			device.Type,
			device.LastHeartbeat.Format("2006-01-02 15:04:05"),
		})
	}

	rendered, _ := pterm.DefaultTable.WithHasHeader().WithData(table).Srender()
	return rendered
}

func (r deviceResult) String() string {
	return r.Table()
}

func (r deviceResult) Data() any {
	return r.Devices
}

func init() {
	rootCmd.AddCommand(devicesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	devicesCmd.AddCommand(listDevicesCmd)
}
