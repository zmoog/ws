package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zmoog/ws/v2/ws/identity"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the Wavin API",
	Long:  `Login to the Wavin API.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Login to the Wavin API")

		im := identity.NewManager(
			viper.GetString("username"),
			viper.GetString("password"),
			viper.GetString("web_api_key"),
		)

		token, err := im.GetToken()
		if err != nil {
			return fmt.Errorf("failed to get token: %w", err)
		}

		fmt.Println("Login successful", token)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
