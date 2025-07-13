/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zmoog/ws/v2/feedback"
)

var (
	cfgFile string
	output  string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ws",
	Short: "Display and control Wavin Sentio floor heating system",
	Long: `ws is a command line tool to display and control the Wavin Sentio floor heating system.

It allows you to list locations and rooms, (in future versions) set the desired temperature
and humidity, and more.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet("username") {
			_ = cmd.MarkFlagRequired("username")
		}
		if !viper.IsSet("password") {
			_ = cmd.MarkFlagRequired("password")
		}

		output := viper.GetString("output")

		switch output {
		case "table":
			feedback.SetFormat(feedback.Table)
		case "text":
			feedback.SetFormat(feedback.Text)
		case "json":
			feedback.SetFormat(feedback.JSON)
		default:
			feedback.Error(fmt.Sprintf("invalid output format: %s", output))
			feedback.SetFormat(feedback.Table)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ws/config)")

	rootCmd.PersistentFlags().StringP("username", "u", "", "The username to use for authentication")
	rootCmd.PersistentFlags().StringP("password", "p", "", "The password to use for authentication")
	rootCmd.PersistentFlags().StringP("web-api-key", "k", "AIzaSyBlAtNI7-2jitPul9I-O4EZcT-n0sIay-g", "The Firebase Web API key for the Wavin Sentio app")
	rootCmd.PersistentFlags().StringP("api-endpoint", "e", "https://blaze.wavinsentio.com/wavin.blaze.v1.BlazeDeviceService", "The API endpoint to use")

	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "The format to use for output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configDir := filepath.Join(home, ".ws")

		// Search config in home directory with name ".ws" (without extension).
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("WS")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	_ = viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	_ = viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	_ = viper.BindPFlag("web_api_key", rootCmd.PersistentFlags().Lookup("web-api-key"))
	_ = viper.BindPFlag("api_endpoint", rootCmd.PersistentFlags().Lookup("api-endpoint"))
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}
