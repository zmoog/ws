/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zmoog/ws/ws"
	"github.com/zmoog/ws/ws/identity"
)

var (
	ulc string
)

// roomsCmd represents the rooms command
var roomsCmd = &cobra.Command{
	Use:   "rooms",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rooms called")
	},
}

// listCmd represents the list command
var listRoomsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the rooms",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")

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
			log.Fatalf("failed to list rooms: %v", err)
		}

		table := pterm.TableData{}
		table = append(table, []string{
			"Name",
			"Thermo",
			"Dryer",
			"TempDesired",
			"TempCurrent",
		})

		for _, room := range rooms {
			table = append(table, []string{
				room.Name,
				room.Thermo,
				room.Dryer,
				fmt.Sprintf("%.1f", room.TempDesired),
				fmt.Sprintf("%.1f", room.TempCurrent),
			})
		}

		// pterm.DefaultTable.WithData(table).Render()
		pterm.DefaultTable.WithHasHeader().WithData(table).Render()
	},
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
	listRoomsCmd.MarkFlagRequired("location-id")
}