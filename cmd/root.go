package cmd

import (
	"os"

	"github.com/ditatompel/wa-status-archiver/internal/config"
	"github.com/ditatompel/wa-status-archiver/internal/database"

	"github.com/spf13/cobra"
)

const AppVer = "0.0.1"

var LogLevel string

var rootCmd = &cobra.Command{
	Use:     "wa-status-archiver",
	Short:   "WA Status Archiver",
	Long:    `A "bot" that listen to WA WebSocket and download all videos
and images from your contact status updates.`,
	Version: AppVer,
}

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
