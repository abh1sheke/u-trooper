package cmd

import (
	"os"

	localog "github.com/abh1sheke/utrooper/log"
	"github.com/abh1sheke/utrooper/viewer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var url, proxy string
var duration, instances, views int
var logLevel uint32

var rootCmd = &cobra.Command{
	Use:   "utroop",
	Short: "utroop Is a youtube view-bot application",
	Run: func(cmd *cobra.Command, args []string) {
		localog.Init(logLevel)
		log.WithFields(log.Fields{
			"URL":       url,
			"INSTANCES": instances,
		}).Info("Starting up")
		viewer.StartViewing(views, instances, duration, &url)
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	log.SetOutput(os.Stdout)
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "youtube video URL")
	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.PersistentFlags().StringVarP(&proxy, "proxy", "p", "", "proxy server URL")
	rootCmd.PersistentFlags().IntVarP(&duration, "duration", "d", 50, "watch duration (seconds)")
	rootCmd.PersistentFlags().IntVarP(&instances, "instances", "i", 1, "number of browser instances open simultaneously")
	rootCmd.PersistentFlags().IntVarP(&views, "views", "v", 1, "number of desired views")
	rootCmd.PersistentFlags().Uint32VarP(&logLevel, "log", "l", 5, "log with level n (0-6)")
}
