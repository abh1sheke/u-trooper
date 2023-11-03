package cmd

import "github.com/spf13/cobra"

var url, proxy string
var duration, simul, views int

var rootCmd = &cobra.Command{
	Use:   "utroop",
	Short: "utroop Is a youtube view-bot application",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "youtube video URL")
	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.PersistentFlags().StringVarP(&proxy, "proxy", "p", "", "proxy server URL")
	rootCmd.PersistentFlags().IntVarP(&duration, "duration", "d", 50, "watch duration")
	rootCmd.PersistentFlags().IntVarP(&simul, "simul", "s", 1, "number of browser instances open simultaneously")
	rootCmd.PersistentFlags().IntVarP(&views, "views", "v", 1, "number of desired views")
}
