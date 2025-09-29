/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zmoog/ws/v2/feedback"
	"github.com/zmoog/ws/v2/ws"
	"github.com/zmoog/ws/v2/ws/identity"
)

var (
	ulc string
)

// roomsCmd represents the rooms command
var roomsCmd = &cobra.Command{
	Use:   "rooms",
	Short: "Manage rooms",
	Long:  `Manage rooms in your account.`,
}

// listCmd represents the list command
var listRoomsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the rooms",
	Long:  `List the rooms in a location.`,
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

		device, err := client.GetDevice(ulc)
		if err != nil {
			return fmt.Errorf("failed to get device: %w", err)
		}

		_ = feedback.PrintResult(roomsListResult{device: device})

		return nil
	},
}

type roomsListResult struct {
	device ws.Device
}

func (r roomsListResult) Table() string {
	var sb strings.Builder

	table := pterm.TableData{}
	table = append(table, []string{
		"Name",
		"Temperature state",
		"Temperature (desired)",
		"Temperature (current)",
		"Humidity (current)",
		"Dehumidification state",
	})

	for _, room := range r.device.LastConfig.Sentio.Rooms {
		table = append(table, []string{
			room.Title,
			room.TemperatureState,
			fmt.Sprintf("%.1f", room.SetpointTemperature),
			fmt.Sprintf("%.1f", room.AirTemperature),
			fmt.Sprintf("%.1f", room.Humidity),
			room.DehumidifierState,
		})
	}

	rendered, err := pterm.DefaultTable.WithHasHeader().WithData(table).Srender()
	if err != nil {
		return fmt.Sprintf("failed to render table: %s", err)
	}
	sb.WriteString(rendered)

	// Add a newline after the table
	sb.WriteString("\n")

	for _, room := range r.device.LastConfig.Sentio.OutdoorTemperatureSensors {
		sb.WriteString(fmt.Sprintf("Outdoor temperature: %.1f\n\n", room.OutdoorTemperature))
	}

	return sb.String()
}

func (r roomsListResult) String() string {
	return r.Table()
}

func (r roomsListResult) Data() any {
	return r.device.LastConfig.Sentio.Rooms
}

func init() {
	rootCmd.AddCommand(roomsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// roomsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// roomsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	roomsCmd.AddCommand(listRoomsCmd)

	listRoomsCmd.Flags().StringVarP(&ulc, "device-name", "d", "", "Device name")
	_ = listRoomsCmd.MarkFlagRequired("device-name")
}
