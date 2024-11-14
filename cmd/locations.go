/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zmoog/ws/ws"
	"github.com/zmoog/ws/ws/identity"
)

// locationsCmd represents the locations command
var locationsCmd = &cobra.Command{
	Use:   "locations",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("username") == "" || viper.GetString("password") == "" {
			return fmt.Errorf("username and password are required")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("locations called", viper.GetString("username"), viper.GetString("password"))
		return nil
	},
}

// listCmd represents the list command
var listLocationsCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("username") == "" || viper.GetString("password") == "" {
			return fmt.Errorf("username and password are required")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("list called")

		identityManager := identity.NewManager(viper.GetString("username"), viper.GetString("password"))
		client := ws.NewClient(identityManager, viper.GetString("api_endpoint"))

		locations, err := client.ListLocations()
		if err != nil {
			return fmt.Errorf("failed to list locations: %v", err)
		}

		table := pterm.TableData{}
		table = append(table, []string{
			"Ulc",
			"Registration",
			"Serial Number",
			"Mode",
			"Vacation On",
			"Outdoor Temperature",
			"DST",
		})

		for _, location := range locations {
			// fmt.Printf("location: %+v\n", location)
			table = append(table, []string{
				location.Ulc,
				location.Registration,
				fmt.Sprintf("%d", location.SerialNumber),
				location.Attributes.Mode,
				fmt.Sprintf("%t", location.Attributes.VacationOn),
				fmt.Sprintf("%.1f", location.Attributes.Outdoor.Temperature),
				fmt.Sprintf("%t", location.Attributes.Dst),
			})
		}

		pterm.DefaultTable.WithHasHeader().WithData(table).Render()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(locationsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// locationsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// locationsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	locationsCmd.AddCommand(listLocationsCmd)
}
