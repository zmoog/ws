/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zmoog/ws/feedback"
	"github.com/zmoog/ws/ws"
	"github.com/zmoog/ws/ws/identity"
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
			viper.GetString("username"),
			viper.GetString("password"),
		)
		client := ws.NewClient(
			identityManager,
			viper.GetString("api_endpoint"),
		)

		rooms, err := client.ListRooms(ulc)
		if err != nil {
			return fmt.Errorf("failed to list rooms: %w", err)
		}

		_ = feedback.PrintResult(roomsListResult{rooms: rooms})

		return nil
	},
}

type roomsListResult struct {
	rooms []ws.Room
}

func (r roomsListResult) Table() string {

	table := pterm.TableData{}
	table = append(table, []string{
		"Name",
		"Status",
		"Temperature (desired)",
		"Temperature (current)",
		"Humidity (current)",
	})

	for _, room := range r.rooms {
		table = append(table, []string{
			room.Name,
			room.Status,
			fmt.Sprintf("%.1f", room.TempDesired),
			fmt.Sprintf("%.1f", room.TempCurrent),
			fmt.Sprintf("%.1f", room.HumidityCurrent),
		})
	}

	if err := pterm.DefaultTable.WithHasHeader().WithData(table).Render(); err != nil {
		return fmt.Sprintf("failed to render table: %v", err)
	}

	return ""
}

func (r roomsListResult) String() string {
	return r.Table()
}

func (r roomsListResult) Data() any {
	return r.rooms
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

	listRoomsCmd.Flags().StringVarP(&ulc, "location-id", "l", "", "Location ID")
	_ = listRoomsCmd.MarkFlagRequired("location-id")
}
