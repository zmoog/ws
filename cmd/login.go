package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zmoog/ws/v2/feedback"
	"github.com/zmoog/ws/v2/ws/identity"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the Wavin API",
	Long:  `Login to the Wavin API.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		feedback.Println("Login to the Wavin API")
		im := identity.NewManager(
			identity.Config{
				Username:  viper.GetString("username"),
				Password:  viper.GetString("password"),
				WebApiKey: viper.GetString("web_api_key"),
			},
		)

		token, err := im.GetToken()
		if err != nil {
			return fmt.Errorf("failed to get token: %w", err)
		}

		feedback.Println(
			fmt.Sprintf("Login successful, expires at %s", token.ExpiresAt.Format("2006-01-02 15:04:05 -0700 MST")),
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
