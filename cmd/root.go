package cmd

import (
	"os"

	"wabot/internal/config"
	"wabot/internal/database"
	"github.com/spf13/cobra"
)

const AppVer = "0.0.1"

var LogLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wabot",
	Short: "ditatompel's WhatsApp bot",
	Long: `A WhatsApp bot run for ditatompel's project.`,
	Version: AppVer,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	config.LoadAll(".env")

	LogLevel = config.AppCfg().LogLevel
	// connect to DB
	if err := database.ConnectDB(); err != nil {
		panic(err)
	}
}
